package ki

import (
	"net/http"
	"net/url"
	"testing"
)

func TestLocation_NewLocation(t *testing.T) {
	loc := NewLocation(http.MethodGet, "/users/{id}")

	if loc.Method() != http.MethodGet {
		t.Errorf("expected method GET, got %s", loc.Method())
	}

	if loc.Pattern() != "/users/{id}" {
		t.Errorf("expected pattern /users/{id}, got %s", loc.Pattern())
	}
}

func TestLocation_WithPrefix(t *testing.T) {
	loc := NewLocation(http.MethodGet, "/users/{id}").WithPrefix("/api/v1")

	url := loc.URL()

	expected := "/api/v1/users/{id}"
	if url.Path != expected {
		t.Errorf("expected path %s, got %s", expected, url.Path)
	}
}

func TestLocation_WithPathParams(t *testing.T) {
	loc := NewLocation(http.MethodGet, "/users/{id}").WithPathParams("123")

	url := loc.URL()

	expected := "/users/123"
	if url.Path != expected {
		t.Errorf("expected path %s, got %s", expected, url.Path)
	}
}

func TestLocation_WithQuery(t *testing.T) {
	q := url.Values{}
	q.Add("page", "2")
	q.Add("sort", "asc")

	loc := NewLocation(http.MethodGet, "/items").WithQuery(q)

	url := loc.URL()

	expected := "page=2&sort=asc"
	if url.RawQuery != expected {
		t.Errorf("expected raw query %s, got %s", expected, url.RawQuery)
	}
}

func TestLocation_WithQueryParam(t *testing.T) {
	loc := NewLocation(http.MethodGet, "/search").WithQueryParam("q", "golang").WithQueryParam("limit", "10")

	url := loc.URL()

	expected := "limit=10&q=golang"
	if url.RawQuery != expected {
		t.Errorf("expected raw query %s, got %s", expected, url.RawQuery)
	}
}

func TestLocation_MultiplePathParams(t *testing.T) {
	loc := NewLocation(http.MethodGet, "/users/{userId}/posts/{postId}").WithPathParams("42", "99")

	url := loc.URL()

	expected := "/users/42/posts/99"
	if url.Path != expected {
		t.Errorf("expected path %s, got %s", expected, url.Path)
	}
}

func TestLocation_ExtraPathParamsIgnored(t *testing.T) {
	loc := NewLocation(http.MethodGet, "/projects/{id}").WithPathParams("1", "extra", "params")

	url := loc.URL()

	expected := "/projects/1"
	if url.Path != expected {
		t.Errorf("expected path %s, got %s", expected, url.Path)
	}
}

func TestLocation_InsufficientPathParamsKeepsPlaceholder(t *testing.T) {
	loc := NewLocation(http.MethodGet, "/teams/{id}/members/{mid}").WithPathParams("10")

	url := loc.URL()

	expected := "/teams/10/members/{mid}"
	if url.Path != expected {
		t.Errorf("expected path %s, got %s", expected, url.Path)
	}
}

func TestLocation_URLParsingError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic due to invalid URL, but did not panic")
		}
	}()

	badPattern := "://bad-url"
	loc := NewLocation(http.MethodGet, badPattern)
	loc.URL()
}
