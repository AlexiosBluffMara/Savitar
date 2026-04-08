package runtime

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	stdRuntime "runtime"
	"time"

	"github.com/alexiosbluffmara/savitar/internal/config"
	"github.com/alexiosbluffmara/savitar/internal/contracts"
	"github.com/alexiosbluffmara/savitar/internal/customization"
	"github.com/alexiosbluffmara/savitar/internal/gateway"
	"github.com/alexiosbluffmara/savitar/internal/mcp"
	"github.com/alexiosbluffmara/savitar/internal/memory"
	"github.com/alexiosbluffmara/savitar/internal/models"
	"github.com/alexiosbluffmara/savitar/internal/persona"
	"github.com/alexiosbluffmara/savitar/internal/repo"
	"github.com/alexiosbluffmara/savitar/internal/session"
)

type Check struct {
	Name      string
	Required  bool
	Available bool
	Details   string
}

type Status struct {
	RootDir            string
	ConfigPath         string
	ConfigLoaded       bool
	AgentName          string
	ModelProfiles      int
	EnabledMCPServers  int
	AgentCount         int
	SkillCount         int
	EnabledSurfaces    int
	TotalSurfaces      int
	SessionPath        string
	SessionInitialized bool
}

type DiscordStatus struct {
	Enabled                   bool
	TokenEnv                  string
	TokenPresent              bool
	DisplayName               string
	OperatorUserCount         int
	RequireMention            bool
	RespondInDirectMessages   bool
	AllowCloudRepliesInGuilds bool
	AllowCloudRepliesInDMs    bool
	AllowLiveWebLookupInGuilds bool
	AllowLiveWebLookupInDMs    bool
	UseMessageContentIntent   bool
	AllowedChannelIDs         []string
	PerUserCooldownSeconds    int
	MaxConcurrentReplies      int
	MaxResponseChars          int
	PresenceText              string
	ImmediateAck              bool
	TriggerMode               string
}

type IntegrationStatus struct {
	Name              string
	Enabled           bool
	AuthSource        string
	CredentialPresent bool
	TokenEnv          string
	CLIPath           string
	Details           string
}

type Runtime struct {
	loaded   config.Loaded
	router   models.Router
	rootDir  string
	mcpReg   *mcp.Registry
	memStore *memory.Store
}

func New(loaded config.Loaded) *Runtime {
	return NewAtRoot(loaded, ".")
}

func NewAtRoot(loaded config.Loaded, rootDir string) *Runtime {
	rt := &Runtime{
		loaded:  loaded,
		router:  models.DefaultRouter(),
		rootDir: rootDir,
	}
	rt.mcpReg = mcp.NewRegistry(rootDir, loaded.Config.MCP.Servers)
	_ = rt.mcpReg.Load() // best-effort; fails gracefully if .vscode/mcp.json absent
	rt.memStore = memory.NewStore(rt.resolvePath(loaded.Config.Memory.KnowledgeDir))
	return rt
}

func (r *Runtime) Config() config.Loaded {
	return r.loaded
}

func (r *Runtime) Router() models.Router {
	return r.router
}

func (r *Runtime) RootDir() string {
	return r.rootDir
}

func (r *Runtime) ModelProfiles() []models.Profile {
	return r.router.Profiles()
}

func (r *Runtime) Contracts() []contracts.Contract {
	return contracts.Default()
}

func (r *Runtime) Persona() persona.Profile {
	return persona.FromConfig(r.loaded.Config.Agent)
}

func (r *Runtime) Agents() ([]customization.Agent, error) {
	return customization.DiscoverAgents(r.rootDir)
}

func (r *Runtime) Skills() ([]customization.Skill, error) {
	return customization.DiscoverSkills(r.rootDir)
}

func (r *Runtime) Session() (session.Report, error) {
	return session.NewStore(r.resolvePath(r.loaded.Config.Memory.SessionIndex)).Load()
}

func (r *Runtime) InitSession() (session.Report, error) {
	state := session.DefaultState()
	state.CurrentModelProfile = r.loaded.Config.Models.LocalDefault.Profile
	state.LastCommand = "session init"
	return session.NewStore(r.resolvePath(r.loaded.Config.Memory.SessionIndex)).Init(state)
}

func (r *Runtime) GatewayPlan() gateway.Plan {
	return gateway.BuildPlan(r.loaded.Config)
}

