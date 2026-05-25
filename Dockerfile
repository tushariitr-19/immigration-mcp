FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o immigration-mcp-server ./cmd/server/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/immigration-mcp-server .
ENTRYPOINT ["./immigration-mcp-server"]