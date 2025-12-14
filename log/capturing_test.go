package log_test

import (
	"testing"

	"github.com/protolambda/proto-log/log"
)

func TestCaptureLogger(t *testing.T) {
	lgr := log.TestLogger(t, log.CapturingMod())
	logs, ok := log.FindHandler[log.Capturer](lgr.Handler())
	assertTrue(t, ok)

	msg := "foo bar"
	lgr.Info(msg, "a", 1)
	msgFilter := log.MessageFilter(msg)
	rec := logs.FindLog(msgFilter)
	assertEqual(t, msg, rec.Record.Message)
	assertEqual(t, 1, rec.AttrValue("a").(int64))

	lgr.Debug("bug")
	containsFilter := log.MessageContainsFilter("bug")
	l := logs.FindLog(containsFilter)
	assertNotNil(t, l) // should capture all logs, not only above level

	msgClear := "clear"
	lgr.Error(msgClear)
	levelFilter := log.LevelFilter(log.LevelError)
	msgFilter = log.MessageFilter(msgClear)
	assertNotNil(t, logs.FindLog(levelFilter, msgFilter))
	logs.Clear()
	containsFilter = log.MessageContainsFilter(msgClear)
	l = logs.FindLog(containsFilter)
	assertNil(t, l)

	lgrb := lgr.With("b", 2)
	msgOp := "optimistic"
	lgrb.Info(msgOp, "c", 3)
	containsFilter = log.MessageContainsFilter(msgOp)
	recOp := logs.FindLog(containsFilter)
	assertNotNil(t, recOp) // should still capture logs from derived logger
	assertEqual(t, 3, recOp.AttrValue("c").(int64))
	// Note: "b" attributes won't be visible on captured record
}

func TestCaptureLoggerAttributesFilter(t *testing.T) {
	lgr := log.TestLogger(t, log.CapturingMod())
	logs, ok := log.FindHandler[log.Capturer](lgr.Handler())
	assertTrue(t, ok)

	msg := "foo bar"
	lgr.Info(msg, "a", "test")
	lgr.Info(msg, "a", "test 2")
	lgr.Info(msg, "a", "random")
	msgFilter := log.MessageFilter(msg)
	attrFilter := log.AttributesFilter("a", "random")

	rec := logs.FindLog(msgFilter, attrFilter)
	assertEqual(t, msg, rec.Record.Message)
	assertEqual(t, "random", rec.AttrValue("a").(string))

	recs := logs.FindLogs(msgFilter, attrFilter)
	assertEqual(t, len(recs), 1)
}

func TestCaptureLoggerNested(t *testing.T) {
	lgrInner := log.TestLogger(t, log.CapturingMod())
	logs, ok := log.FindHandler[log.Capturer](lgrInner.Handler())
	assertTrue(t, ok)

	lgrInner.Info("hi", "a", "test")

	lgrChildX := lgrInner.With("name", "childX")
	lgrChildX.Info("hello", "b", "42")

	lgrChildY := lgrInner.With("name", "childY")
	lgrChildY.Info("hola", "c", "7")

	lgrInner.Info("hello universe", "greeting", "from Inner")

	lgrChildX.Info("hello world", "greeting", "from X")

	assertEqual(t, len(logs.FindLogs(log.AttributesFilter("name", "childX"))), 2) // X logged twice
	assertEqual(t, len(logs.FindLogs(log.AttributesFilter("name", "childY"))), 1) // Y logged once

	assertEqual(t, len(logs.FindLogs(
		log.AttributesContainsFilter("greeting", "from"))), 2) // two greetings
	assertEqual(t, len(logs.FindLogs(
		log.AttributesContainsFilter("greeting", "from"),
		log.AttributesFilter("name", "childX"))), 1) // only one greeting from X

	assertEqual(t, len(logs.FindLogs(
		log.AttributesFilter("a", "test"))), 1) // root logger logged 'a' once
}
