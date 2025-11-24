package ki

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBufferedResponseWriter_WriteAndFlush(t *testing.T) {
	rec := httptest.NewRecorder()
	brw := NewBufferedResponseWriter(rec)

	data := []byte("hello world")
	n, err := brw.Write(data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n != len(data) {
		t.Errorf("expected to write %d bytes, wrote %d", len(data), n)
	}

	if brw.Size() != len(data) {
		t.Errorf("expected buffer size %d, got %d", len(data), brw.Size())
	}

	if rec.Body.Len() != 0 {
		t.Errorf("expected recorder body to be empty before flush, got %q", rec.Body.String())
	}

	brw.Flush()

	if rec.Body.String() != string(data) {
		t.Errorf("expected %q in response, got %q", data, rec.Body.String())
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestBufferedResponseWriter_WriteAfterFlush(t *testing.T) {
	rec := httptest.NewRecorder()
	brw := NewBufferedResponseWriter(rec)

	_, _ = brw.Write([]byte("first"))
	brw.Flush()

	_, _ = brw.Write([]byte("second"))

	expected := "firstsecond"
	if rec.Body.String() != expected {
		t.Errorf("expected %q in response, got %q", expected, rec.Body.String())
	}
}

func TestBufferedResponseWriter_StatusCode(t *testing.T) {
	rec := httptest.NewRecorder()
	brw := NewBufferedResponseWriter(rec)

	brw.WriteHeader(http.StatusTeapot)

	if brw.StatusCode() != http.StatusTeapot {
		t.Errorf("expected status code %d, got %d", http.StatusTeapot, brw.StatusCode())
	}

	_, _ = brw.Write([]byte("I'm a teapot"))
	brw.Flush()

	if rec.Code != http.StatusTeapot {
		t.Errorf("expected status code %d, got %d", http.StatusTeapot, rec.Code)
	}
}

func TestBufferedResponseWriter_Header(t *testing.T) {
	rec := httptest.NewRecorder()
	brw := NewBufferedResponseWriter(rec)

	brw.Header().Set("X-Test", "true")
	_, _ = brw.Write([]byte("test"))
	brw.Flush()

	if got := rec.Header().Get("X-Test"); got != "true" {
		t.Errorf("expected header X-Test=true, got %q", got)
	}
}

func TestBufferedResponseWriter_FlushIdempotent(t *testing.T) {
	rec := httptest.NewRecorder()
	brw := NewBufferedResponseWriter(rec)

	_, _ = brw.Write([]byte("data"))
	brw.Flush()
	brw.Flush()

	if rec.Body.String() != "data" {
		t.Errorf("expected body %q, got %q", "data", rec.Body.String())
	}
}