func (r *Runtime) DiscordStatus() DiscordStatus {
	cfg := r.loaded.Config.Transports.Discord
	triggerMode := "mentions in any joined channel"
	if len(cfg.AllowedChannelIDs) > 0 {
		triggerMode = "mentions in allowlisted channels"
	}
	if !cfg.RequireMention {
		if len(cfg.AllowedChannelIDs) > 0 {
			triggerMode = "all messages in allowlisted channels"
		} else {
			triggerMode = "all guild messages"
		}
	}
	if cfg.RespondInDirectMessages {
		triggerMode += " + direct messages"
	}

	return DiscordStatus{
		Enabled:                   cfg.Enabled,
		TokenEnv:                  cfg.BotTokenEnv,
		TokenPresent:              os.Getenv(cfg.BotTokenEnv) != "",
		DisplayName:               cfg.DisplayName,
		OperatorUserCount:         len(cfg.OperatorUserIDs),
		RequireMention:            cfg.RequireMention,
		RespondInDirectMessages:   cfg.RespondInDirectMessages,
		AllowCloudRepliesInGuilds: cfg.AllowCloudRepliesInGuilds,
		AllowCloudRepliesInDMs:    cfg.AllowCloudRepliesInDMs,
		AllowLiveWebLookupInGuilds: cfg.AllowLiveWebLookupInGuilds,
		AllowLiveWebLookupInDMs:    cfg.AllowLiveWebLookupInDMs,
		UseMessageContentIntent:   cfg.UseMessageContentIntent,
		AllowedChannelIDs:         append([]string(nil), cfg.AllowedChannelIDs...),
		PerUserCooldownSeconds:    cfg.PerUserCooldownSeconds,
		MaxConcurrentReplies:      cfg.MaxConcurrentReplies,
		MaxResponseChars:          cfg.MaxResponseChars,
		PresenceText:              cfg.PresenceText,
		ImmediateAck:              cfg.ImmediateAck,
		TriggerMode:               triggerMode,
	}
}

func (r *Runtime) IntegrationStatuses() []IntegrationStatus {
	cfg := r.loaded.Config.Integrations
	return []IntegrationStatus{
		{
			Name:              "ollama",
			Enabled:           cfg.Ollama.Enabled,
			AuthSource:        "api-key-env",
			CredentialPresent: envPresent(cfg.Ollama.APIKeyEnv),
			TokenEnv:          cfg.Ollama.APIKeyEnv,
			CLIPath:           integrationCommandPath(r.rootDir, "ollama"),
			Details:           fmt.Sprintf("%s model %s", cfg.Ollama.BaseURL, cfg.Ollama.CloudModel),
		},
		{
			Name:              "github",
			Enabled:           cfg.GitHub.Enabled,
			AuthSource:        githubAuthSource(cfg.GitHub),
			CredentialPresent: githubCredentialPresent(cfg.GitHub),
			TokenEnv:          cfg.GitHub.TokenEnv,
			CLIPath:           integrationCommandPath(r.rootDir, "gh"),
			Details:           githubDetails(cfg.GitHub),
		},
		{
			Name:              "huggingface",
			Enabled:           cfg.HuggingFace.Enabled,
			AuthSource:        "token-env",
			CredentialPresent: envPresent(cfg.HuggingFace.TokenEnv),
			TokenEnv:          cfg.HuggingFace.TokenEnv,
			CLIPath:           integrationCommandPath(r.rootDir, "hf"),
			Details:           cfg.HuggingFace.CacheDir,
		},
		{
			Name:              "kaggle",
			Enabled:           cfg.Kaggle.Enabled,
			AuthSource:        "token-env",
			CredentialPresent: envPresent(cfg.Kaggle.TokenEnv),
			TokenEnv:          cfg.Kaggle.TokenEnv,
			CLIPath:           integrationCommandPath(r.rootDir, "kaggle"),
			Details:           cfg.Kaggle.ConfigDir,
		},
	}
}

func (r *Runtime) Status() (Status, error) {
	agents, err := r.Agents()
	if err != nil {
		return Status{}, err
	}

	skills, err := r.Skills()
	if err != nil {
		return Status{}, err
	}

	sessionReport, err := r.Session()
	if err != nil {
		return Status{}, err
	}

	enabledMCPServers := 0
	for _, server := range r.loaded.Config.MCP.Servers {
		if server.Enabled {
			enabledMCPServers++
		}
	}

	plan := r.GatewayPlan()

	return Status{
		RootDir:            r.rootDir,
		ConfigPath:         r.loaded.Path,
		ConfigLoaded:       r.loaded.Exists,
		AgentName:          r.loaded.Config.Agent.Name,
		ModelProfiles:      len(r.ModelProfiles()),
		EnabledMCPServers:  enabledMCPServers,
		AgentCount:         len(agents),
		SkillCount:         len(skills),
		EnabledSurfaces:    plan.EnabledCount(),
		TotalSurfaces:      len(plan.Surfaces),
		SessionPath:        sessionReport.Path,
		SessionInitialized: sessionReport.Exists,
	}, nil
}

