package chat

import (
	"context"

	"github.com/alexiosbluffmara/savitar/internal/models"
)

type Turn struct {
	Role    string
	Content string
}

// ToolContext holds a pre-fetched tool result to be injected into the reply.
type ToolContext struct {
	ServerName string
	ToolName   string
	Result     string
}

type Request struct {
	Surface           string
	ConversationID    string
	SenderDisplayName string
	UserInput         string
	History           []Turn
	Task              models.Task
	Route             models.Decision
	ReplyLimit        int
	ToolContexts      []ToolContext // pre-fetched tool results
	MemoryContext     []string      // loaded memory pack bodies
}

type Generator interface {
	Generate(context.Context, Request) (string, error)
}
