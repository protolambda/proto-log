package log

import (
	"context"
	"io"
	"log/slog"
	"os"
)

type T interface {
	// Output returns a Writer that writes to the same test output stream as TB.Log.
	// The output is indented like TB.Log lines, but Output does not add source locations or newlines.
	// The output is internally line buffered, and a call to TB.Log or the end of the test will implicitly
	// flush the buffer, followed by a newline.
	// After a test function and all its parents return, neither Output nor the Write method may be called.
	//
	// This was introduced in Go 1.25
	Output() io.Writer

	Error(args ...any)
	FailNow()
}

// TestLogger creates a TerminalHandler configured for testing.
// All log-output is written to the T.Output().
// Color is enabled.
// Source-info is enabled.
// Crit-level logs will be followed up with a T.FailNow().
func TestLogger(t T, mods ...HandlerMod) Logger {
	var h slog.Handler
	wd, err := os.Getwd()
	if err != nil {
		t.Error("failed to get work-dir:", err)
		t.FailNow()
		return nil
	}
	h = TerminalHandler(t.Output(),
		WithColor(true),
		WithIncludeSource(true),
		WithSourceRelDir(wd))
	for _, m := range mods {
		h = m(h)
	}
	// FailNow after Crit-level logs
	h = PostProcessMod(func(ctx context.Context, r slog.Record) {
		if r.Level >= LevelCrit {
			t.Error("CRIT-level log")
			t.FailNow()
		}
	})(h)
	return New(h)
}
