//go:build integration

package tests

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/tushariitr-19/immigration-mcp/models"
)

func newSession(t *testing.T) (*mcp.ClientSession, func()) {
	t.Helper()
	serverBin := os.Getenv("IMMIGRATION_MCP_SERVER")
	if serverBin == "" {
		serverBin = "./immigration-mcp-server"
	}

	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "v1.0.0"}, nil)
	cmd := exec.Command(serverBin)
	cmd.Env = os.Environ()

	transport := &mcp.CommandTransport{Command: cmd}
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}

	return session, func() { session.Close() }
}

// get_visa_bulletin tests

func TestGetVisaBulletinIntegration(t *testing.T) {
	session, cleanup := newSession(t)
	defer cleanup()

	res, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "get_visa_bulletin",
		Arguments: map[string]any{"month": "may", "year": "2026"},
	})
	if err != nil || res.IsError {
		t.Fatalf("get_visa_bulletin failed: %v", err)
	}

	var bulletin models.VisaBulletin
	if err := json.Unmarshal([]byte(res.Content[0].(*mcp.TextContent).Text), &bulletin); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	t.Logf("bulletin month: %s year: %s countries: %d", bulletin.Month, bulletin.Year, len(bulletin.EmploymentBased))

	if bulletin.Month != "may" {
		t.Errorf("expected month may got %s", bulletin.Month)
	}
	if bulletin.EmploymentBased["India"].EB2 == "" {
		t.Error("expected India EB2 date to be non-empty")
	}
	if bulletin.PublishedURL == "" {
		t.Error("expected published_url to be non-empty")
	}
}

func TestGetVisaBulletinDefaultsIntegration(t *testing.T) {
	session, cleanup := newSession(t)
	defer cleanup()

	// No month/year — should default to current
	res, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "get_visa_bulletin",
		Arguments: map[string]any{},
	})
	if err != nil || res.IsError {
		t.Fatalf("get_visa_bulletin with defaults failed: %v", err)
	}

	var bulletin models.VisaBulletin
	if err := json.Unmarshal([]byte(res.Content[0].(*mcp.TextContent).Text), &bulletin); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	t.Logf("bulletin month: %s year: %s", bulletin.Month, bulletin.Year)

	if bulletin.Month == "" {
		t.Error("expected month to be non-empty")
	}
	if len(bulletin.EmploymentBased) == 0 {
		t.Error("expected employment_based to be non-empty")
	}
}

func TestGetVisaBulletinChinaIntegration(t *testing.T) {
	session, cleanup := newSession(t)
	defer cleanup()

	res, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "get_visa_bulletin",
		Arguments: map[string]any{"month": "may", "year": "2026"},
	})
	if err != nil || res.IsError {
		t.Fatalf("get_visa_bulletin failed: %v", err)
	}

	var bulletin models.VisaBulletin
	if err := json.Unmarshal([]byte(res.Content[0].(*mcp.TextContent).Text), &bulletin); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	t.Logf("China EB1: %s EB2: %s EB3: %s", bulletin.EmploymentBased["China"].EB1, bulletin.EmploymentBased["China"].EB2, bulletin.EmploymentBased["China"].EB3)

	china, ok := bulletin.EmploymentBased["China"]
	if !ok {
		t.Fatal("expected China to be in employment_based")
	}
	if china.EB1 == "" || china.EB2 == "" || china.EB3 == "" {
		t.Error("expected all China EB dates to be non-empty")
	}
}

// check_priority_date tests

func TestCheckPriorityDateIntegration(t *testing.T) {
	session, cleanup := newSession(t)
	defer cleanup()

	res, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "check_priority_date",
		Arguments: map[string]any{
			"country":       "India",
			"category":      "EB2",
			"priority_date": "2015-03-10",
		},
	})
	if err != nil || res.IsError {
		t.Fatalf("check_priority_date failed: %v", err)
	}

	var result models.PriorityDateResult
	if err := json.Unmarshal([]byte(res.Content[0].(*mcp.TextContent).Text), &result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	t.Logf("eligible: %v days_behind: %d message: %s", result.Eligible, result.DaysBehind, result.Message)

	if result.Eligible {
		t.Error("expected not eligible for India EB2 2015-03-10")
	}
	if result.DaysBehind == 0 {
		t.Error("expected days_behind to be non-zero")
	}
}

func TestCheckPriorityDateEligibleIntegration(t *testing.T) {
	session, cleanup := newSession(t)
	defer cleanup()

	// India EB2 cutoff is 2014-07-15 — a date before that should be eligible
	res, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "check_priority_date",
		Arguments: map[string]any{
			"country":       "India",
			"category":      "EB2",
			"priority_date": "2013-01-01",
		},
	})
	if err != nil || res.IsError {
		t.Fatalf("check_priority_date failed: %v", err)
	}

	var result models.PriorityDateResult
	if err := json.Unmarshal([]byte(res.Content[0].(*mcp.TextContent).Text), &result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	t.Logf("eligible: %v message: %s", result.Eligible, result.Message)

	if !result.Eligible {
		t.Errorf("expected eligible for India EB2 2013-01-01, got message: %s", result.Message)
	}
}

