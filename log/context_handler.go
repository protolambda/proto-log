package log

import (
	"context"
	"log/slog"
)

// ContextHandler allows you to change the context of logging calls
// that would otherwise use the default context.Background.
type ContextHandler struct {
	inner slog.Handler
	ctx   context.Context
}

var _ Handler = (*ContextHandler)(nil)

func ContextMod() HandlerMod {
	return func(h slog.Handler) slog.Handler {
		return &ContextHandler{inner: h, ctx: context.Background()}
	}
}

func (h *ContextHandler) Unwrap() slog.Handler {
	return h.inner
}

func (h *ContextHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	if ctx == context.Background() {
		ctx = h.ctx
	}
	return h.inner.Enabled(ctx, lvl)
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx == context.Background() {
		ctx = h.ctx
	}
	return h.inner.Handle(ctx, r)
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

func (h *ContextHandler) Context() context.Context {
	return h.ctx
}

func (h *ContextHandler) SetContext(ctx context.Context) {
	h.ctx = ctx
}