func (r *Runtime) Doctor() []Check {
	checks := []Check{
		commandCheck("brew", true),
		commandCheck("go", true),
		commandCheck("node", true),
		commandCheck("npm", true),
		commandCheck("npx", true),
		commandCheck("gh", true),
		commandCheck("jq", false),
		commandCheck("docker", false),
		commandCheck("python3", false),
		commandCheck("uv", false),
		commandCheck("adb", false),
	}

	if stdRuntime.GOOS == "darwin" {
		checks = append(checks, commandCheck("osascript", false))
	}

	discord := r.DiscordStatus()
	if discord.Enabled {
		checks = append(checks, envCheck(discord.TokenEnv, true))
		if !discord.RequireMention && !discord.UseMessageContentIntent {
			checks = append(checks, Check{
				Name:      "discord-message-content-intent",
				Required:  true,
				Available: false,
				Details:   "set transports.discord.useMessageContentIntent=true or keep requireMention=true",
			})
		}
		if !discord.RequireMention && len(discord.AllowedChannelIDs) == 0 {
			checks = append(checks, Check{
				Name:      "discord-channel-allowlist",
				Required:  true,
				Available: false,
				Details:   "set transports.discord.allowedChannelIDs before disabling requireMention",
			})
		}
	}

	for _, status := range r.IntegrationStatuses() {
		if !status.Enabled {
			continue
		}

		switch status.Name {
		case "ollama":
			checks = append(checks, envCheck(status.TokenEnv, true))
			checks = append(checks, commandCheck("ollama", true))
		case "github":
			checks = append(checks, githubCredentialCheck(r.loaded.Config.Integrations.GitHub))
		case "huggingface":
			checks = append(checks, envCheck(status.TokenEnv, true))
			checks = append(checks, integrationCommandCheck(r.rootDir, "hf", true))
		case "kaggle":
			checks = append(checks, envCheck(status.TokenEnv, true))
			checks = append(checks, integrationCommandCheck(r.rootDir, "kaggle", true))
		}
	}

	return checks
}

// MCPStatus probes all configured MCP servers and returns their status.
// timeout controls how long to wait per server.
func (r *Runtime) MCPStatus(timeout time.Duration) []mcp.ServerStatus {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	return r.mcpReg.ProbeAll(timeout)
}

// MCPCallTool invokes a named tool on the given MCP server.
func (r *Runtime) MCPCallTool(ctx context.Context, serverName, toolName string, args map[string]any) (mcp.ToolCallResult, error) {
	return r.mcpReg.CallTool(ctx, serverName, toolName, args)
}

// RepoAnalyze clones and analyzes a repository, returning a structured summary.
func (r *Runtime) RepoAnalyze(ctx context.Context, repoURL string) (*repo.Summary, error) {
	outputDir := r.resolvePath("workspace/repo-analysis")
	analyzer := repo.NewAnalyzer(outputDir)
	return analyzer.Analyze(ctx, repoURL)
}

// MemoryStore returns the memory store for this runtime.
func (r *Runtime) MemoryStore() *memory.Store {
	return r.memStore
}

// MemoryWrite writes a pack to the knowledge store.
func (r *Runtime) MemoryWrite(pack memory.Pack) error {
	return r.memStore.Write(pack)
}

// MemoryBySubject loads all packs for a given subject tag.
func (r *Runtime) MemoryBySubject(subject string) ([]memory.Pack, error) {
	return r.memStore.ReadBySubject(subject)
}

// MemoryList lists all known memory packs.
func (r *Runtime) MemoryList() ([]memory.PackMeta, error) {
	return r.memStore.ListAll()
}

// RepoMarkdownSearch searches markdown files in the local repository and
// returns ranked excerpts plus a lightweight document graph for LLM context.
func (r *Runtime) RepoMarkdownSearch(query string) (memory.MarkdownSearchResult, error) {
	cfg := r.loaded.Config.Knowledge
	return memory.SearchRepoMarkdown(r.rootDir, cfg.RepoMarkdownDirs, query, cfg.MaxRepoResults, cfg.MaxGraphEdges)
}

// SessionList returns summaries for all session files.
func (r *Runtime) SessionList() ([]session.Summary, error) {
	return session.NewStore(r.resolvePath(r.loaded.Config.Memory.SessionIndex)).List()
}

