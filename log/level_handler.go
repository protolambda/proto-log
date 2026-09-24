package log

import (
	"context"
	"log/slog"
)

// LevelHandler filters records below a minimum level.
//
// The level is shared with every handler derived through WithAttrs or WithGroup,
// so changing it on any logger in a With-chain changes it for the whole chain,
// like a [slog.LevelVar] passed to the standard handlers.
// To give a derived logger its own (stricter) level, wrap its handler with another [LevelMod].
type LevelHandler struct {
	inner slog.Handler
	lvl   *slog.LevelVar
}

var _ Handler = (*LevelHandler)(nil)

func LevelMod(minLvl slog.Level) HandlerMod {
	return func(h slog.Handler) slog.Handler {
		out := &LevelHandler{inner: h, lvl: new(slog.LevelVar)}
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
	return &LevelHandler{inner: h.inner.WithAttrs(attrs), lvl: h.lvl}
}

func (h *LevelHandler) WithGroup(name string) slog.Handler {
	return &LevelHandler{inner: h.inner.WithGroup(name), lvl: h.lvl}
}

// MinLevel returns the minimum level. It is safe for concurrent use.
func (h *LevelHandler) MinLevel() slog.Level {
	return h.lvl.Level()
}

// SetMinLevel changes the minimum level for every handler sharing this level,
// i.e. all handlers derived from the same [LevelMod] application. It is safe for concurrent use.
func (h *LevelHandler) SetMinLevel(lvl slog.Level) {
	h.lvl.Set(lvl)
}
