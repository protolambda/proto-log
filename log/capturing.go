package log

import (
	"context"
	"log/slog"
	"slices"
	"sync"
)

type Capturer interface {
	slog.Handler
	Clear()
	FindLog(filters ...LogFilter) *CapturedRecord
	FindLogs(filters ...LogFilter) []*CapturedRecord
	Records() []*CapturedRecord
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

// CapturingHandler provides a log handler that captures all log records and forwards them to a delegate.
// It is safe for concurrent use: handlers derived through WithAttrs and WithGroup share the captured records.
type CapturingHandler struct {
	handler slog.Handler
	store   *captureStore // shared among derived CapturingHandlers
	// attrs are inherited log record attributes, from a logger that this CapturingHandler may be derived from
	attrs *CapturedAttrs
}

type captureStore struct {
	mu   sync.Mutex
	logs []*CapturedRecord
}

var _ Handler = (*CapturingHandler)(nil)

func CapturingMod() HandlerMod {
	return func(h slog.Handler) slog.Handler {
		return &CapturingHandler{handler: h, store: new(captureStore)}
	}
}

func (c *CapturingHandler) Unwrap() slog.Handler {
	return c.handler
}

func (c *CapturingHandler) Handle(ctx context.Context, r slog.Record) error {
	clone := r.Clone() // the caller may reuse the attrs storage after Handle returns
	rec := &CapturedRecord{Parent: c.attrs, Record: &clone}
	c.store.mu.Lock()
	c.store.logs = append(c.store.logs, rec)
	c.store.mu.Unlock()
	return c.handler.Handle(ctx, r)
}

func (c *CapturingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CapturingHandler{
		handler: c.handler.WithAttrs(attrs),
		store:   c.store,
		attrs: &CapturedAttrs{
			Parent:     c.attrs,
			Attributes: attrs,
		},
	}
}

// WithGroup derives a handler that keeps the inherited attributes.
// Captured attribute keys are not qualified by group names.
func (c *CapturingHandler) WithGroup(name string) slog.Handler {
	return &CapturingHandler{
		handler: c.handler.WithGroup(name),
		store:   c.store,
		attrs:   c.attrs,
	}
}

func (c *CapturingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return c.handler.Enabled(ctx, level)
}

// Records returns a snapshot of all captured records, in capture order.
func (c *CapturingHandler) Records() []*CapturedRecord {
	c.store.mu.Lock()
	defer c.store.mu.Unlock()
	return slices.Clone(c.store.logs)
}

func (c *CapturingHandler) Clear() {
	c.store.mu.Lock()
	defer c.store.mu.Unlock()
	c.store.logs = nil // don't reuse the slice: earlier Records snapshots may still reference it
}

func (c *CapturingHandler) FindLog(filters ...LogFilter) *CapturedRecord {
	for _, record := range c.Records() {
		if matchAll(record, filters) {
			return record
		}
	}
	return nil
}

func (c *CapturingHandler) FindLogs(filters ...LogFilter) []*CapturedRecord {
	var logs []*CapturedRecord
	for _, record := range c.Records() {
		if matchAll(record, filters) {
			logs = append(logs, record)
		}
	}
	return logs
}

func matchAll(record *CapturedRecord, filters []LogFilter) bool {
	for _, filter := range filters {
		if !filter(record) {
			return false
		}
	}
	return true
}
