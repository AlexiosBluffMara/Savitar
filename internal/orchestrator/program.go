package orchestrator

import (
	"context"
	"errors"
	"time"

	"github.com/alexiosbluffmara/savitar/internal/gateway"
	"github.com/alexiosbluffmara/savitar/internal/memory"
	"github.com/alexiosbluffmara/savitar/internal/models"
)

// ErrNotImplemented marks the orchestrator as scaffolding until later passes
// fill in the concrete implementation behind the rewrite blueprint.
var ErrNotImplemented = errors.New("orchestrator program method not implemented")

// Intent is the high-level interpretation of an inbound message. It should stay
// product-facing rather than transport-facing so the same intent model works in
// Discord, the operator console, and later expansion surfaces.
type Intent struct {
	Name                 string
	SubjectHint          string
	RequiresOperator     bool
	RequiresLiveEvidence bool
	RequiresReview       bool
}

// SessionState is the rewrite-oriented per-conversation state the orchestrator
// needs while planning a turn. The existing session package can later become
// one concrete implementation of this contract.
type SessionState struct {
	ConversationID      string
	CurrentSurface      string
	LastIntent          string
	LastSubject         string
	CurrentModelProfile string
	UpdatedAt           time.Time
}

// ToolCall is the normalized description of one live evidence request.
type ToolCall struct {
	ServerName string
	ToolName   string
	Arguments  map[string]any
	Reason     string
}

// ToolResult is the normalized outcome of one live evidence request.
type ToolResult struct {
	Call      ToolCall
	Content   string
	IsError   bool
	Reference gateway.Reference
}

// TurnPlan is the orchestrator's pre-execution plan for one inbound message.
// Later passes should make this structure explicit in logs and operator review
// so route choices and tool use are inspectable.
type TurnPlan struct {
	Intent         Intent
	Task           models.Task
	MemoryQuery    memory.Query
	ToolCalls      []ToolCall
	RequiresReview bool
}

// TurnContext is the assembled context available at composition time.
type TurnContext struct {
	Session     SessionState
	Memory      memory.RetrievalResult
	ToolResults []ToolResult
	Route       models.Decision
}

// ReplyDraft is the transport-neutral output from composition before delivery.
type ReplyDraft struct {
	Body           string
	References     []gateway.Reference
	RequiresReview bool
}

// RunRecord is the durable explanation of what the orchestrator did.
type RunRecord struct {
	Surface        gateway.Surface
	ConversationID string
	Intent         string
	Subject        string
	Route          models.Decision
	ToolCalls      []ToolCall
	References     []gateway.Reference
	RequiresReview bool
	StartedAt      time.Time
	CompletedAt    time.Time
}

// Result is the complete output of one handled turn.
type Result struct {
	Delivery gateway.Delivery
	Record   RunRecord
	Reply    ReplyDraft
}

// IntentResolver turns a normalized envelope into a product-level intent and
// task classification.
type IntentResolver interface {
	ResolveIntent(context.Context, gateway.Envelope) (Intent, models.Task, error)
}

// SessionStore abstracts the conversation state the orchestrator needs to load
// and update around each turn.
type SessionStore interface {
	LoadSession(context.Context, gateway.Envelope) (SessionState, error)
	SaveSession(context.Context, SessionState) error
}

// Planner builds the turn plan after intent resolution and session loading.
type Planner interface {
	PlanTurn(context.Context, gateway.Envelope, Intent, models.Task, SessionState) (TurnPlan, error)
}

// ToolBroker executes any live evidence calls requested by the plan.
type ToolBroker interface {
	Execute(context.Context, []ToolCall) ([]ToolResult, error)
}

// RouteSelector chooses the model lane for the turn.
type RouteSelector interface {
	SelectRoute(context.Context, TurnPlan, TurnContext) (models.Decision, error)
}

// Composer generates the final draft after knowledge, tools, and route choice
// are available.
type Composer interface {
	ComposeReply(context.Context, gateway.Envelope, TurnPlan, TurnContext) (ReplyDraft, error)
}

// Recorder persists the durable run record used by operator review.
type Recorder interface {
	RecordRun(context.Context, RunRecord) error
}

// Program is the top-level rewrite seam. Later passes should implement this
// struct rather than growing more reply logic directly inside a transport.
type Program struct {
	IntentResolver IntentResolver
	SessionStore   SessionStore
	Retriever      memory.Retriever
	Planner        Planner
	ToolBroker     ToolBroker
	RouteSelector  RouteSelector
	Composer       Composer
	Recorder       Recorder
	Clock          func() time.Time
}

// HandleTurn is intentionally left as a scaffold. The later implementation
// should execute the rewrite lifecycle defined in docs/roadmap/0006.
func (p Program) HandleTurn(context.Context, gateway.Envelope) (Result, error) {
	return Result{}, ErrNotImplemented
}
