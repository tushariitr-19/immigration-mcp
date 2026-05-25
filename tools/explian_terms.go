package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"

	"github.com/tushariitr-19/immigration-mcp/logger"
	"github.com/tushariitr-19/immigration-mcp/models"
	"github.com/tushariitr-19/immigration-mcp/util"
)

type ExplainTermInput struct {
	Term string `json:"term" jsonschema:"the immigration term to explain e.g. 'priority date', 'H1B', 'EAD', 'AC21'"`
}

var ExplainTermTool = &mcp.Tool{
	Name:        "explain_term",
	Description: "Get a plain English explanation of any US immigration term, form, or concept",
}

func ExplainTermHandler() func(context.Context, *mcp.CallToolRequest, ExplainTermInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input ExplainTermInput) (*mcp.CallToolResult, any, error) {
		logger.Log.Info("explaining term", zap.String("term", input.Term))

		result, err := explainTerm(input.Term)
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

func explainTerm(term string) (*models.ImmigrationTerm, error) {
	normalized := strings.ToLower(strings.TrimSpace(term))

	// Exact match first
	if t, ok := util.ImmigrationTerms[normalized]; ok {
		return &t, nil
	}

	// Fuzzy match
	for key, t := range util.ImmigrationTerms {
		if strings.Contains(normalized, key) || strings.Contains(key, normalized) {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("term not found: '%s' — try terms like 'priority date', 'H1B', 'EAD', 'I-485', 'AC21', 'PERM'", term)
}
