package log

import (
	"context"
	"log/slog"
)

type LevelHandler struct {
	inner slog.Handler
	lvl   slog.Level
}

var _ Handler = (*LevelHandler)(nil)

func LevelMod(minLvl slog.Level) HandlerMod {
	return func(h slog.Handler) slog.Handler {
		return &LevelHandler{inner: h, lvl: minLvl}
	}
}

func (h *LevelHandler) Unwrap() slog.Handler {
	return h.inner
}

func (h *LevelHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	if lvl < h.lvl {
		return false
	}
	return h.inner.Enabled(ctx, lvl)
}

func (h *LevelHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.inner.Handle(ctx, r)
}

func (h *LevelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LevelHandler{
		inner: h.inner.WithAttrs(attrs),
		lvl:   h.lvl,
	}
}

func (h *LevelHandler) WithGroup(name string) slog.Handler {
	return &LevelHandler{
		inner: h.inner.WithGroup(name),
		lvl:   h.lvl,
	}
}

func (h *LevelHandler) MinLevel() slog.Level {
	return h.lvl
}

func (h *LevelHandler) SetMinLevel(lvl slog.Level) {
	h.lvl = lvl
}
