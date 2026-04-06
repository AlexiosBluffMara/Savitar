package gateway

import "github.com/alexiosbluffmara/savitar/internal/config"

type Surface string

const (
	SurfaceDiscord  Surface = "discord"
	SurfaceWhatsApp Surface = "whatsapp"
	SurfaceIMessage Surface = "imessage"
	SurfaceWebUI    Surface = "webui"
)

type Attachment struct {
	Name     string
	MIMEType string
	Path     string
}

type Envelope struct {
	Surface           Surface
	ConversationID    string
	SenderID          string
	SenderDisplayName string
	ReplyToID         string
	Body              string
	Mentions          []string
	Attachments       []Attachment
	Metadata          map[string]string
}

type SurfaceConfig struct {
	Name     Surface
	Enabled  bool
	Kind     string
	Identity string
	Details  string
}

type Plan struct {
	Surfaces []SurfaceConfig
}

func BuildPlan(cfg config.Config) Plan {
	return Plan{
		Surfaces: []SurfaceConfig{
			{
				Name:     SurfaceDiscord,
				Enabled:  cfg.Transports.Discord.Enabled,
				Kind:     "messaging",
				Identity: cfg.Transports.Discord.DisplayName,
				Details:  cfg.Transports.Discord.BotTokenEnv,
			},
			{
				Name:     SurfaceWhatsApp,
				Enabled:  cfg.Transports.WhatsApp.Enabled,
				Kind:     cfg.Transports.WhatsApp.Bridge,
				Identity: cfg.Transports.WhatsApp.DisplayName,
				Details:  cfg.Transports.WhatsApp.DeviceName,
			},
			{
				Name:     SurfaceIMessage,
				Enabled:  cfg.Transports.IMessage.Enabled,
				Kind:     cfg.Transports.IMessage.Bridge,
				Identity: cfg.Transports.IMessage.DisplayName,
				Details:  cfg.Transports.IMessage.AccountEnv,
			},
			{
				Name:     SurfaceWebUI,
				Enabled:  cfg.WebUI.Enabled,
				Kind:     cfg.WebUI.AuthProvider,
				Identity: cfg.Agent.Name,
				Details:  cfg.WebUI.PublicBaseURL,
			},
		},
	}
}

func (p Plan) EnabledCount() int {
	count := 0
	for _, surface := range p.Surfaces {
		if surface.Enabled {
			count++
		}
	}

	return count
}
