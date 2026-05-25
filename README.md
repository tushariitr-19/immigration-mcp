# immigration-mcp

⚠️ **Disclaimer:** This tool is for informational purposes only and does not constitute legal advice. Always consult a qualified immigration attorney for decisions about your specific case. Data is sourced from official government websites (USCIS, State Department) but may not reflect the most recent updates.

An MCP (Model Context Protocol) server for US immigration guidance — live Visa Bulletin, priority date checker, USCIS news, and immigration term explanations.

Built with the official [Go MCP SDK](https://github.com/modelcontextprotocol/go-sdk).

## Why immigration-mcp?

Navigating US immigration is complex and expensive. immigration-mcp gives AI agents access to live, structured immigration data from official government sources — so you can ask questions in plain English and get accurate, up-to-date answers.

- **Free** — powered by public government data sources (USCIS, State Department)
- **Live data** — fetches the latest Visa Bulletin monthly, USCIS news daily
- **Open source** — MIT licensed
- **Production grade** — structured logging, graceful shutdown, unit tested

## Available Tools

| Tool | Description |
|------|-------------|
| `get_visa_bulletin` | Fetch the latest US Visa Bulletin with employment-based priority dates by country and category |
| `check_priority_date` | Check if your priority date is current for I-485 filing |
| `explain_term` | Plain English explanation of any immigration term |

## Example Prompts

Once connected to Claude Desktop:

- *"What are the current EB2 priority dates for India?"*
- *"My priority date is March 2015, I'm from India EB2 — can I file I-485 this month?"*
- *"Has EB2 India moved forward compared to last month?"*
- *"What is the difference between Final Action Date and Date for Filing?"*
- *"Are there any recent USCIS policy changes affecting H1B holders?"*
- *"What documents do I need to file I-485?"*

## Prerequisites

- Go 1.25+
- No API keys required — powered by public government data

## Setup

### 1. Install

```bash
git clone https://github.com/tushariitr-19/immigration-mcp
cd immigration-mcp
go build -o immigration-mcp-server ./cmd/server/
```

### 2. Configure Claude Desktop

Add to your Claude Desktop config (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "immigration-mcp": {
      "command": "/path/to/immigration-mcp-server",
      "env": {
        "DEBUG": "false"
      }
    }
  }
}
```

Optional:
```bash
export DEBUG=true   # enables debug logging
```

## Architecture

```
immigration-mcp/
├── cmd/server/main.go        ← entry point, env vars, graceful shutdown
├── server/server.go          ← MCP server setup, tool registration
├── tools/
│   ├── visa_bulletin.go      ← get_visa_bulletin tool
│   ├── priority_date.go      ← check_priority_date tool
│   └── explain_term.go       ← explain_term tool
├── util/
│   ├── util.go               ← shared helper functions
│   └── constants.go          ← shared constants
├── models/
│   └── models.go             ← shared data models
├── tests/
│   └── visa_bulletin_test.go ← unit tests
├── logger/
│   └── logger.go             ← structured logging via zap
└── Dockerfile
```

Each tool is self-contained — the server is agnostic of what tools do internally. Adding a new tool is a single line in `server/server.go`.

## Running Tests

```bash
# Unit tests only
make test-unit

# Integration tests only
make test-integration

# All tests
make test

# Build binary
make build
```

## Contributing

PRs welcome. To add a new tool:

1. Create `tools/<toolname>.go`
2. Define your input struct and tool definition
3. Register it in `server/server.go` with one line
4. Add unit tests in `tests/<toolname>_test.go`

## License

MIT