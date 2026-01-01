package log

import (
	"context"
	"log/slog"
	"sync/atomic"
)

type LevelHandler struct {
	inner slog.Handler
	lvl   atomic.Int64 // slog.Level
}

var _ Handler = (*LevelHandler)(nil)

func LevelMod(minLvl slog.Level) HandlerMod {
	return func(h slog.Handler) slog.Handler {
		out := &LevelHandler{inner: h}
		out.SetMinLevel(minLvl)
		return out
	}
}

func (h *LevelHandler) Unwrap() slog.Handler {
	return h.inner
}

func (h *LevelHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	if lvl < h.MinLevel() {
		return false
	}
	return h.inner.Enabled(ctx, lvl)
}

func (h *LevelHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.inner.Handle(ctx, r)
}

func (h *LevelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	out := &LevelHandler{
		inner: h.inner.WithAttrs(attrs),
	}
	out.SetMinLevel(h.MinLevel())
	return out
}

func (h *LevelHandler) WithGroup(name string) slog.Handler {
	out := &LevelHandler{
		inner: h.inner.WithGroup(name),
	}
	out.SetMinLevel(h.MinLevel())
	return out
}

func (h *LevelHandler) MinLevel() slog.Level {
	return slog.Level(h.lvl.Load())
}

func (h *LevelHandler) SetMinLevel(lvl slog.Level) {
	h.lvl.Store(int64(lvl))
}