func (r *Runtime) Plan() []string {
	return []string{
		"Lock the product to Discord, curated knowledge, evidence-backed replies, and one authenticated operator console.",
		"Build the orchestration program that turns an inbound envelope into a route, evidence set, reply draft, run record, and delivery plan.",
		"Implement retrieval, citations, session continuity, and durable run records before adding more public surfaces.",
		"Replace the direct Discord reply path with the orchestrator and expose operator review state through backend services.",
		"Ship the authenticated operator web UI, then expand to later transports and automation only after the core workflow is coherent.",
	}
}

func commandCheck(name string, required bool) Check {
	path, err := exec.LookPath(name)
	if err != nil {
		return Check{
			Name:      name,
			Required:  required,
			Available: false,
			Details:   "not found in PATH",
		}
	}

	return Check{
		Name:      name,
		Required:  required,
		Available: true,
		Details:   fmt.Sprintf("found at %s", path),
	}
}

func envCheck(name string, required bool) Check {
	if name == "" {
		return Check{
			Name:      "env",
			Required:  required,
			Available: false,
			Details:   "env var name not configured",
		}
	}

	if os.Getenv(name) == "" {
		return Check{
			Name:      name,
			Required:  required,
			Available: false,
			Details:   "env var not set",
		}
	}

	return Check{
		Name:      name,
		Required:  required,
		Available: true,
		Details:   "env var is set",
	}
}

func envPresent(name string) bool {
	if name == "" {
		return false
	}
	return os.Getenv(name) != ""
}

func commandPath(name string) string {
	prefixes := []string{"/opt/homebrew/bin", "/usr/local/bin", "/usr/bin", "/bin"}
	if homeDir, err := os.UserHomeDir(); err == nil {
		prefixes = append([]string{filepath.Join(homeDir, ".local", "bin")}, prefixes...)
	}

	for _, prefix := range prefixes {
		candidate := filepath.Join(prefix, name)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate
		}
	}

	path, err := exec.LookPath(name)
	if err != nil {
		return ""
	}
	return path
}

func integrationCommandPath(rootDir string, name string) string {
	_ = rootDir
	return commandPath(name)
}

func integrationCommandCheck(rootDir string, name string, required bool) Check {
	path := integrationCommandPath(rootDir, name)
	if path == "" {
		return Check{
			Name:      name,
			Required:  required,
			Available: false,
			Details:   "not found in PATH or workspace .venv/bin",
		}
	}
	return Check{
		Name:      name,
		Required:  required,
		Available: true,
		Details:   fmt.Sprintf("found at %s", path),
	}
}

func githubAuthSource(cfg config.GitHubIntegrationConfig) string {
	if cfg.PreferGHCLIAuth {
		return "gh-cli-or-env"
	}
	return "token-env"
}

func githubDetails(cfg config.GitHubIntegrationConfig) string {
	if cfg.PreferGHCLIAuth {
		return "gh credential store preferred; GH_TOKEN overrides when set"
	}
	return "GH_TOKEN-backed GitHub API access"
}

func githubCredentialPresent(cfg config.GitHubIntegrationConfig) bool {
	if envPresent(cfg.TokenEnv) {
		return true
	}
	if !cfg.PreferGHCLIAuth {
		return false
	}
	return ghAuthReady()
}

func githubCredentialCheck(cfg config.GitHubIntegrationConfig) Check {
	if envPresent(cfg.TokenEnv) {
		return Check{
			Name:      cfg.TokenEnv,
			Required:  true,
			Available: true,
			Details:   "env var is set",
		}
	}
	if !cfg.PreferGHCLIAuth {
		return envCheck(cfg.TokenEnv, true)
	}
	if commandPath("gh") == "" {
		return Check{
			Name:      "github-auth",
			Required:  true,
			Available: false,
			Details:   "gh not found and GH_TOKEN not set",
		}
	}
	if ghAuthReady() {
		return Check{
			Name:      "github-auth",
			Required:  true,
			Available: true,
			Details:   "gh credential store is authenticated",
		}
	}
	return Check{
		Name:      "github-auth",
		Required:  true,
		Available: false,
		Details:   "GH_TOKEN not set and gh auth status is not authenticated",
	}
}

func ghAuthReady() bool {
	path := commandPath("gh")
	if path == "" {
		return false
	}
	cmd := exec.Command(path, "auth", "status")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func (r *Runtime) resolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	return filepath.Clean(filepath.Join(r.rootDir, path))
}
