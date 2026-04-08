package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Loaded struct {
	Config Config
	Path   string
	Exists bool
}

type Config struct {
	Agent        AgentConfig        `json:"agent"`
	Models       ModelsConfig       `json:"models"`
	MCP          MCPConfig          `json:"mcp"`
	Integrations IntegrationsConfig `json:"integrations"`
	Transports   TransportsConfig   `json:"transports"`
	WebUI        WebUIConfig        `json:"webui"`
	Automation   AutomationConfig   `json:"automation"`
	Memory       MemoryConfig       `json:"memory"`
	Knowledge    KnowledgeConfig    `json:"knowledge"`
	Operator     OperatorConfig     `json:"operator"`
	Skills       SkillsConfig       `json:"skills"`
}

type AgentConfig struct {
	Name    string             `json:"name"`
	Style   string             `json:"style"`
	Persona AgentPersonaConfig `json:"persona"`
}

type AgentPersonaConfig struct {
	Tone              string `json:"tone"`
	DesignBias        string `json:"designBias"`
	CommentaryDensity string `json:"commentaryDensity"`
	PublicBio         string `json:"publicBio"`
	DisclosurePolicy  string `json:"disclosurePolicy"`
}

type ModelsConfig struct {
	LocalDefault    ModelProfileConfig  `json:"localDefault"`
	CopilotProfiles []CopilotLaneConfig `json:"copilotProfiles"`
}

type ModelProfileConfig struct {
	Provider        string  `json:"provider"`
	Profile         string  `json:"profile"`
	Model           string  `json:"model"`
	Endpoint        string  `json:"endpoint"`
	UsageMultiplier float64 `json:"usageMultiplier"`
}

type CopilotLaneConfig struct {
	Name            string  `json:"name"`
	UsageMultiplier float64 `json:"usageMultiplier"`
	Purpose         string  `json:"purpose"`
}

type MCPConfig struct {
	Servers []MCPServerConfig `json:"servers"`
}

type MCPServerConfig struct {
	Name    string `json:"name"`
	Mode    string `json:"mode"`
	Enabled bool   `json:"enabled"`
}

type IntegrationsConfig struct {
	Ollama      OllamaIntegrationConfig      `json:"ollama"`
	GitHub      GitHubIntegrationConfig      `json:"github"`
	HuggingFace HuggingFaceIntegrationConfig `json:"huggingface"`
	Kaggle      KaggleIntegrationConfig      `json:"kaggle"`
}

type OllamaIntegrationConfig struct {
	Enabled    bool   `json:"enabled"`
	BaseURL    string `json:"baseURL"`
	APIKeyEnv  string `json:"apiKeyEnv"`
	CloudModel string `json:"cloudModel"`
}

type GitHubIntegrationConfig struct {
	Enabled         bool   `json:"enabled"`
	TokenEnv        string `json:"tokenEnv"`
	PreferGHCLIAuth bool   `json:"preferGHCLIAuth"`
}

type HuggingFaceIntegrationConfig struct {
	Enabled  bool   `json:"enabled"`
	TokenEnv string `json:"tokenEnv"`
	CacheDir string `json:"cacheDir"`
}

type KaggleIntegrationConfig struct {
	Enabled   bool   `json:"enabled"`
	TokenEnv  string `json:"tokenEnv"`
	ConfigDir string `json:"configDir"`
}

type TransportsConfig struct {
	Discord  DiscordConfig  `json:"discord"`
	WhatsApp WhatsAppConfig `json:"whatsapp"`
	IMessage IMessageConfig `json:"imessage"`
}

