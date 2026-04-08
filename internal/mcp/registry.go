package mcp

import (
	"context"
	"fmt"
	"time"

	"github.com/alexiosbluffmara/savitar/internal/config"
)

// ServerStatus is the probed state of a single MCP server.
type ServerStatus struct {
	Name      string
	Mode      string // "http" or "stdio"
	Enabled   bool
	Reachable bool
	Tools     []Tool
	Info      ServerInfo
	Err       string
}

// Registry loads MCP server configs and can probe each server for status.
type Registry struct {
	rootDir    string
	savitarCfg []config.MCPServerConfig
	vscode     []VSCodeServerConfig
}

// NewRegistry creates a Registry. Call Load to resolve server configs.
func NewRegistry(rootDir string, savitarCfg []config.MCPServerConfig) *Registry {
	return &Registry{rootDir: rootDir, savitarCfg: savitarCfg}
}

// Load reads the .vscode/mcp.json file and merges it with the savitar config.
// It is safe to call more than once; the result replaces any prior VSCode data.
func (r *Registry) Load() error {
	vscode, err := LoadVSCodeConfig(r.rootDir)
	if err != nil {
		return err
	}
	r.vscode = vscode
	return nil
}

// resolved returns the merged server list. VSCode entries are the source of
// truth for connection details; the savitar config provides the Enabled flag.
func (r *Registry) resolved() []mergedServer {
	enabledByName := map[string]bool{}
	for _, s := range r.savitarCfg {
		enabledByName[s.Name] = s.Enabled
	}

	out := make([]mergedServer, 0, len(r.vscode))
	for _, v := range r.vscode {
		enabled, ok := enabledByName[v.Name]
		if !ok {
			enabled = true // if not in savitar config, treat as enabled
		}
		out = append(out, mergedServer{vscode: v, enabled: enabled})
	}

	// Include any savitar-only entries (not in .vscode/mcp.json) so they
	// appear in status even if unconfigured.
	vscodeSeen := map[string]bool{}
	for _, v := range r.vscode {
		vscodeSeen[v.Name] = true
	}
	for _, s := range r.savitarCfg {
		if !vscodeSeen[s.Name] {
			out = append(out, mergedServer{
				vscode:  VSCodeServerConfig{Name: s.Name, Type: s.Mode},
				enabled: s.Enabled,
			})
		}
	}
	return out
}

type mergedServer struct {
	vscode  VSCodeServerConfig
	enabled bool
}

func (m mergedServer) mode() string {
	if m.vscode.IsHTTP() {
		return "http"
	}
	if m.vscode.IsStdio() {
		return "stdio"
	}
	return "unknown"
}

// authEnvFor returns the env-var name that holds the auth token for well-known
// HTTP servers. This is a best-effort heuristic; it can be extended over time.
func authEnvFor(name string) string {
	switch name {
	case "github":
		return "GH_TOKEN"
	}
	return ""
}

// ProbeAll probes every enabled server and returns a status list.
// timeout controls how long to wait per server.
func (r *Registry) ProbeAll(timeout time.Duration) []ServerStatus {
	servers := r.resolved()
	statuses := make([]ServerStatus, 0, len(servers))
	for _, s := range servers {
		if !s.enabled {
			statuses = append(statuses, ServerStatus{
				Name:    s.vscode.Name,
				Mode:    s.mode(),
				Enabled: false,
			})
			continue
		}
		status := r.probe(s, timeout)
		statuses = append(statuses, status)
	}
	return statuses
}

func (r *Registry) probe(s mergedServer, timeout time.Duration) ServerStatus {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	st := ServerStatus{
		Name:    s.vscode.Name,
		Mode:    s.mode(),
		Enabled: s.enabled,
	}

	var client Client
	switch {
	case s.vscode.IsHTTP():
		client = NewHTTPClient(s.vscode, authEnvFor(s.vscode.Name))
	case s.vscode.IsStdio():
		client = NewStdioClient(s.vscode)
	default:
		st.Err = "no connection details (not in .vscode/mcp.json or config)"
		return st
	}
	defer func() { _ = client.Close() }()

	info, err := client.Initialize(ctx)
	if err != nil {
		st.Err = fmt.Sprintf("initialize: %v", err)
		return st
	}
	st.Reachable = true
	st.Info = info

	tools, err := client.ListTools(ctx)
	if err != nil {
		st.Err = fmt.Sprintf("list tools: %v", err)
		return st
	}
	st.Tools = tools
	return st
}

// CallTool opens a fresh client for the named server, initializes it, calls
// the tool, and closes the client. This is for single one-shot calls from the
// reply engine; a long-running bot should keep clients alive between calls.
func (r *Registry) CallTool(ctx context.Context, serverName, toolName string, args map[string]any) (ToolCallResult, error) {
	servers := r.resolved()
	for _, s := range servers {
		if s.vscode.Name != serverName || !s.enabled {
			continue
		}

		var client Client
		switch {
		case s.vscode.IsHTTP():
			client = NewHTTPClient(s.vscode, authEnvFor(s.vscode.Name))
		case s.vscode.IsStdio():
			client = NewStdioClient(s.vscode)
		default:
			return ToolCallResult{}, fmt.Errorf("no connection details for MCP server %q", serverName)
		}
		defer func() { _ = client.Close() }()

		if _, err := client.Initialize(ctx); err != nil {
			return ToolCallResult{}, fmt.Errorf("initialize %q: %w", serverName, err)
		}
		return client.CallTool(ctx, toolName, args)
	}
	return ToolCallResult{}, fmt.Errorf("MCP server %q not found or not enabled", serverName)
}
