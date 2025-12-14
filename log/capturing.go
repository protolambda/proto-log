package log

import (
	"context"
	"log/slog"
)

type Capturer interface {
	slog.Handler
	Clear()
	FindLog(filters ...LogFilter) *CapturedRecord
	FindLogs(filters ...LogFilter) []*CapturedRecord
}

var _ Capturer = (*CapturingHandler)(nil)

// CapturedAttrs forms a chain of inherited attributes, to traverse on captured log records.
type CapturedAttrs struct {
	Parent     *CapturedAttrs
	Attributes []slog.Attr
}

// Attrs calls f on each Attr in the [CapturedAttrs].
// Iteration stops if f returns false.
func (r *CapturedAttrs) Attrs(f func(slog.Attr) bool) {
	for _, a := range r.Attributes {
		if !f(a) {
			return
		}
	}
	if r.Parent != nil {
		r.Parent.Attrs(f)
	}
}

// CapturedRecord is a wrapped around a regular log-record,
// to preserve the inherited attributes context, without mutating the record or reordering attributes.
type CapturedRecord struct {
	Parent *CapturedAttrs
	*slog.Record
}

// Attrs calls f on each Attr in the [CapturedRecord].
// Iteration stops if f returns false.
func (r *CapturedRecord) Attrs(f func(slog.Attr) bool) {
	searching := true
	r.Record.Attrs(func(a slog.Attr) bool {
		searching = f(a)
		return searching
	})
	if !searching { // if we found it already, then don't traverse the remainder
		return
	}
	if r.Parent != nil {
		r.Parent.Attrs(f)
	}
}

func (r *CapturedRecord) AttrValue(key string) (v any) {
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == key {
			v = a.Value.Any()
			return false
		}
		return true // try next
	})
	return
}

// CapturingHandler provides a log handler that captures all log records and optionally forwards them to a delegate.
// Note that it is not thread safe.
type CapturingHandler struct {
	handler slog.Handler
	Logs    *[]*CapturedRecord // shared among derived CapturingHandlers
	// attrs are inherited log record attributes, from a logger that this CapturingHandler may be derived from
	attrs *CapturedAttrs
}

var _ Handler = (*CapturingHandler)(nil)

func CapturingMod() HandlerMod {
	return func(h slog.Handler) slog.Handler {
		return &CapturingHandler{handler: h, Logs: new([]*CapturedRecord)}
	}
}

func (c *CapturingHandler) Unwrap() slog.Handler {
	return c.handler
}

func (c *CapturingHandler) Handle(ctx context.Context, r slog.Record) error {
	*c.Logs = append(*c.Logs, &CapturedRecord{
		Parent: c.attrs,
		Record: &r,
	})
	return c.handler.Handle(ctx, r)
}

func (c *CapturingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CapturingHandler{
		handler: c.handler.WithAttrs(attrs),
		Logs:    c.Logs,
		attrs: &CapturedAttrs{
			Parent:     c.attrs,
			Attributes: attrs,
		},
	}
}

func (c *CapturingHandler) WithGroup(name string) slog.Handler {
	return &CapturingHandler{
		handler: c.handler.WithGroup(name),
		Logs:    c.Logs,
	}
}

func (c *CapturingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return c.handler.Enabled(ctx, level)
}

func (c *CapturingHandler) Clear() {
	*c.Logs = (*c.Logs)[:0] // reuse slice
}

func (c *CapturingHandler) FindLog(filters ...LogFilter) *CapturedRecord {
	for _, record := range *c.Logs {
		match := true
		for _, filter := range filters {
			if !filter(record) {
				match = false
				break
			}
		}
		if match {
			return record
		}
	}
	return nil
}

func (c *CapturingHandler) FindLogs(filters ...LogFilter) []*CapturedRecord {
	var logs []*CapturedRecord
	for _, record := range *c.Logs {
		match := true
		for _, filter := range filters {
			if !filter(record) {
				match = false
				break
			}
		}
		if match {
			logs = append(logs, record)
		}
	}
	return logs
}
