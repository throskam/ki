package ki

import (
	"testing"
)

func TestRegistry_AddsAndGetsLocation(t *testing.T) {
	reg := NewRegistry()
	reg.Add("home", "GET", "/home")

	loc := reg.Get("home")

	if loc.Method() != "GET" {
		t.Errorf("expected method GET, got %s", loc.Method())
	}
	if loc.Pattern() != "/home" {
		t.Errorf("expected pattern /home, got %s", loc.Pattern())
	}
}

func TestRegistry_PanicsOnDuplicateAdd(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on duplicate Add")
		}
	}()

	reg := NewRegistry()
	reg.Add("home", "GET", "/home")
	reg.Add("home", "GET", "/home-again")
}

func TestRegistry_RemovesLocation(t *testing.T) {
	reg := NewRegistry()
	reg.Add("about", "GET", "/about")
	reg.Remove("about")

	if reg.Has("about") {
		t.Errorf("expected route to be removed")
	}
}

func TestRegistry_HasReturnsCorrectValues(t *testing.T) {
	reg := NewRegistry()
	reg.Add("ping", "GET", "/ping")

	if !reg.Has("ping") {
		t.Errorf("expected to have key 'ping'")
	}
	if reg.Has("pong") {
		t.Errorf("did not expect to have key 'pong'")
	}
}

func TestRegistry_ChildRegistryInherits(t *testing.T) {
	parent := NewRegistry()
	child := parent.Child("/v1")
	child.Add("status", "GET", "/status")

	if !parent.Has("status") {
		t.Errorf("expected parent to have key from child")
	}

	loc := parent.Get("status")

	if loc.Method() != "GET" {
		t.Errorf("expected GET, got %s", loc.Method())
	}
	if loc.Pattern() != "/status" {
		t.Errorf("expected /status, got %s", loc.Pattern())
	}

	url := loc.URL()
	if url.Path != "/v1/status" {
		t.Errorf("expected path /v1/status, got %s", url.Path)
	}
}

func TestRegistry_PanicsOnMissingKey(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when getting missing key")
		}
	}()

	reg := NewRegistry()
	reg.Get("nonexistent") // should panic
}