package log

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"sync"
)

type terminalHandler struct {
	out *terminalOutput // shared among derived handlers
	cfg *FormatConfig

	// attrs are the flattened attributes inherited through WithAttrs, with group-qualified keys.
	attrs []slog.Attr
	// groupPrefix qualifies the keys of attributes added after WithGroup, e.g. "a.b.".
	groupPrefix string
}

// terminalOutput is the state that handlers writing to the same output share,
// so that concurrent writes do not interleave and field padding stays aligned across sub-loggers.
type terminalOutput struct {
	mu sync.Mutex
	wr io.Writer

	// fieldPadding is a map with maximum field value lengths seen until now
	// to allow padding log contexts in a bit smarter way.
	fieldPadding map[string]int

	buf   bytes.Buffer
	attrs []slog.Attr // scratch space for flattened record attributes
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
//
// Attribute values are padded to the longest value seen so far for the same key,
// to align subsequent log lines. Handlers derived with WithAttrs and WithGroup share this padding state.
// Groups are flattened into dot-separated keys.
func TerminalHandler(wr io.Writer, opts ...FormatOption) slog.Handler {
	out := &terminalHandler{
		out: &terminalOutput{
			wr:           wr,
			fieldPadding: make(map[string]int),
		},
		cfg: &FormatConfig{
			UseColor:      false,
			IncludeSource: false,
			ExcludeTime:   false,
			SourceRelDir:  "",
		},
	}
	out.cfg.Apply(opts...)
	return out
}

// Handle writes the formatted record, and returns any error from the underlying writer.
func (h *terminalHandler) Handle(_ context.Context, r slog.Record) error {
	o := h.out
	o.mu.Lock()
	defer o.mu.Unlock()
	o.attrs = o.attrs[:0]
	r.Attrs(func(a slog.Attr) bool {
		o.attrs = appendFlatAttr(o.attrs, h.groupPrefix, a)
		return true
	})
	o.buf.Reset()
	h.format(&o.buf, r, o.attrs)
	_, err := o.wr.Write(o.buf.Bytes())
	return err
}

func (h *terminalHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true
}

func (h *terminalHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	c := *h
	c.groupPrefix = h.groupPrefix + name + "."
	return &c
}

func (h *terminalHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	c := *h
	c.attrs = h.attrs[:len(h.attrs):len(h.attrs)] // force a copy on append, siblings must not share storage
	for _, a := range attrs {
		c.attrs = appendFlatAttr(c.attrs, h.groupPrefix, a)
	}
	return &c
}

// appendFlatAttr appends a to dst, following the slog.Handler rules:
// values are resolved, empty attributes and empty groups are dropped,
// groups with an empty key are inlined, and other groups qualify the keys of their members.
func appendFlatAttr(dst []slog.Attr, prefix string, a slog.Attr) []slog.Attr {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return dst
	}
	if a.Value.Kind() != slog.KindGroup {
		a.Key = prefix + a.Key
		return append(dst, a)
	}
	if a.Key != "" {
		prefix = prefix + a.Key + "."
	}
	for _, member := range a.Value.Group() {
		dst = appendFlatAttr(dst, prefix, member)
	}
	return dst
}

// ResetFieldPadding zeroes the field-padding for all attribute pairs.
func (h *terminalHandler) ResetFieldPadding() {
	h.out.mu.Lock()
	h.out.fieldPadding = make(map[string]int)
	h.out.mu.Unlock()
}
