package mcp

import "context"

// Client is the interface for a connected MCP server session.
type Client interface {
	// Initialize starts the MCP handshake and returns server info.
	Initialize(ctx context.Context) (ServerInfo, error)
	// ListTools returns the tools advertised by the server.
	ListTools(ctx context.Context) ([]Tool, error)
	// CallTool invokes a named tool with the given arguments.
	CallTool(ctx context.Context, name string, args map[string]any) (ToolCallResult, error)
	// Close releases the server connection.
	Close() error
}
