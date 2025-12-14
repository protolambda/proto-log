package log_test

import (
	"bytes"
	"testing"

	"github.com/protolambda/proto-log/log"
)

func TestTerminalHandler(t *testing.T) {
	var buf bytes.Buffer
	h := log.TerminalHandler(&buf, log.WithColor(false))
	logger := log.New(h)
	logger.Debug("Hello world", "foo", "1", "bar", 2)
	got := string(buf.Bytes())
	t.Log(got)
	assertSubstring(t, got, `DEBUG`)
	assertSubstring(t, got, `Hello world`)
	assertSubstring(t, got, `foo=1`)
	assertSubstring(t, got, `bar=2`)
}
