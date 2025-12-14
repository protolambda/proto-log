package log

import (
	"fmt"
	"log/slog"
	"math"
	"strings"
)

const (
	LevelMaxVerbosity slog.Level = math.MinInt
	LevelTrace        slog.Level = -8
	LevelDebug                   = slog.LevelDebug
	LevelInfo                    = slog.LevelInfo
	LevelWarn                    = slog.LevelWarn
	LevelError                   = slog.LevelError
	LevelCrit         slog.Level = 12
)

// LevelFromString returns the implied slog.Level from a string name.
// This is case-insensitive, and allows log-level aliases.
// This does not handle +/- integer level suffixes.
func LevelFromString(lvlString string) (slog.Level, error) {
	lvlString = strings.ToLower(lvlString) // ignore case
	switch lvlString {
	case "trace", "trce":
		return LevelTrace, nil
	case "debug", "dbug", "dbg":
		return LevelDebug, nil
	case "info", "inf":
		return LevelInfo, nil
	case "warn", "wrn":
		return LevelWarn, nil
	case "error", "eror", "err":
		return LevelError, nil
	case "crit":
		return LevelCrit, nil
	default:
		return LevelDebug, fmt.Errorf("unknown level: %q", lvlString)
	}
}

// LevelAlignedString returns a 5-character string containing the name of a Lvl.
func LevelAlignedString(l slog.Level) string {
	switch l {
	case LevelTrace:
		return "TRACE"
	case slog.LevelDebug:
		return "DEBUG"
	case slog.LevelInfo:
		return "INFO "
	case slog.LevelWarn:
		return "WARN "
	case slog.LevelError:
		return "ERROR"
	case LevelCrit:
		return "CRIT "
	default:
		return "unknown level"
	}
}

// LevelString returns a string containing the name of a Lvl.
func LevelString(l slog.Level) string {
	switch l {
	case LevelTrace:
		return "trace"
	case slog.LevelDebug:
		return "debug"
	case slog.LevelInfo:
		return "info"
	case slog.LevelWarn:
		return "warn"
	case slog.LevelError:
		return "error"
	case LevelCrit:
		return "crit"
	default:
		return "unknown"
	}
}
