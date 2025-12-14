package log

import (
	"context"
	"log/slog"
)

// PostProcessHandler allows you to post-process logs
type PostProcessHandler struct {
	inner slog.Handler
	fn    HandlerFunc
}

var _ Handler = (*PostProcessHandler)(nil)

func PostProcessMod(fn HandlerFunc) HandlerMod {
	return func(h slog.Handler) slog.Handler {
		return &PostProcessHandler{inner: h, fn: fn}
	}
}

func (h *PostProcessHandler) Unwrap() slog.Handler {
	return h.inner
}

func (h *PostProcessHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return h.inner.Enabled(ctx, lvl)
}

func (h *PostProcessHandler) Handle(ctx context.Context, r slog.Record) error {
	defer h.fn(ctx, r)
	return h.inner.Handle(ctx, r)
}

func (h *PostProcessHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PostProcessHandler{
		inner: h.inner.WithAttrs(attrs),
		fn:    h.fn,
	}
}

func (h *PostProcessHandler) WithGroup(name string) slog.Handler {
	return &PostProcessHandler{
		inner: h.inner.WithGroup(name),
		fn:    h.fn,
	}
}
