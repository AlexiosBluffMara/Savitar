package models

type Complexity string

const (
	ComplexityRoutine  Complexity = "routine"
	ComplexityStandard Complexity = "standard"
	ComplexityComplex  Complexity = "complex"
)

type Task struct {
	Complexity       Complexity
	PrivateContext   bool
	RequiresLocal    bool
	LatencySensitive bool
}

type Profile struct {
	Name            string
	Provider        string
	Model           string
	Purpose         string
	UsageMultiplier float64
}

type Decision struct {
	Profile Profile
	Reason  string
}

type Router struct {
	local        Profile
	routineLane  Profile
	standardLane Profile
	complexLane  Profile
}

func DefaultProfiles() []Profile {
	return []Profile{
		{
			Name:            "local-default",
			Provider:        "ollama",
			Model:           "gemma4:e4b",
			Purpose:         "Private or latency-sensitive local work on Apple Silicon.",
			UsageMultiplier: 0,
		},
		{
			Name:            "copilot-0x",
			Provider:        "copilot",
			Model:           "managed-routine-lane",
			Purpose:         "Routine work with effectively unlimited usage.",
			UsageMultiplier: 0,
		},
		{
			Name:            "copilot-0.33x",
			Provider:        "copilot",
			Model:           "managed-standard-lane",
			Purpose:         "Medium-complexity synthesis and planning.",
			UsageMultiplier: 0.33,
		},
		{
			Name:            "copilot-1x",
			Provider:        "copilot",
			Model:           "managed-complex-lane",
			Purpose:         "High-complexity reasoning and multi-step execution.",
			UsageMultiplier: 1,
		},
	}
}

func DefaultRouter() Router {
	profiles := DefaultProfiles()
	return Router{
		local:        profiles[0],
		routineLane:  profiles[1],
		standardLane: profiles[2],
		complexLane:  profiles[3],
	}
}

func (r Router) Profiles() []Profile {
	return []Profile{r.local, r.routineLane, r.standardLane, r.complexLane}
}

func (r Router) Route(task Task) Decision {
	if task.RequiresLocal || task.PrivateContext || task.LatencySensitive {
		return Decision{
			Profile: r.local,
			Reason:  "local lane selected for privacy or latency sensitivity",
		}
	}

	switch task.Complexity {
	case ComplexityComplex:
		return Decision{Profile: r.complexLane, Reason: "complex work escalates to the 1x lane"}
	case ComplexityStandard:
		return Decision{Profile: r.standardLane, Reason: "standard work uses the 0.33x lane"}
	default:
		return Decision{Profile: r.routineLane, Reason: "routine work stays on the 0x lane"}
	}
}
