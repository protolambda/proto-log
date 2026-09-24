package log_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/protolambda/proto-log/log"
)

func TestCapturingHandlerConcurrent(t *testing.T) {
	lgr := log.New(log.TerminalHandler(io.Discard), log.CapturingMod())
	logs, ok := log.FindHandler[log.Capturer](lgr.Handler())
	assertTrue(t, ok)

	const workers, perWorker = 8, 100
	var wg sync.WaitGroup
	for i := range workers {
		sub := lgr.With("worker", i)
		wg.Go(func() {
			for range perWorker {
				sub.Error("work")
				_ = logs.FindLogs(log.MessageFilter("work"))
			}
		})
	}
	wg.Wait()
	assertEqual(t, workers*perWorker, len(logs.Records()))
	logs.Clear()
	assertEqual(t, 0, len(logs.Records()))
}

func TestCapturingHandlerWithGroupKeepsAttrs(t *testing.T) {
	lgr := log.New(log.TerminalHandler(io.Discard), log.CapturingMod())
	logs, ok := log.FindHandler[log.Capturer](lgr.Handler())
	assertTrue(t, ok)
	lgr.With("a", "x").WithGroup("g").Error("hi")
	assertNotNil(t, logs.FindLog(log.AttributesFilter("a", "x")))
}

func TestTerminalHandlerConcurrent(t *testing.T) {
	var buf bytes.Buffer // not safe for concurrent use: relies on the shared handler lock
	lgr := log.New(log.TerminalHandler(&buf, log.WithExcludeTime(true)))
	var wg sync.WaitGroup
	for i := range 8 {
		sub := lgr.With("worker", i)
		wg.Go(func() {
			for range 100 {
				sub.Info("work", "k", "v")
			}
		})
	}
	wg.Wait()
	lines := strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n")
	assertEqual(t, 800, len(lines))
	for _, line := range lines {
		assertSubstring(t, line, "k=v")
	}
}

func TestLevelSharedWithDerived(t *testing.T) {
	lgr := log.New(log.DiscardHandler(), log.LevelMod(log.LevelInfo))
	sub := lgr.With("a", 1).WithGroup("g")
	assertTrue(t, !sub.Enabled(context.Background(), log.LevelDebug))
	lh, ok := log.FindHandler[*log.LevelHandler](lgr.Handler())
	assertTrue(t, ok)
	lh.SetMinLevel(log.LevelDebug)
	subLH, ok := log.FindHandler[*log.LevelHandler](sub.Handler())
	assertTrue(t, ok)
	assertEqual(t, log.LevelDebug, subLH.MinLevel())
}

type ctxKey struct{}

func TestWithContext(t *testing.T) {
	var got []context.Context
	lgr := log.New(log.TerminalHandler(io.Discard), log.PostProcessMod(func(ctx context.Context, _ slog.Record) {
		got = append(got, ctx)
	}))
	ctx := context.WithValue(context.Background(), ctxKey{}, "a")
	withCtx := lgr.WithContext(ctx)
	assertTrue(t, lgr.Context() == context.Background()) // parent is not mutated
	assertTrue(t, withCtx.Context() == ctx)
	again := withCtx.WithContext(context.TODO())
	assertTrue(t, again.Context() == context.TODO())
	assertTrue(t, withCtx.Context() == ctx)

	withCtx.Info("default")
	withCtx.InfoContext(nil, "nil")
	withCtx.InfoContext(context.TODO(), "todo")
	assertEqual(t, 3, len(got))
	assertTrue(t, got[0] == ctx)
	assertTrue(t, got[1] == ctx)
	assertTrue(t, got[2] == context.TODO())
}

type errWriter struct{}

func (errWriter) Write([]byte) (int, error) { return 0, errors.New("broken") }

func TestTerminalHandlerReturnsWriteError(t *testing.T) {
	h := log.TerminalHandler(errWriter{})
	err := h.Handle(context.Background(), slog.NewRecord(time.Time{}, log.LevelInfo, "x", 0))
	assertTrue(t, err != nil)
}
