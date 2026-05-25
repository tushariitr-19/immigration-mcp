package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"

	"github.com/tushariitr-19/immigration-mcp/logger"
	"github.com/tushariitr-19/immigration-mcp/models"
	"github.com/tushariitr-19/immigration-mcp/util"
)

type CheckPriorityDateInput struct {
	Country      string `json:"country" jsonschema:"country of birth e.g. India, China, Mexico, Philippines, Worldwide"`
	Category     string `json:"category" jsonschema:"employment-based category e.g. EB1, EB2, EB3"`
	PriorityDate string `json:"priority_date" jsonschema:"your priority date in YYYY-MM-DD format e.g. 2015-03-10"`
}

var CheckPriorityDateTool = &mcp.Tool{
	Name:        "check_priority_date",
	Description: "Check if your priority date is current for I-485 filing based on the latest Visa Bulletin",
}

func CheckPriorityDateHandler() func(context.Context, *mcp.CallToolRequest, CheckPriorityDateInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input CheckPriorityDateInput) (*mcp.CallToolResult, any, error) {
		logger.Log.Info("checking priority date",
			zap.String("country", input.Country),
			zap.String("category", input.Category),
			zap.String("priority_date", input.PriorityDate),
		)

		// Fetch latest visa bulletin
		bulletin, err := fetchVisaBulletin("", "")
		if err != nil {
			return nil, nil, fmt.Errorf("failed to fetch visa bulletin: %w", err)
		}

		result, err := checkPriorityDate(bulletin, input.Country, input.Category, input.PriorityDate)
		if err != nil {
			return nil, nil, err
		}

		out, err := json.Marshal(result)
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

func checkPriorityDate(bulletin *models.VisaBulletin, country, category, priorityDate string) (*models.PriorityDateResult, error) {
	// Normalize inputs
	country = util.NormalizeInputCountry(country)
	category = util.NormalizeCategory(category)

	// Find country in bulletin
	countryDates, ok := bulletin.EmploymentBased[country]
	if !ok {
		return nil, fmt.Errorf("country not found in bulletin: %s", country)
	}

	// Get cutoff date for category
	var cutoff string
	switch category {
	case "EB1":
		cutoff = countryDates.EB1
	case "EB2":
		cutoff = countryDates.EB2
	case "EB3":
		cutoff = countryDates.EB3
	default:
		return nil, fmt.Errorf("invalid category: %s, must be EB1, EB2 or EB3", category)
	}

	// Handle Current and Unavailable
	if cutoff == "Current" {
		return &models.PriorityDateResult{
			Country:    country,
			Category:   category,
			YourDate:   priorityDate,
			CutoffDate: cutoff,
			Eligible:   true,
			Message:    "Your category is current — you are eligible to file I-485",
		}, nil
	}

	if cutoff == "Unavailable" {
		return &models.PriorityDateResult{
			Country:    country,
			Category:   category,
			YourDate:   priorityDate,
			CutoffDate: cutoff,
			Eligible:   false,
			Message:    "This category is currently unavailable — no visas are being issued",
		}, nil
	}

	// Parse dates
	yourDate, err := time.Parse("2006-01-02", priorityDate)
	if err != nil {
		return nil, fmt.Errorf("invalid priority date format, use YYYY-MM-DD: %w", err)
	}

	cutoffDate, err := time.Parse("2006-01-02", cutoff)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cutoff date: %w", err)
	}

	eligible := yourDate.Before(cutoffDate)

	result := &models.PriorityDateResult{
		Country:    country,
		Category:   category,
		YourDate:   priorityDate,
		CutoffDate: cutoff,
		Eligible:   eligible,
	}

	if eligible {
		result.Message = fmt.Sprintf("Your priority date %s is before the cutoff %s — you are eligible to file I-485", priorityDate, cutoff)
	} else {
		days := int(yourDate.Sub(cutoffDate).Hours() / 24)
		result.DaysBehind = days
		months := days / 30
		remainingDays := days % 30
		result.Message = fmt.Sprintf("Your priority date %s is %d months and %d days behind the cutoff %s — not yet eligible to file I-485", priorityDate, months, remainingDays, cutoff)
	}

	return result, nil
}
