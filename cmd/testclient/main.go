package main

import (
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "v1.0.0"}, nil)

	cmd := exec.Command("go", "run", "./cmd/server/")
	cmd.Env = os.Environ()

	transport := &mcp.CommandTransport{Command: cmd}
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// List tools
	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		log.Fatalf("ListTools failed: %v", err)
	}
	for _, t := range tools.Tools {
		log.Printf("tool: %s — %s", t.Name, t.Description)
	}

	// Test get_visa_bulletin
	res, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "get_visa_bulletin",
		Arguments: map[string]any{"month": "may", "year": "2026"},
	})
	if err != nil {
		log.Fatalf("get_visa_bulletin failed: %v", err)
	}

	for _, c := range res.Content {
		log.Print(c.(*mcp.TextContent).Text)
	}

	// Test check_priority_date
	res2, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "check_priority_date",
		Arguments: map[string]any{
			"country":       "India",
			"category":      "EB2",
			"priority_date": "2015-03-10",
		},
	})
	if err != nil {
		log.Fatalf("check_priority_date failed: %v", err)
	}
	for _, c := range res2.Content {
		log.Print(c.(*mcp.TextContent).Text)
	}
}
