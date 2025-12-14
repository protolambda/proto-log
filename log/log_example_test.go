package log_test

import (
	"context"
	"os"

	"github.com/protolambda/proto-log/log"
)

func ExampleWithIncludeSource() {
	workDir, _ := os.Getwd()
	h := log.LogfmtHandler(os.Stdout,
		log.WithExcludeTime(true),
		log.WithIncludeSource(true),
		log.WithSourceRelDir(workDir),
	)
	logger := log.New(h)
	logger.Info("Hello world", "foo", 1, "bar", true)
	// Output:
	// lvl=info source=log_example_test.go:18 msg="Hello world" foo=1 bar=true
}

func ExampleNew() {
	h := log.TerminalHandler(os.Stdout, log.WithColor(false), log.WithExcludeTime(true))
	logger := log.New(h)
	logger.Info("Hello world", "foo", 1, "bar", true)
	logger.InfoContext(context.Background(), "log with context", "hi", 123)
	// Output:
	// INFO  Hello world                              foo=1 bar=true
	// INFO  log with context                         hi=123
}

func ExampleJSONHandler() {
	h := log.JSONHandler(os.Stdout, log.WithExcludeTime(true))
	logger := log.New(h)
	logger.Info("Hello JSON", "example", map[string]int{"hello": 123})
	// Output: {"lvl":"info","msg":"Hello JSON","example":{"hello":123}}
}

func ExampleLogfmtHandler() {
	h := log.LogfmtHandler(os.Stdout, log.WithExcludeTime(true))
	logger := log.New(h)
	logger.Info("Hello Logfmt", "example", map[string]int{"hello": 123})
	// Output: lvl=info msg="Hello Logfmt" example=map[hello:123]
}

func ExampleLevelHandler_SetMinLevel() {
	h := log.TerminalHandler(os.Stdout,
		log.WithColor(false),
		log.WithExcludeTime(true),
	)
	logger := log.New(h,
		log.LevelMod(log.LevelInfo),
	)
	logger.Info("Hello info world", "foobar", 123)

	subLogger := logger.With("name", "alice")
	subLogger.Info("Report from sub-logger")
	// By getting the LevelHandler, the level of this logger can be adjusted
	lh, ok := log.FindHandler[*log.LevelHandler](subLogger.Handler())
	if !ok {
		panic("log handler does not have a LevelHandler mod")
	}
	lh.SetMinLevel(log.LevelDebug)

	logger.Debug("Hidden debug message")
	subLogger.Debug("Hello debug world from sub-logger")

	// Output:
	// INFO  Hello info world                         foobar=123
	// INFO  Report from sub-logger                   name=alice
	// DEBUG Hello debug world from sub-logger        name=alice
}
