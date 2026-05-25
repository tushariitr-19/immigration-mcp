.PHONY: build test test-unit test-integration clean

# Build the server binary
build:
	go build -o immigration-mcp-server ./cmd/server/

# Run unit tests only
test-unit:
	go test ./tests/... -v

# Run integration tests only (requires built binary)
test-integration: build
	IMMIGRATION_MCP_SERVER=$(PWD)/immigration-mcp-server go test ./tests/... -v -tags integration -run Integration

# Run all tests
test: test-unit test-integration

# Clean built binaries
clean:
	rm -f immigration-mcp-server