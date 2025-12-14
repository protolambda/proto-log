package log

import (
	"fmt"
	"log/slog"
	"math/big"
	"path/filepath"
	"reflect"
	"time"
)

func (cfg *FormatConfig) BuiltinReplace(_ []string, attr slog.Attr, logfmt bool) slog.Attr {
	switch attr.Key {
	case slog.TimeKey:
		if cfg.ExcludeTime {
			return slog.Attr{}
		}
		if attr.Value.Kind() == slog.KindTime {
			if logfmt {
				return slog.String("t", attr.Value.Time().Format(timeFormat))
			} else {
				return slog.Attr{Key: "t", Value: attr.Value}
			}
		}
	case slog.LevelKey:
		if l, ok := attr.Value.Any().(slog.Level); ok {
			attr = slog.Any("lvl", LevelString(l))
			return attr
		}
	case slog.SourceKey:
		if !cfg.IncludeSource {
			return slog.Attr{}
		} else if cfg.SourceRelDir != "" {
			// Hacky, to access and mutate the inner Source data, for adjustment of the filepath
			s := attr.Value.Resolve().Any().(*slog.Source)
			file, err := filepath.Rel(cfg.SourceRelDir, s.File)
			if err == nil {
				s.File = file
			}
		}
	}

	switch v := attr.Value.Any().(type) {
	case time.Time:
		if logfmt {
			attr = slog.String(attr.Key, v.Format(timeFormat))
		}
	case *big.Int:
		if v == nil {
			attr.Value = slog.StringValue("<nil>")
		} else {
			attr.Value = slog.StringValue(v.String())
		}
	case u256:
		if v == nil {
			attr.Value = slog.StringValue("<nil>")
		} else {
			attr.Value = slog.StringValue(v.Dec())
		}
	case fmt.Stringer:
		if v == nil || (reflect.ValueOf(v).Kind() == reflect.Pointer && reflect.ValueOf(v).IsNil()) {
			attr.Value = slog.StringValue("<nil>")
		} else {
			attr.Value = slog.StringValue(v.String())
		}
	}
	return attr
}
