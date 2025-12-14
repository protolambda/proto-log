package log

import (
	"io"
	"log/slog"
)

// LogfmtHandler returns a handler which prints records in logfmt format, an easy machine-parseable but human-readable
// format for key/value pairs.
//
// For more details see: http://godoc.org/github.com/kr/logfmt
func LogfmtHandler(wr io.Writer, opts ...FormatOption) slog.Handler {
	var cfg FormatConfig
	cfg.Apply(opts...)
	hOpts := &slog.HandlerOptions{
		AddSource: cfg.IncludeSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return cfg.BuiltinReplace(groups, a, true)
		},
		Level: LevelMaxVerbosity,
	}
	return slog.NewTextHandler(wr, hOpts)
}
