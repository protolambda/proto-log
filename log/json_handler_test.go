package log_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/protolambda/proto-log/log"
)

func TestJSONHandler(t *testing.T) {
	wd, _ := os.Getwd()
	var buf bytes.Buffer
	h := log.JSONHandler(&buf, log.WithExcludeTime(true), log.WithIncludeSource(true), log.WithSourceRelDir(wd))
	logger := log.New(h)
	logger.Debug("Hello world", "foo", "1", "bar", 2)
	got := string(buf.Bytes())
	got = strings.TrimSpace(got)
	t.Log(got)
	expected := `{"lvl":"debug","source":{"function":"github.com/protolambda/proto-log/log_test.TestJSONHandler","file":"json_handler_test.go","line":17},"msg":"Hello world","foo":"1","bar":2}`
	assertEqual(t, got, expected)
}
