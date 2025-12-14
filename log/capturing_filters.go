package log

import (
	"log/slog"
	"strings"
)

type LogFilter func(record *CapturedRecord) bool

func LevelFilter(level slog.Level) LogFilter {
	return func(r *CapturedRecord) bool {
		return r.Record.Level == level
	}
}

func AttributesFilter(key, value string) LogFilter {
	return func(r *CapturedRecord) bool {
		found := false
		r.Attrs(func(a slog.Attr) bool {
			if a.Key == key && a.Value.String() == value {
				found = true
				return false
			}
			return true // try next
		})
		return found
	}
}

func AttributesContainsFilter(key, value string) LogFilter {
	return func(r *CapturedRecord) bool {
		found := false
		r.Attrs(func(a slog.Attr) bool {
			if a.Key == key && strings.Contains(a.Value.String(), value) {
				found = true
				return false
			}
			return true // try next
		})
		return found
	}
}

func MessageFilter(message string) LogFilter {
	return func(r *CapturedRecord) bool {
		return r.Record.Message == message
	}
}

func MessageContainsFilter(message string) LogFilter {
	return func(r *CapturedRecord) bool {
		return strings.Contains(r.Record.Message, message)
	}
}

func ErrContainsFilter(errMessage string) LogFilter {
	return func(r *CapturedRecord) bool {
		found := false
		r.Attrs(func(a slog.Attr) bool {
			if a.Key != "err" {
				return true
			}
			if err, ok := a.Value.Any().(error); ok && strings.Contains(err.Error(), errMessage) {
				found = true
				return false
			}
			return true
		})
		return found
	}
}
