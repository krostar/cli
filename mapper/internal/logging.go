package mapper

import (
	"io"

	"github.com/krostar/test"
	testlogging "github.com/krostar/test/logging"
)

// TestWriter creates an io.Writer implementation suitable for testing CLI output.
// It wraps the provided test.TestingT with a prefixed writer to help identify
// CLI-specific output in test logs by adding a "[CLI]: " prefix to all written data.
//
// This is particularly useful when testing CLI applications where you need to:
// - Distinguish CLI output from other test output
// - Capture and verify command output in tests
// - Provide a controlled output destination for CLI commands
//
// The returned writer will forward all data to the test logging system while
// ensuring proper test helper marking for accurate test failure reporting.
func TestWriter(t test.TestingT) io.Writer {
	return &writerWithPrefix{
		t:      t,
		w:      testlogging.NewWriter(t),
		prefix: []byte("[CLI]: "),
	}
}

// writerWithPrefix is an io.Writer implementation that adds a consistent prefix
// to all written output before delegating to an underlying writer.
// This is particularly useful for differentiating CLI output in test logs.
type writerWithPrefix struct {
	t      test.TestingT
	w      io.Writer
	prefix []byte
}

func (pw *writerWithPrefix) Write(p []byte) (int, error) {
	pw.t.Helper()

	w, err := pw.w.Write(append(pw.prefix, p...))
	if w > len(p) {
		w = len(p)
	}

	return w, err
}
