package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// VSCodeServerConfig is the raw config for a single MCP server from .vscode/mcp.json.
type VSCodeServerConfig struct {
	Name    string            // populated from the map key
	Type    string            `json:"type"`    // "http" or "" (stdio)
	URL     string            `json:"url"`     // for HTTP servers
	Command string            `json:"command"` // for stdio servers
	Args    []string          `json:"args"`    // for stdio servers
	Env     map[string]string `json:"env"`     // for stdio servers
}

// IsHTTP returns true when this server uses HTTP transport.
func (c VSCodeServerConfig) IsHTTP() bool {
	return c.Type == "http" && c.URL != ""
}

// IsStdio returns true when this server uses stdio transport.
func (c VSCodeServerConfig) IsStdio() bool {
	return c.Command != ""
}

// vscodeMCPFile is the on-disk shape of .vscode/mcp.json.
type vscodeMCPFile struct {
	Servers map[string]VSCodeServerConfig `json:"servers"`
}

// LoadVSCodeConfig reads the .vscode/mcp.json file relative to rootDir and
// returns the server configs in declaration order (map iteration order is
// unstable, but names are preserved). Returns nil, nil if the file does not
// exist.
func LoadVSCodeConfig(rootDir string) ([]VSCodeServerConfig, error) {
	path := filepath.Join(rootDir, ".vscode", "mcp.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var f vscodeMCPFile
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, err
	}

	servers := make([]VSCodeServerConfig, 0, len(f.Servers))
	for name, cfg := range f.Servers {
		cfg.Name = name
		servers = append(servers, cfg)
	}
	return servers, nil
}
