package ki

import (
	"context"
	"log/slog"
	"testing"
)

func TestSetLoggerLevelByText_ValidLevels(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("unexpected panic for input %q: %v", tt.input, r)
				}
			}()

			SetLoggerLevelByText(tt.input)

			if level := Logger.Enabled(context.Background(), tt.expected); !level {
				t.Errorf("expected logger to be enabled for level %q", tt.input)
			}
		})
	}
}

func TestSetLoggerLevelByText_InvalidLevel(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid level, but got none")
		}
	}()

	SetLoggerLevelByText("invalid-level")
}