func TestCheckPriorityDateCurrentCategoryIntegration(t *testing.T) {
	session, cleanup := newSession(t)
	defer cleanup()

	// Worldwide EB1 is Current
	res, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "check_priority_date",
		Arguments: map[string]any{
			"country":       "Worldwide",
			"category":      "EB1",
			"priority_date": "2020-01-01",
		},
	})
	if err != nil || res.IsError {
		t.Fatalf("check_priority_date failed: %v", err)
	}

	var result models.PriorityDateResult
	if err := json.Unmarshal([]byte(res.Content[0].(*mcp.TextContent).Text), &result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	t.Logf("eligible: %v cutoff: %s message: %s", result.Eligible, result.CutoffDate, result.Message)

	if !result.Eligible {
		t.Error("expected eligible for Worldwide EB1 which is Current")
	}
}

func TestCheckPriorityDateCaseInsensitiveIntegration(t *testing.T) {
	session, cleanup := newSession(t)
	defer cleanup()

	// lowercase inputs
	res, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "check_priority_date",
		Arguments: map[string]any{
			"country":       "india",
			"category":      "eb2",
			"priority_date": "2015-03-10",
		},
	})
	if err != nil || res.IsError {
		t.Fatalf("check_priority_date failed: %v", err)
	}

	var result models.PriorityDateResult
	if err := json.Unmarshal([]byte(res.Content[0].(*mcp.TextContent).Text), &result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	t.Logf("country: %s category: %s eligible: %v", result.Country, result.Category, result.Eligible)

	if result.Country != "India" {
		t.Errorf("expected country India got %s", result.Country)
	}
	if result.Category != "EB2" {
		t.Errorf("expected category EB2 got %s", result.Category)
	}
}

// explain_term tests

func TestExplainTermExactMatchIntegration(t *testing.T) {
	session, cleanup := newSession(t)
	defer cleanup()

	res, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "explain_term",
		Arguments: map[string]any{"term": "priority date"},
	})
	if err != nil || res.IsError {
		t.Fatalf("explain_term failed: %v", err)
	}

	var result models.ImmigrationTerm
	if err := json.Unmarshal([]byte(res.Content[0].(*mcp.TextContent).Text), &result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if result.Term != "Priority Date" {
		t.Errorf("expected Priority Date got %s", result.Term)
	}
	if result.Simple == "" {
		t.Error("expected simple explanation to be non-empty")
	}
	t.Logf("term: %s simple: %s", result.Term, result.Simple)
}

func TestExplainTermFuzzyMatchIntegration(t *testing.T) {
	session, cleanup := newSession(t)
	defer cleanup()

	// fuzzy match — "what is h1b visa" should match "h1b"
	res, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "explain_term",
		Arguments: map[string]any{"term": "what is h1b visa"},
	})
	if err != nil || res.IsError {
		t.Fatalf("explain_term failed: %v", err)
	}

	var result models.ImmigrationTerm
	if err := json.Unmarshal([]byte(res.Content[0].(*mcp.TextContent).Text), &result); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if result.Term == "" {
		t.Error("expected a result for fuzzy match")
	}
	t.Logf("term: %s simple: %s", result.Term, result.Simple)
}

func TestExplainTermNotFoundIntegration(t *testing.T) {
	session, cleanup := newSession(t)
	defer cleanup()

	res, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "explain_term",
		Arguments: map[string]any{"term": "nonexistent term xyz"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should return an error result not a panic
	if !res.IsError {
		t.Error("expected IsError to be true for unknown term")
	}
	t.Logf("error response received as expected")
}

func TestExplainTermCaseInsensitiveIntegration(t *testing.T) {
	session, cleanup := newSession(t)
	defer cleanup()

	tests := []string{"PRIORITY DATE", "Priority Date", "priority date", "H1B", "h1b"}

	for _, term := range tests {
		t.Run(term, func(t *testing.T) {
			res, err := session.CallTool(context.Background(), &mcp.CallToolParams{
				Name:      "explain_term",
				Arguments: map[string]any{"term": term},
			})
			if err != nil || res.IsError {
				t.Fatalf("explain_term failed for %s: %v", term, err)
			}

			var result models.ImmigrationTerm
			if err := json.Unmarshal([]byte(res.Content[0].(*mcp.TextContent).Text), &result); err != nil {
				t.Fatalf("failed to parse response: %v", err)
			}

			if result.Term == "" {
				t.Errorf("expected result for term %s", term)
			}
			t.Logf("term: %s → %s", term, result.Term)
		})
	}
}
