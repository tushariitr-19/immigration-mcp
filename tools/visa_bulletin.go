package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
	"golang.org/x/net/html"

	"github.com/tushariitr-19/immigration-mcp/logger"
	"github.com/tushariitr-19/immigration-mcp/models"
	"github.com/tushariitr-19/immigration-mcp/util"
)

type GetVisaBulletinInput struct {
	Month string `json:"month,omitempty" jsonschema:"the month to fetch e.g. 'may', defaults to current month"`
	Year  string `json:"year,omitempty" jsonschema:"the year to fetch e.g. '2026', defaults to current year"`
}

var GetVisaBulletinTool = &mcp.Tool{
	Name:        "get_visa_bulletin",
	Description: "Fetch the latest US Visa Bulletin with employment-based priority dates by country and category",
}

func GetVisaBulletinHandler() func(context.Context, *mcp.CallToolRequest, GetVisaBulletinInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input GetVisaBulletinInput) (*mcp.CallToolResult, any, error) {
		logger.Log.Info("fetching visa bulletin", zap.String("month", input.Month), zap.String("year", input.Year))

		bulletin, err := fetchVisaBulletin(input.Month, input.Year)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to fetch visa bulletin: %w", err)
		}

		out, err := json.Marshal(bulletin)
		if err != nil {
			return nil, nil, err
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(out)},
			},
		}, nil, nil
	}
}

func fetchVisaBulletin(month, year string) (*models.VisaBulletin, error) {
	now := time.Now()
	if month == "" {
		month = strings.ToLower(now.Format("January"))
	}
	if year == "" {
		year = now.Format("2006")
	}
	url := util.BuildBulletinURL(month, year)
	logger.Log.Info("fetching from URL", zap.String("url", url))

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	bulletin, err := parseVisaBulletin(string(body), month, year)
	if err != nil {
		return nil, err
	}
	bulletin.PublishedURL = url
	return bulletin, nil
}

func parseVisaBulletin(body, month, year string) (*models.VisaBulletin, error) {
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Extract all tables as 2D slices preserving row structure
	var tables [][][]string
	var extractTables func(*html.Node)
	extractTables = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			var rows [][]string
			var extractRows func(*html.Node)
			extractRows = func(n *html.Node) {
				if n.Type == html.ElementNode && n.Data == "tr" {
					var cells []string
					var extractCells func(*html.Node)
					extractCells = func(n *html.Node) {
						if n.Type == html.ElementNode && (n.Data == "td" || n.Data == "th") {
							var text strings.Builder
							var extractText func(*html.Node)
							extractText = func(n *html.Node) {
								if n.Type == html.TextNode {
									text.WriteString(strings.TrimSpace(n.Data))
								}
								for c := n.FirstChild; c != nil; c = c.NextSibling {
									extractText(c)
								}
							}
							extractText(n)
							cells = append(cells, strings.TrimSpace(text.String()))
							return
						}
						for c := n.FirstChild; c != nil; c = c.NextSibling {
							extractCells(c)
						}
					}
					extractCells(n)
					if len(cells) > 0 {
						rows = append(rows, cells)
					}
					return
				}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					extractRows(c)
				}
			}
			extractRows(n)
			if len(rows) > 0 {
				tables = append(tables, rows)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractTables(c)
		}
	}
	extractTables(doc)

	eb := make(map[string]models.CategoryDates)

	for _, table := range tables {
		if !util.IsEmploymentBasedTable(table) {
			continue
		}

		if len(table) < 2 {
			continue
		}

		// Map country names to column indices from header row
		header := table[0]
		colIdx := make(map[string]int)
		for i, cell := range header {
			if strings.Contains(strings.ToLower(cell), "chargeability") {
				colIdx["Worldwide"] = i
			} else {
				normalized := util.NormalizeCountry(cell)
				if normalized != "Worldwide" {
					colIdx[normalized] = i
				}
			}
		}

		logger.Log.Info("EB table columns", zap.Any("colIdx", colIdx))

		// Parse EB1 EB2 EB3 rows
		ebData := make(map[string][]string)
		for _, row := range table[1:] {
			if len(row) == 0 {
				continue
			}
			label := strings.ToLower(strings.TrimSpace(row[0]))
			var ebNum int
			switch {
			case strings.HasPrefix(label, "1"):
				ebNum = 1
			case strings.HasPrefix(label, "2"):
				ebNum = 2
			case strings.HasPrefix(label, "3"):
				ebNum = 3
			default:
				continue
			}

			for country, idx := range colIdx {
				if idx >= len(row) {
					continue
				}
				if ebData[country] == nil {
					ebData[country] = make([]string, 4)
				}
				ebData[country][ebNum] = util.ParseDate(row[idx])
			}
		}

		for country, dates := range ebData {
			eb[country] = models.CategoryDates{
				EB1: dates[1],
				EB2: dates[2],
				EB3: dates[3],
			}
		}
		break
	}

	if year == "" {
		year = time.Now().Format("2006")
	}
	if month == "" {
		month = strings.ToLower(time.Now().Format("January"))
	}

	return &models.VisaBulletin{
		Month:           month,
		Year:            year,
		EmploymentBased: eb,
	}, nil
}
