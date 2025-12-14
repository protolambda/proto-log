package log_test

import (
	"log/slog"
	"testing"

	"github.com/protolambda/proto-log/log"
)

type handlerA struct {
	slog.Handler
}

func wrapA(inner slog.Handler) slog.Handler {
	return &handlerA{Handler: inner}
}

type handlerB struct {
	slog.Handler
}

func wrapB(inner slog.Handler) slog.Handler {
	return &handlerB{Handler: inner}
}

func (w *handlerB) Unwrap() slog.Handler {
	return w.Handler
}

type handlerC struct {
	slog.Handler
}

func wrapC(inner slog.Handler) slog.Handler {
	return &handlerC{Handler: inner}
}

func (w *handlerC) Unwrap() slog.Handler {
	return w.Handler
}

func TestFindHandler(t *testing.T) {
	t.Run("nested", func(t *testing.T) {
		a := wrapA(nil)
		b := wrapB(a)
		c := wrapC(b)
		h := c
		got1, ok := log.FindHandler[*handlerA](h)
		assertTrue(t, ok)
		assertEqual(t, a.(*handlerA), got1)
		got2, ok := log.FindHandler[*handlerB](h)
		assertTrue(t, ok)
		assertEqual(t, b.(*handlerB), got2)
		got3, ok := log.FindHandler[*handlerC](h)
		assertTrue(t, ok)
		assertEqual(t, c.(*handlerC), got3)
	})
}
