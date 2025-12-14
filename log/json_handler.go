package log

import (
	"io"
	"log/slog"
)

// JSONHandler returns a handler which prints records in JSON format.
func JSONHandler(wr io.Writer, opts ...FormatOption) slog.Handler {
	var cfg FormatConfig
	cfg.Apply(opts...)
	hOpts := &slog.HandlerOptions{
		AddSource: cfg.IncludeSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return cfg.BuiltinReplace(groups, a, false)
		},
		Level: LevelMaxVerbosity,
	}
	return slog.NewJSONHandler(wr, hOpts)
}
