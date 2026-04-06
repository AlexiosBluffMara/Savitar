package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

// HTTPClient communicates with a remote MCP server over HTTP using JSON-RPC.
type HTTPClient struct {
	cfg        VSCodeServerConfig
	httpClient *http.Client
	token      string
	idSeq      atomic.Int64
}

// NewHTTPClient creates an HTTP MCP client. authEnv is the name of an
// environment variable that holds a Bearer token, or "" for unauthenticated.
func NewHTTPClient(cfg VSCodeServerConfig, authEnv string) *HTTPClient {
	token := ""
	if authEnv != "" {
		token = strings.TrimSpace(os.Getenv(authEnv))
	}
	return &HTTPClient{
		cfg:        cfg,
		token:      token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *HTTPClient) nextID() int {
	return int(c.idSeq.Add(1))
}

func (c *HTTPClient) post(ctx context.Context, req rpcRequest) (*rpcResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.URL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", "savitar/0.1")
	if c.token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("HTTP 401: auth required — set the auth token env var for %q", c.cfg.Name)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, c.cfg.URL)
	}

	var rpcResp rpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &rpcResp, nil
}

// Initialize performs the MCP initialize handshake.
func (c *HTTPClient) Initialize(ctx context.Context) (ServerInfo, error) {
	resp, err := c.post(ctx, newInitializeRequest(c.nextID()))
	if err != nil {
		return ServerInfo{}, fmt.Errorf("initialize %q: %w", c.cfg.Name, err)
	}
	if resp.Error != nil {
		return ServerInfo{}, fmt.Errorf("initialize %q: %s", c.cfg.Name, resp.Error.Message)
	}
	if resp.Result == nil {
		return ServerInfo{}, fmt.Errorf("initialize %q: empty result", c.cfg.Name)
	}
	return resp.Result.ServerInfo, nil
}

// ListTools returns the tools the server advertises.
func (c *HTTPClient) ListTools(ctx context.Context) ([]Tool, error) {
	resp, err := c.post(ctx, newListToolsRequest(c.nextID()))
	if err != nil {
		return nil, fmt.Errorf("list tools %q: %w", c.cfg.Name, err)
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("list tools %q: %s", c.cfg.Name, resp.Error.Message)
	}
	if resp.Result == nil {
		return nil, nil
	}
	return resp.Result.Tools, nil
}

// CallTool invokes a named tool.
func (c *HTTPClient) CallTool(ctx context.Context, name string, args map[string]any) (ToolCallResult, error) {
	resp, err := c.post(ctx, newCallToolRequest(c.nextID(), name, args))
	if err != nil {
		return ToolCallResult{ServerName: c.cfg.Name, ToolName: name, Err: err}, err
	}
	if resp.Error != nil {
		e := fmt.Errorf("%s", resp.Error.Message)
		return ToolCallResult{ServerName: c.cfg.Name, ToolName: name, Err: e}, e
	}
	if resp.Result == nil {
		return ToolCallResult{ServerName: c.cfg.Name, ToolName: name}, nil
	}
	return ToolCallResult{
		ServerName: c.cfg.Name,
		ToolName:   name,
		Content:    resp.Result.Content,
		IsError:    resp.Result.IsError,
	}, nil
}

// Close is a no-op for HTTP clients.
func (c *HTTPClient) Close() error { return nil }
