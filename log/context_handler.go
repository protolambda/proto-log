package log

import (
	"context"
	"log/slog"
)

// ContextHandler replaces the context of logging calls that use the default context
// with a configured default context.
//
// A call uses the default context when its context is nil or is exactly [context.Background],
// which is what the non-Context logging methods (e.g. Info) pass.
// Any other context, including [context.TODO] and contexts derived from [context.Background],
// is passed through unchanged.
//
// A ContextHandler is immutable: use [Logger.WithContext] to derive a logger with a different default context.
type ContextHandler struct {
	inner slog.Handler
	ctx   context.Context
}

var _ Handler = (*ContextHandler)(nil)

// ContextMod wraps a handler with a [ContextHandler] that uses [context.Background] as default context.
func ContextMod() HandlerMod {
	return ContextModWith(context.Background())
}

// ContextModWith wraps a handler with a [ContextHandler] that uses ctx as default context.
func ContextModWith(ctx context.Context) HandlerMod {
	if ctx == nil {
		ctx = context.Background()
	}
	return func(h slog.Handler) slog.Handler {
		return &ContextHandler{inner: h, ctx: ctx}
	}
}

func (h *ContextHandler) Unwrap() slog.Handler {
	return h.inner
}

func (h *ContextHandler) resolve(ctx context.Context) context.Context {
	if ctx == nil || ctx == context.Background() {
		return h.ctx
	}
	return ctx
}

func (h *ContextHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return h.inner.Enabled(h.resolve(ctx), lvl)
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.inner.Handle(h.resolve(ctx), r)
}

func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{
		inner: h.inner.WithAttrs(attrs),
		ctx:   h.ctx,
	}
}

func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{
		inner: h.inner.WithGroup(name),
		ctx:   h.ctx,
	}
}

// Context returns the default context.
func (h *ContextHandler) Context() context.Context {
	return h.ctx
}
