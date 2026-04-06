package mcp

// JSON-RPC 2.0 wire types and MCP protocol types.

type rpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id,omitempty"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Result  *rpcResult      `json:"result,omitempty"`
	Error   *rpcErrorObject `json:"error,omitempty"`
}

type rpcResult struct {
	// initialize result
	ProtocolVersion string     `json:"protocolVersion,omitempty"`
	ServerInfo      ServerInfo `json:"serverInfo,omitempty"`
	Capabilities    any        `json:"capabilities,omitempty"`

	// tools/list result
	Tools []Tool `json:"tools,omitempty"`

	// tools/call result
	Content []ToolContent `json:"content,omitempty"`
	IsError bool          `json:"isError,omitempty"`
}

type rpcErrorObject struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ServerInfo is metadata returned by the server on initialize.
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Tool is a capability advertised by an MCP server.
type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	InputSchema any    `json:"inputSchema,omitempty"`
}

// ToolContent is a piece of content returned by a tool call.
type ToolContent struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// ToolResult aggregates the content from a tool call into a single string.
func (t ToolContent) String() string {
	return t.Text
}

// ToolCallResult is the result of calling a tool.
type ToolCallResult struct {
	ServerName string
	ToolName   string
	Content    []ToolContent
	IsError    bool
	Err        error
}

// Text returns the concatenated text from all content items.
func (r ToolCallResult) Text() string {
	out := make([]string, 0, len(r.Content))
	for _, c := range r.Content {
		if c.Text != "" {
			out = append(out, c.Text)
		}
	}
	result := ""
	for i, s := range out {
		if i > 0 {
			result += "\n"
		}
		result += s
	}
	return result
}

// initializeParams is the parameters for the initialize request.
type initializeParams struct {
	ProtocolVersion string         `json:"protocolVersion"`
	Capabilities    map[string]any `json:"capabilities"`
	ClientInfo      clientInfo     `json:"clientInfo"`
}

type clientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

const protocolVersion = "2024-11-05"

func newInitializeRequest(id int) rpcRequest {
	return rpcRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  "initialize",
		Params: initializeParams{
			ProtocolVersion: protocolVersion,
			Capabilities:    map[string]any{},
			ClientInfo:      clientInfo{Name: "savitar", Version: "0.1"},
		},
	}
}

func newInitializedNotification() rpcRequest {
	return rpcRequest{
		JSONRPC: "2.0",
		Method:  "notifications/initialized",
	}
}

func newListToolsRequest(id int) rpcRequest {
	return rpcRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  "tools/list",
		Params:  map[string]any{},
	}
}

func newCallToolRequest(id int, name string, args map[string]any) rpcRequest {
	if args == nil {
		args = map[string]any{}
	}
	return rpcRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  "tools/call",
		Params: map[string]any{
			"name":      name,
			"arguments": args,
		},
	}
}
