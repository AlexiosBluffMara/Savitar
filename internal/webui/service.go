package webui

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alexiosbluffmara/savitar/internal/operator"
	savitarruntime "github.com/alexiosbluffmara/savitar/internal/runtime"
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
	Runtime       *savitarruntime.Runtime
	Operator      *operator.Service
	Authenticator Authenticator
	DemoMode      bool
}

// Routes returns the initial authenticated API surface for the operator UI.
func (s Service) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleHome)
	mux.HandleFunc("/api/dashboard", s.handleDashboard)
	mux.HandleFunc("/api/snapshot", s.handleSnapshot)
	mux.HandleFunc("/api/reviews", s.handleReviews)
	mux.HandleFunc("/api/reviews/decision", s.handleDecision)
	return mux
}

func (s Service) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !s.authorized(w, r) {
		return
	}
	dashboard, err := s.dashboardData(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := dashboardTemplate.Execute(w, dashboard); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s Service) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !s.authorized(w, r) {
		return
	}
	dashboard, err := s.dashboardData(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dashboard)
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

// RunDemo starts the local demo web UI and blocks until the server exits.
func RunDemo(out io.Writer, errOut io.Writer, rt *savitarruntime.Runtime, addr string) int {
	if !rt.Config().Config.WebUI.Enabled {
		fmt.Fprintln(out, "warning: webui.enabled=false in config; serving because the command was invoked explicitly")
	}

	operatorService := operator.Service{
		Store: NewRuntimeStore(rt),
		Clock: time.Now,
	}
	service := Service{
		Runtime:       rt,
		Operator:      &operatorService,
		Authenticator: DemoAuthenticator{},
		DemoMode:      true,
	}

	if addr == "" {
		addr = ":8080"
	}
	fmt.Fprintf(out, "starting web UI prototype at %s\n", addr)
	fmt.Fprintln(out, "mode: local demo only; Google OAuth, rate limits, and public hardening are still pending")
	if err := http.ListenAndServe(addr, service.Routes()); err != nil {
		fmt.Fprintf(errOut, "web UI exited with error: %v\n", err)
		return 1
	}
	return 0
}