type DiscordConfig struct {
	Enabled                   bool     `json:"enabled"`
	BotTokenEnv               string   `json:"botTokenEnv"`
	DisplayName               string   `json:"displayName"`
	OperatorUserIDs           []string `json:"operatorUserIDs"`
	RequireMention            bool     `json:"requireMention"`
	RespondInDirectMessages   bool     `json:"respondInDirectMessages"`
	AllowCloudRepliesInGuilds bool     `json:"allowCloudRepliesInGuildChannels"`
	AllowCloudRepliesInDMs    bool     `json:"allowCloudRepliesInDirectMessages"`
	AllowLiveWebLookupInGuilds bool    `json:"allowLiveWebLookupInGuildChannels"`
	AllowLiveWebLookupInDMs    bool    `json:"allowLiveWebLookupInDirectMessages"`
	UseMessageContentIntent   bool     `json:"useMessageContentIntent"`
	AllowedChannelIDs         []string `json:"allowedChannelIDs"`
	PerUserCooldownSeconds    int      `json:"perUserCooldownSeconds"`
	MaxConcurrentReplies      int      `json:"maxConcurrentReplies"`
	MaxResponseChars          int      `json:"maxResponseChars"`
	PresenceText              string   `json:"presenceText"`
	ImmediateAck              bool     `json:"immediateAck"`
}

type WhatsAppConfig struct {
	Enabled     bool   `json:"enabled"`
	Bridge      string `json:"bridge"`
	DeviceName  string `json:"deviceName"`
	NumberEnv   string `json:"numberEnv"`
	DisplayName string `json:"displayName"`
}

type IMessageConfig struct {
	Enabled        bool     `json:"enabled"`
	Bridge         string   `json:"bridge"`
	AllowedSenders []string `json:"allowedSenders"`
	AccountEnv     string   `json:"accountEnv"`
	DisplayName    string   `json:"displayName"`
}

type WebUIConfig struct {
	Enabled                 bool     `json:"enabled"`
	PublicBaseURL           string   `json:"publicBaseURL"`
	AuthProvider            string   `json:"authProvider"`
	GoogleClientIDEnv       string   `json:"googleClientIDEnv"`
	GoogleClientSecretEnv   string   `json:"googleClientSecretEnv"`
	SessionSecretEnv        string   `json:"sessionSecretEnv"`
	AllowedEmailDomains     []string `json:"allowedEmailDomains"`
	RequireOperatorApproval bool     `json:"requireOperatorApproval"`
}

type AutomationConfig struct {
	PlaywrightProfileDir string                `json:"playwrightProfileDir"`
	PeekabooEnabled      bool                  `json:"peekabooEnabled"`
	AllowShellExecution  bool                  `json:"allowShellExecution"`
	ScreenCapture        ScreenCaptureConfig   `json:"screenCapture"`
	RemoteHosts          []RemoteHostConfig    `json:"remoteHosts"`
	MobileCompanion      MobileCompanionConfig `json:"mobileCompanion"`
}

type ScreenCaptureConfig struct {
	ScreenshotsDir      string `json:"screenshotsDir"`
	RecordingsDir       string `json:"recordingsDir"`
	AllowDesktopControl bool   `json:"allowDesktopControl"`
}

type RemoteHostConfig struct {
	Name      string `json:"name"`
	Transport string `json:"transport"`
	Enabled   bool   `json:"enabled"`
}

type MobileCompanionConfig struct {
	Enabled    bool   `json:"enabled"`
	Platform   string `json:"platform"`
	DeviceName string `json:"deviceName"`
	Phase      string `json:"phase"`
}

type MemoryConfig struct {
	KnowledgeDir string `json:"knowledgeDir"`
	SnapshotDir  string `json:"snapshotDir"`
	SessionIndex string `json:"sessionIndex"`
}

