package server

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"

	"github.com/tushariitr-19/immigration-mcp/logger"
	"github.com/tushariitr-19/immigration-mcp/tools"
)

const version = "v0.1.0"

type Server struct {
	mcp *mcp.Server
}

func New() *Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "immigration-mcp",
		Version: version,
	}, nil)

	mcp.AddTool(s, tools.GetVisaBulletinTool, tools.GetVisaBulletinHandler())
	mcp.AddTool(s, tools.CheckPriorityDateTool, tools.CheckPriorityDateHandler())
	mcp.AddTool(s, tools.ExplainTermTool, tools.ExplainTermHandler())

	logger.Log.Info("registered tools", zap.String("version", version))

	return &Server{mcp: s}
}

func (s *Server) Run(ctx context.Context) error {
	return s.mcp.Run(ctx, &mcp.StdioTransport{})
}
