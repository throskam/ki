package ki

import (
	"context"
	"log/slog"

	"golang.org/x/text/language"
)

type contextKey string

const (
	languageContextKey  contextKey = "language"
	loggerContextKey    contextKey = "logger"
	registryContextKey  contextKey = "registry"
	requestIDContextKey contextKey = "request-id"
)

// GetLocation returns the location for the given key from the registry in the context.
// Use with Locator middleware.
func GetLocation(ctx context.Context, key string) Location {
	registry := ctx.Value(registryContextKey).(*Registry)

	return registry.Get(key)
}

// SetRegistry sets the registry in the context.
func SetRegistry(ctx context.Context, registry *Registry) context.Context {
	return context.WithValue(ctx, registryContextKey, registry)
}

// GetRequestID returns the request ID from the context.
// Use with RequestID middleware.
func GetRequestID(ctx context.Context) string {
	reqID, ok := ctx.Value(requestIDContextKey).(string)
	if !ok {
		return ""
	}

	return reqID
}

// SetRequestID sets the request ID in the context.
func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey, requestID)
}

// MustGetLanguage returns the language from the context.
// If the language is not set, it panics.
// Use with Language middleware.
func MustGetLanguage(ctx context.Context) language.Tag {
	lang, ok := ctx.Value(languageContextKey).(language.Tag)
	if !ok {
		panic("no language")
	}

	return lang
}

// SetLanguage sets the language in the context.
func SetLanguage(ctx context.Context, lang language.Tag) context.Context {
	return context.WithValue(ctx, languageContextKey, lang)
}

// MustGetLogger returns the logger from the context.
// If the logger is not set, it panics.
// Use with Logger middleware.
func MustGetLogger(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerContextKey).(*slog.Logger)
	if !ok {
		panic("no logger")
	}

	return logger
}

// SetLogger sets the logger in the context.
func SetLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

