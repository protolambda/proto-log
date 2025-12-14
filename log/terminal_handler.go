package log

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"sync"
)

type terminalHandler struct {
	mu sync.Mutex
	wr io.Writer

	cfg *FormatConfig

	attrs []slog.Attr

	// fieldPadding is a map with maximum field value lengths seen until now
	// to allow padding log contexts in a bit smarter way.
	fieldPadding map[string]int

	buf *bytes.Buffer
}

// TerminalHandler returns a handler which formats log records at all levels optimized for human readability on
// a terminal with color-coded level output and terser human friendly timestamp.
// This format should only be used for interactive programs or while developing.
//
//	[LEVEL] [TIME] MESSAGE key=value key=value ...
//
// Example:
//
//	[DBUG] [May 16 20:58:45] remove route ns=haproxy addr=127.0.0.1:50002
func TerminalHandler(wr io.Writer, opts ...FormatOption) slog.Handler {
	out := &terminalHandler{
		wr:           wr,
		fieldPadding: make(map[string]int),
		cfg: &FormatConfig{
			UseColor:      false,
			IncludeSource: false,
			ExcludeTime:   true,
			SourceRelDir:  "",
		},
	}
	out.cfg.Apply(opts...)
	return out
}

func (h *terminalHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	buf := h.format(r)
	h.wr.Write(buf)
	h.buf.Reset()
	return nil
}

func (h *terminalHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true
}

func (h *terminalHandler) WithGroup(name string) slog.Handler {
	panic("not implemented")
}

func (h *terminalHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &terminalHandler{
		wr:           h.wr,
		cfg:          h.cfg,
		attrs:        append(h.attrs, attrs...),
		fieldPadding: make(map[string]int),
		buf:          nil,
	}
}

// ResetFieldPadding zeroes the field-padding for all attribute pairs.
func (h *terminalHandler) ResetFieldPadding() {
	h.mu.Lock()
	h.fieldPadding = make(map[string]int)
	h.mu.Unlock()
}
