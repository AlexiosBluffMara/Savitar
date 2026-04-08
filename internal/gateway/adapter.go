package gateway

import (
	"context"
	"errors"
)

// ErrNotImplemented is returned by later-pass adapter methods that are present
// as scaffolding before a concrete transport or trust implementation exists.
var ErrNotImplemented = errors.New("gateway adapter method not implemented")

// Reference is the normalized citation metadata attached to a reply or delivery.
// It is intentionally transport-neutral so the same evidence record can be shown
// in Discord, the operator console, and later review logs.
type Reference struct {
	Label   string
	Source  string
	Locator string
	Kind    string
}

// DeliveryStatus records the observed result of an outbound send attempt.
type DeliveryStatus string

const (
	DeliveryPending  DeliveryStatus = "pending"
	DeliverySent     DeliveryStatus = "sent"
	DeliveryDeferred DeliveryStatus = "deferred"
	DeliveryFailed   DeliveryStatus = "failed"
)

// Delivery is the outbound counterpart to Envelope. The orchestrator produces
// this value after trust, routing, evidence, and composition are complete.
type Delivery struct {
	Envelope         Envelope
	Body             string
	References       []Reference
	RequiresApproval bool
	Metadata         map[string]string
}

// DeliveryResult is the transport-facing outcome of an outbound send.
type DeliveryResult struct {
	Status     DeliveryStatus
	Surface    Surface
	ExternalID string
	Details    string
}

// TrustDecision captures the policy outcome for an inbound or outbound action.
// Later implementations can expand this with risk scores or approval reasons
// without changing the gateway role in the rewrite.
type TrustDecision struct {
	Allowed        bool
	RequiresReview bool
	Reason         string
}

// InboundAdapter streams normalized inbound traffic from a single surface into
// the shared orchestration path.
type InboundAdapter interface {
	Surface() Surface
	Start(context.Context, func(context.Context, Envelope) error) error
}

// OutboundAdapter delivers a normalized reply back to one surface.
type OutboundAdapter interface {
	Surface() Surface
	Deliver(context.Context, Delivery) (DeliveryResult, error)
}

// TrustEvaluator applies per-surface policy before the runtime commits to a
// reply or operator-visible action.
type TrustEvaluator interface {
	Evaluate(context.Context, Envelope) (TrustDecision, error)
}