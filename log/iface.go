package log

import (
	"context"
	"log/slog"
)

// Handler is a slog.Handler you can unwrap,
// to access inner handler functionality.
type Handler interface {
	slog.Handler
	Unwrap() slog.Handler
}

type SLogLogger interface {
	// Debug logs a message at the debug level with context key/value pairs
	Debug(msg string, args ...any)

	// Info logs a message at the info level with context key/value pairs
	Info(msg string, args ...any)

	// Warn logs a message at the warn level with context key/value pairs
	Warn(msg string, args ...any)

	// Error logs a message at the error level with context key/value pairs
	Error(msg string, args ...any)

	// DebugContext logs at [LevelDebug] with the given context.
	DebugContext(ctx context.Context, msg string, args ...any)

	// InfoContext logs at [LevelInfo] with the given context.
	InfoContext(ctx context.Context, msg string, args ...any)

	// WarnContext logs at [LevelWarn] with the given context.
	WarnContext(ctx context.Context, msg string, args ...any)

	// ErrorContext logs at [LevelError] with the given context.
	ErrorContext(ctx context.Context, msg string, args ...any)

	// Log logs a message at the specified level
	Log(ctx context.Context, level slog.Level, msg string, attrs ...any)

	// LogAttrs is a more efficient version of [Logger.Log] that accepts only Attrs.
	LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr)

	// Handler returns the underlying handler of the inner logger.
	Handler() slog.Handler

	// Enabled reports whether l emits log records at the given context and level.
	Enabled(ctx context.Context, level slog.Level) bool
}

var _ SLogLogger = (*slog.Logger)(nil)

// ExtendedSLogLogger adds Trace and Crit levels to the common SLogLogger interface.
type ExtendedSLogLogger interface {
	SLogLogger

	// Trace log a message at the trace level with context key/value pairs.
	Trace(msg string, args ...any)

	// Crit logs a message at the crit level with context key/value pairs.
	Crit(msg string, args ...any)

	// TraceContext logs at [LevelTrace] with the given context.
	TraceContext(ctx context.Context, msg string, args ...any)

	// CritContext logs at [LevelCrit] with the given context.
	CritContext(ctx context.Context, msg string, args ...any)
}

type Logger interface {
	ExtendedSLogLogger

	// With returns a Logger that includes the given attributes
	// in each output operation. Arguments are converted to
	// attributes as if by [Logger.Log].
	With(args ...any) Logger

	// WithGroup returns a Logger that starts a group, if name is non-empty.
	// The keys of all attributes added to the Logger will be qualified by the given
	// name. (How that qualification happens depends on the [Handler.WithGroup]
	// method of the Logger's Handler.)
	//
	// If name is empty, WithGroup returns the receiver.
	WithGroup(name string) Logger

	// Context returns the default context that is used when logging
	Context() context.Context

	// WithContext creates a clone, with the given context as new default context.
	WithContext(ctx context.Context) Logger
}

// HandlerFunc fits the slog.Handler.Handle method
type HandlerFunc func(ctx context.Context, r slog.Record)
