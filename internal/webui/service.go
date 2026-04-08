package webui

import (
	"encoding/json"
	"net/http"

	"github.com/alexiosbluffmara/savitar/internal/operator"
)

// Session is the authenticated operator identity for a web request.
type Session struct {
	Email    string
	Approved bool
	Roles    []string
}

// Authenticator resolves an operator session from an inbound request.
type Authenticator interface {
	SessionFromRequest(*http.Request) (Session, error)
}

// Service is the backend HTTP seam for the operator console. The rewrite keeps
// this service backend-first so authentication, review, and run visibility can
// be built before choosing how much frontend complexity is necessary.
type Service struct {
	Operator      *operator.Service
	Authenticator Authenticator
}

// Routes returns the initial authenticated API surface for the operator UI.
func (s Service) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/snapshot", s.handleSnapshot)
	mux.HandleFunc("/api/reviews", s.handleReviews)
	mux.HandleFunc("/api/reviews/decision", s.handleDecision)
	return mux
}

func (s Service) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !s.authorized(w, r) {
		return
	}
	if s.Operator == nil {
		http.Error(w, "operator service not wired", http.StatusNotImplemented)
		return
	}
	snapshot, err := s.Operator.Snapshot(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(snapshot)
}

func (s Service) handleReviews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !s.authorized(w, r) {
		return
	}
	if s.Operator == nil {
		http.Error(w, "operator service not wired", http.StatusNotImplemented)
		return
	}
	items, err := s.Operator.PendingReviews(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(items)
}

func (s Service) handleDecision(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !s.authorized(w, r) {
		return
	}
	if s.Operator == nil {
		http.Error(w, "operator service not wired", http.StatusNotImplemented)
		return
	}
	var decision operator.Decision
	if err := json.NewDecoder(r.Body).Decode(&decision); err != nil {
		http.Error(w, "invalid decision payload", http.StatusBadRequest)
		return
	}
	if err := s.Operator.ApplyDecision(r.Context(), decision); err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (s Service) authorized(w http.ResponseWriter, r *http.Request) bool {
	if s.Authenticator == nil {
		http.Error(w, "authenticator not wired", http.StatusNotImplemented)
		return false
	}
	session, err := s.Authenticator.SessionFromRequest(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return false
	}
	if !session.Approved {
		http.Error(w, "operator approval required", http.StatusForbidden)
		return false
	}
	return true
}