type KnowledgeConfig struct {
	CatalogDir              string   `json:"catalogDir"`
	IndexDir                string   `json:"indexDir"`
	MaxPacksPerReply        int      `json:"maxPacksPerReply"`
	MaxExcerptsPerReply     int      `json:"maxExcerptsPerReply"`
	RequireSourceMetadata   bool     `json:"requireSourceMetadata"`
	RepoMarkdownDirs        []string `json:"repoMarkdownDirs"`
	MaxRepoResults          int      `json:"maxRepoResults"`
	MaxGraphEdges           int      `json:"maxGraphEdges"`
	EnableLiveWebLookup     bool     `json:"enableLiveWebLookup"`
	LiveWebProvider         string   `json:"liveWebProvider"`
	WebLookupTimeoutSeconds int      `json:"webLookupTimeoutSeconds"`
}

type OperatorConfig struct {
	ReviewQueueDir             string   `json:"reviewQueueDir"`
	RunLogDir                  string   `json:"runLogDir"`
	ApprovedEmails             []string `json:"approvedEmails"`
	RequireReviewForToolErrors bool     `json:"requireReviewForToolErrors"`
}

type SkillsConfig struct {
	WorkspaceDir       string `json:"workspaceDir"`
	UserDir            string `json:"userDir"`
	AllowRemoteCatalog bool   `json:"allowRemoteCatalog"`
}

func DefaultPath() string {
	return "config/savitar.local.json"
}

