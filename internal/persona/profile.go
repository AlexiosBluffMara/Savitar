package persona

import "github.com/alexiosbluffmara/savitar/internal/config"

type Profile struct {
	Name              string
	Style             string
	Tone              string
	DesignBias        string
	CommentaryDensity string
	PublicBio         string
	DisclosurePolicy  string
}

func FromConfig(cfg config.AgentConfig) Profile {
	name := cfg.Name
	if name == "" {
		name = "Savitar"
	}

	return Profile{
		Name:              name,
		Style:             cfg.Style,
		Tone:              cfg.Persona.Tone,
		DesignBias:        cfg.Persona.DesignBias,
		CommentaryDensity: cfg.Persona.CommentaryDensity,
		PublicBio:         cfg.Persona.PublicBio,
		DisclosurePolicy:  cfg.Persona.DisclosurePolicy,
	}
}

func (p Profile) RequiresExplicitDisclosure() bool {
	return p.DisclosurePolicy != ""
}
