package handlers

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Handler defines the interface for registering MCP tools
type Handler interface {
	RegisterTools(mcpServer *server.MCPServer, readOnly bool) error
}

// ToolDefinition defines a single MCP tool with its metadata and handler
type ToolDefinition struct {
	Name        string
	Description string
	ReadOnly    bool // If true, tool is available even in read-only mode
	BuildTool   func() mcp.Tool
	Handler     func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}
