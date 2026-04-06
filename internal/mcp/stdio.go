package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"sync/atomic"
)

// StdioClient communicates with a local MCP server process over stdin/stdout
// using newline-delimited JSON-RPC 2.0.
type StdioClient struct {
	cfg     VSCodeServerConfig
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	scanner *bufio.Scanner
	mu      sync.Mutex
	idSeq   atomic.Int64
	started bool
}

// NewStdioClient creates a new StdioClient for the given server config.
// Call Initialize to start the process and perform the MCP handshake.
func NewStdioClient(cfg VSCodeServerConfig) *StdioClient {
	return &StdioClient{cfg: cfg}
}

func (c *StdioClient) start(ctx context.Context) error {
	if c.started {
		return nil
	}
	if c.cfg.Command == "" {
		return fmt.Errorf("stdio client %q: no command configured", c.cfg.Name)
	}

	cmd := exec.CommandContext(ctx, c.cfg.Command, c.cfg.Args...)
	for k, v := range c.cfg.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdio client %q: stdin pipe: %w", c.cfg.Name, err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdio client %q: stdout pipe: %w", c.cfg.Name, err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("stdio client %q: start: %w", c.cfg.Name, err)
	}

	c.cmd = cmd
	c.stdin = stdin
	c.scanner = bufio.NewScanner(stdout)
	c.scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	c.started = true
	return nil
}

func (c *StdioClient) nextID() int {
	return int(c.idSeq.Add(1))
}

func (c *StdioClient) send(req rpcRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(c.stdin, "%s\n", data)
	return err
}

func (c *StdioClient) recv() (*rpcResponse, error) {
	if !c.scanner.Scan() {
		if err := c.scanner.Err(); err != nil {
			return nil, err
		}
		return nil, io.EOF
	}
	line := c.scanner.Bytes()
	var resp rpcResponse
	if err := json.Unmarshal(line, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal rpc response: %w (raw: %s)", err, string(line))
	}
	return &resp, nil
}

// sendAndReceive is the core request/response cycle.
// It skips notifications (responses without an ID matching ours) while waiting.
func (c *StdioClient) sendAndReceive(req rpcRequest) (*rpcResponse, error) {
	if err := c.send(req); err != nil {
		return nil, err
	}
	for {
		resp, err := c.recv()
		if err != nil {
			return nil, err
		}
		// Notifications have no ID; skip them.
		if resp.ID == nil {
			continue
		}
		return resp, nil
	}
}

// Initialize starts the process, runs the MCP handshake, and returns server info.
func (c *StdioClient) Initialize(ctx context.Context) (ServerInfo, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.start(ctx); err != nil {
		return ServerInfo{}, err
	}

	resp, err := c.sendAndReceive(newInitializeRequest(c.nextID()))
	if err != nil {
		return ServerInfo{}, fmt.Errorf("initialize: %w", err)
	}
	if resp.Error != nil {
		return ServerInfo{}, fmt.Errorf("initialize: %s", resp.Error.Message)
	}
	if resp.Result == nil {
		return ServerInfo{}, fmt.Errorf("initialize: empty result")
	}

	// Send the initialized notification (fire-and-forget).
	_ = c.send(newInitializedNotification())

	return resp.Result.ServerInfo, nil
}

// ListTools returns the tools the server advertises.
func (c *StdioClient) ListTools(ctx context.Context) ([]Tool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	resp, err := c.sendAndReceive(newListToolsRequest(c.nextID()))
	if err != nil {
		return nil, fmt.Errorf("list tools: %w", err)
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("list tools: %s", resp.Error.Message)
	}
	if resp.Result == nil {
		return nil, nil
	}
	return resp.Result.Tools, nil
}

// CallTool invokes a named tool and returns its result.
func (c *StdioClient) CallTool(ctx context.Context, name string, args map[string]any) (ToolCallResult, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	resp, err := c.sendAndReceive(newCallToolRequest(c.nextID(), name, args))
	if err != nil {
		return ToolCallResult{ServerName: c.cfg.Name, ToolName: name}, fmt.Errorf("call tool %q: %w", name, err)
	}
	if resp.Error != nil {
		return ToolCallResult{ServerName: c.cfg.Name, ToolName: name}, fmt.Errorf("call tool %q: %s", name, resp.Error.Message)
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

// Close kills the server process.
func (c *StdioClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.started {
		return nil
	}
	_ = c.stdin.Close()
	if c.cmd != nil && c.cmd.Process != nil {
		_ = c.cmd.Process.Kill()
		_ = c.cmd.Wait()
	}
	c.started = false
	return nil
}
