package models

import "testing"

func TestRouteUsesLocalLaneForPrivateTasks(t *testing.T) {
	router := DefaultRouter()
	decision := router.Route(Task{Complexity: ComplexityComplex, PrivateContext: true})

	if decision.Profile.Name != "local-default" {
		t.Fatalf("expected local-default, got %q", decision.Profile.Name)
	}
}

func TestRouteUsesComplexLaneForComplexWork(t *testing.T) {
	router := DefaultRouter()
	decision := router.Route(Task{Complexity: ComplexityComplex})

	if decision.Profile.Name != "copilot-1x" {
		t.Fatalf("expected copilot-1x, got %q", decision.Profile.Name)
	}
}

func TestRouteUsesRoutineLaneByDefault(t *testing.T) {
	router := DefaultRouter()
	decision := router.Route(Task{})

	if decision.Profile.Name != "copilot-0x" {
		t.Fatalf("expected copilot-0x, got %q", decision.Profile.Name)
	}
}
