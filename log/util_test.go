package log_test

import (
	"strings"
	"testing"
)

func assertTrue(t *testing.T, v bool) {
	if !v {
		t.Helper()
		t.Error("expected true")
		t.FailNow()
	}
}

func assertNotNil[V any](t *testing.T, v *V) {
	if v == nil {
		t.Helper()
		t.Error("expected non-nil value")
		t.FailNow()
	}
}

func assertNil[V any](t *testing.T, v *V) {
	if v != nil {
		t.Helper()
		t.Error("expected nil value, but got:", v)
		t.FailNow()
	}
}

func assertEqual[V comparable](t *testing.T, a V, b V) {
	if a != b {
		t.Helper()
		t.Errorf("expected to be equal:\nA: %v\nB: %v", a, b)
		t.FailNow()
	}
}

func assertSubstring(t *testing.T, v string, sub string) {
	if !strings.Contains(v, sub) {
		t.Helper()
		t.Errorf("expected %q to be substring of %q", sub, v)
		t.FailNow()
	}
}
