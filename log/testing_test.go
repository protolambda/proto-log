package log

import "testing"

func TestNewTestLogger(t *testing.T) {
	logger := TestLogger(t)
	logger.Debug("Hello world", "foo", 123, "bar", "hi")
	subLogger := logger.With("nest", "1")
	subLogger.Debug("Testing sub-logger")
	//subLogger.Crit("example") // To make the test fail
}
