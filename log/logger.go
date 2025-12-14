// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"context"
	"log/slog"
	"runtime"
	"time"

	"github.com/protolambda/proto-log/log/internal/logsettings"
)

// A loggerImpl records structured information about each call to its
// Log, Debug, Info, Warn, and Error methods.
// For each call, it creates a [Record] and passes it to a [Handler].
//
// To create a new Logger, call [New] or a Logger method
// that begins "With".
type loggerImpl struct {
	handler slog.Handler // for structured logging
}

// New creates a new Logger with the given non-nil Handler.
func New(h slog.Handler, mods ...HandlerMod) Logger {
	if h == nil {
		panic("nil Handler")
	}
	for _, mod := range mods {
		h = mod(h)
	}
	if _, ok := FindHandler[*ContextHandler](h); !ok {
		// if there is no ContextHandler in the stack, add it
		h = ContextMod()(h)
	}
	return &loggerImpl{handler: h}
}

func (l *loggerImpl) clone() *loggerImpl {
	c := *l
	return &c
}

// Handler returns l's Handler.
func (l *loggerImpl) Handler() slog.Handler { return l.handler }

// With returns a Logger that includes the given attributes
// in each output operation. Arguments are converted to
// attributes as if by [Logger.Log].
func (l *loggerImpl) With(args ...any) Logger {
	if len(args) == 0 {
		return l
	}
	c := l.clone()
	c.handler = l.handler.WithAttrs(argsToAttrSlice(args))
	return c
}

func argsToAttrSlice(args []any) []slog.Attr {
	var (
		attr  slog.Attr
		attrs []slog.Attr
	)
	for len(args) > 0 {
		attr, args = argsToAttr(args)
		attrs = append(attrs, attr)
	}
	return attrs
}

const badKey = "!BADKEY"

// argsToAttr turns a prefix of the nonempty args slice into an Attr
// and returns the unconsumed portion of the slice.
// If args[0] is an Attr, it returns it.
// If args[0] is a string, it treats the first two elements as
// a key-value pair.
// Otherwise, it treats args[0] as a value with a missing key.
func argsToAttr(args []any) (slog.Attr, []any) {
	switch x := args[0].(type) {
	case string:
		if len(args) == 1 {
			return slog.String(badKey, x), nil
		}
		return slog.Any(x, args[1]), args[2:]

	case slog.Attr:
		return x, args[1:]

	default:
		return slog.Any(badKey, x), args[1:]
	}
}

// WithGroup returns a Logger that starts a group, if name is non-empty.
// The keys of all attributes added to the Logger will be qualified by the given
// name. (How that qualification happens depends on the [Handler.WithGroup]
// method of the Logger's Handler.)
//
// If name is empty, WithGroup returns the receiver.
func (l *loggerImpl) WithGroup(name string) Logger {
	if name == "" {
		return l
	}
	c := l.clone()
	c.handler = l.handler.WithGroup(name)
	return c
}

// Enabled reports whether l emits log records at the given context and level.
func (l *loggerImpl) Enabled(ctx context.Context, level slog.Level) bool {
	if ctx == nil {
		ctx = context.Background()
	}
	return l.Handler().Enabled(ctx, level)
}

// Log emits a log record with the current time and the given level and message.
// The Record's Attrs consist of the Logger's attributes followed by
// the Attrs specified by args.
//
// The attribute arguments are processed as follows:
//   - If an argument is an Attr, it is used as is.
//   - If an argument is a string and this is not the last argument,
//     the following argument is treated as the value and the two are combined
//     into an Attr.
//   - Otherwise, the argument is treated as a value with key "!BADKEY".
func (l *loggerImpl) Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	l.log(ctx, level, msg, args...)
}

// LogAttrs is a more efficient version of [Logger.Log] that accepts only Attrs.
func (l *loggerImpl) LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	l.logAttrs(ctx, level, msg, attrs...)
}

// Trace logs at [LevelTrace].
func (l *loggerImpl) Trace(msg string, args ...any) {
	l.log(context.Background(), LevelTrace, msg, args...)
}

// TraceContext logs at [LevelTrace] with the given context.
func (l *loggerImpl) TraceContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelTrace, msg, args...)
}

// Debug logs at [LevelDebug].
func (l *loggerImpl) Debug(msg string, args ...any) {
	l.log(context.Background(), LevelDebug, msg, args...)
}

// DebugContext logs at [LevelDebug] with the given context.
func (l *loggerImpl) DebugContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelDebug, msg, args...)
}

// Info logs at [LevelInfo].
func (l *loggerImpl) Info(msg string, args ...any) {
	l.log(context.Background(), LevelInfo, msg, args...)
}

// InfoContext logs at [LevelInfo] with the given context.
func (l *loggerImpl) InfoContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelInfo, msg, args...)
}

// Warn logs at [LevelWarn].
func (l *loggerImpl) Warn(msg string, args ...any) {
	l.log(context.Background(), LevelWarn, msg, args...)
}

// WarnContext logs at [LevelWarn] with the given context.
func (l *loggerImpl) WarnContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelWarn, msg, args...)
}

// Error logs at [LevelError].
func (l *loggerImpl) Error(msg string, args ...any) {
	l.log(context.Background(), LevelError, msg, args...)
}

// ErrorContext logs at [LevelError] with the given context.
func (l *loggerImpl) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelError, msg, args...)
}

// Crit logs at [LevelCrit].
func (l *loggerImpl) Crit(msg string, args ...any) {
	l.log(context.Background(), LevelCrit, msg, args...)
}

// CritContext logs at [LevelCrit] with the given context.
func (l *loggerImpl) CritContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, LevelCrit, msg, args...)
}

// log is the low-level logging method for methods that take ...any.
// It must always be called directly by an exported logging method
// or function, because it uses a fixed call depth to obtain the pc.
func (l *loggerImpl) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if !l.Enabled(ctx, level) {
		return
	}
	var pc uintptr
	if !logsettings.IgnorePC {
		var pcs [1]uintptr
		// skip [runtime.Callers, this function, this function's caller]
		runtime.Callers(3, pcs[:])
		pc = pcs[0]
	}
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = l.Handler().Handle(ctx, r)
}

// logAttrs is like [Logger.log], but for methods that take ...Attr.
func (l *loggerImpl) logAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	if !l.Enabled(ctx, level) {
		return
	}
	var pc uintptr
	if !logsettings.IgnorePC {
		var pcs [1]uintptr
		// skip [runtime.Callers, this function, this function's caller]
		runtime.Callers(3, pcs[:])
		pc = pcs[0]
	}
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.AddAttrs(attrs...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = l.Handler().Handle(ctx, r)
}

// Context returns the default context that is used when logging
func (l *loggerImpl) Context() context.Context {
	h, ok := FindHandler[*ContextHandler](l.handler)
	if !ok {
		return context.Background()
	}
	return h.Context()
}

// WithContext creates a clone, with the given context as new default context.
func (l *loggerImpl) WithContext(ctx context.Context) Logger {
	c := l.clone()
	c.handler = l.handler.WithAttrs(nil)
	h, ok := FindHandler[*ContextHandler](c.handler)
	if !ok {
		panic("expected context handler")
	}
	h.SetContext(ctx)
	return c
}
