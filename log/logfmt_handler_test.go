package log_test

import (
	"bytes"
	"testing"

	"github.com/protolambda/proto-log/log"
)

func TestLogfmtHandler(t *testing.T) {
	var buf bytes.Buffer
	h := log.LogfmtHandler(&buf)
	logger := log.New(h)
	logger.Debug("Hello world", "foo", "1", "bar", 2)
	got := string(buf.Bytes())
	t.Log(got)
	assertSubstring(t, got, `t=`) // ignore timestamp itself
	assertSubstring(t, got, `lvl=debug msg="Hello world" foo=1 bar=2`)
}