func Default() Config {
	return Config{
		Agent: AgentConfig{
			Name:  "Savitar",
			Style: "Direct, human, opinionated, memory-aware",
			Persona: AgentPersonaConfig{
				Tone:              "clear and conversational",
				DesignBias:        "prefer readable, pragmatic systems over novelty",
				CommentaryDensity: "medium",
				PublicBio:         "Systems-minded operator with a strong point of view and a calm, human conversational style",
				DisclosurePolicy:  "Speak naturally, but never misrepresent system identity when trust, consent, billing, or safety are involved",
			},
		},
		Models: ModelsConfig{
			LocalDefault: ModelProfileConfig{
				Provider:        "ollama",
				Profile:         "local-default",
				Model:           "gemma4:e4b",
				Endpoint:        "http://127.0.0.1:11434",
				UsageMultiplier: 0,
			},
			CopilotProfiles: []CopilotLaneConfig{
				{Name: "copilot-0x", UsageMultiplier: 0, Purpose: "Routine work with effectively unlimited usage"},
				{Name: "copilot-0.33x", UsageMultiplier: 0.33, Purpose: "Medium-complexity synthesis"},
				{Name: "copilot-1x", UsageMultiplier: 1, Purpose: "High-complexity reasoning and execution"},
			},
		},
		MCP: MCPConfig{
			Servers: []MCPServerConfig{
				{Name: "github", Mode: "remote", Enabled: true},
				{Name: "context7", Mode: "local", Enabled: true},
				{Name: "tavily", Mode: "local", Enabled: true},
				{Name: "playwright", Mode: "local", Enabled: true},
				{Name: "peekaboo", Mode: "local", Enabled: true},
			},
		},
		Integrations: IntegrationsConfig{
			Ollama: OllamaIntegrationConfig{
				Enabled:    false,
				BaseURL:    "https://api.ollama.com",
				APIKeyEnv:  "OLLAMA_API_KEY",
				CloudModel: "gemma4:31b-cloud",
			},
			GitHub: GitHubIntegrationConfig{
				Enabled:         false,
				TokenEnv:        "GH_TOKEN",
				PreferGHCLIAuth: true,
			},
			HuggingFace: HuggingFaceIntegrationConfig{
				Enabled:  false,
				TokenEnv: "HF_TOKEN",
				CacheDir: "~/.cache/huggingface",
			},
			Kaggle: KaggleIntegrationConfig{
				Enabled:   false,
				TokenEnv:  "KAGGLE_API_TOKEN",
				ConfigDir: "~/.kaggle",
			},
		},
		Transports: TransportsConfig{
			Discord: DiscordConfig{
				Enabled:                   false,
				BotTokenEnv:               "SAVITAR_DISCORD_BOT_TOKEN",
				DisplayName:               "Savitar",
				OperatorUserIDs:           []string{},
				RequireMention:            true,
				RespondInDirectMessages:   true,
				AllowCloudRepliesInGuilds: false,
				AllowCloudRepliesInDMs:    false,
				AllowLiveWebLookupInGuilds: false,
				AllowLiveWebLookupInDMs:    false,
				UseMessageContentIntent:   false,
				AllowedChannelIDs:         []string{},
				PerUserCooldownSeconds:    5,
				MaxConcurrentReplies:      2,
				MaxResponseChars:          1800,
				PresenceText:              "the repo",
				ImmediateAck:              true,
			},
			WhatsApp: WhatsAppConfig{Enabled: false, Bridge: "android-companion", DeviceName: "pixel-fold", NumberEnv: "SAVITAR_WHATSAPP_PRIMARY_NUMBER", DisplayName: "Savitar"},
			IMessage: IMessageConfig{Enabled: false, Bridge: "messages-applescript", AllowedSenders: []string{}, AccountEnv: "SAVITAR_IMESSAGE_ACCOUNT", DisplayName: "Savitar"},
		},
		WebUI: WebUIConfig{
			Enabled:                 false,
			PublicBaseURL:           "https://savitar.example.com",
			AuthProvider:            "google-oauth",
			GoogleClientIDEnv:       "SAVITAR_GOOGLE_CLIENT_ID",
			GoogleClientSecretEnv:   "SAVITAR_GOOGLE_CLIENT_SECRET",
			SessionSecretEnv:        "SAVITAR_WEB_SESSION_SECRET",
			AllowedEmailDomains:     []string{},
			RequireOperatorApproval: true,
		},
		Automation: AutomationConfig{
			PlaywrightProfileDir: ".savitar/browser",
			PeekabooEnabled:      true,
			AllowShellExecution:  true,
			ScreenCapture: ScreenCaptureConfig{
				ScreenshotsDir:      ".savitar/screenshots",
				RecordingsDir:       ".savitar/recordings",
				AllowDesktopControl: false,
			},
			RemoteHosts: []RemoteHostConfig{
				{Name: "local-mac-mini", Transport: "local-shell", Enabled: true},
			},
			MobileCompanion: MobileCompanionConfig{
				Enabled:    false,
				Platform:   "android",
				DeviceName: "pixel-fold",
				Phase:      "later",
			},
		},
		Memory: MemoryConfig{
			KnowledgeDir: ".savitar/knowledge",
			SnapshotDir:  ".savitar/snapshots",
			SessionIndex: ".savitar/session-index.json",
		},
		Knowledge: KnowledgeConfig{
			CatalogDir:              ".savitar/knowledge-catalog",
			IndexDir:                ".savitar/knowledge-index",
			MaxPacksPerReply:        3,
			MaxExcerptsPerReply:     6,
			RequireSourceMetadata:   true,
			RepoMarkdownDirs:        []string{"README.md", "docs/adr", "docs/roadmap", "docs/hackathon"},
			MaxRepoResults:          3,
			MaxGraphEdges:           6,
			EnableLiveWebLookup:     true,
			LiveWebProvider:         "duckduckgo-json",
			WebLookupTimeoutSeconds: 8,
		},
		Operator: OperatorConfig{
			ReviewQueueDir:             ".savitar/reviews",
			RunLogDir:                  ".savitar/runs",
			ApprovedEmails:             []string{},
			RequireReviewForToolErrors: true,
		},
		Skills: SkillsConfig{
			WorkspaceDir:       ".github/skills",
			UserDir:            "~/.copilot/skills",
			AllowRemoteCatalog: false,
		},
	}
}

func Load(path string) (Loaded, error) {
	loaded := Loaded{
		Config: Default(),
		Path:   path,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return loaded, nil
		}

		return Loaded{}, err
	}

	loaded.Exists = true
	if err := json.Unmarshal(data, &loaded.Config); err != nil {
		return Loaded{}, err
	}

	return loaded, nil
}
