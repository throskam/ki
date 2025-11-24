package ki

import (
	"log/slog"
	"os"
)

// Logger is the global logger.
var Logger *slog.Logger

var loggerLevel *slog.LevelVar

// SetLoggerLevelByText sets the logger level to the given text level.
// The text level must be one of the following: debug, info, warn, error.
func SetLoggerLevelByText(s string) {
	var level slog.Level

	err := level.UnmarshalText([]byte(s))
	if err != nil {
		panic("unknown slog level")
	}

	loggerLevel.Set(level)
}

func init() {
	loggerLevel = &slog.LevelVar{}

	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: loggerLevel,
	}))
}