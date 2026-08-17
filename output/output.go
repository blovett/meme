package output

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

var (
	// Stdout is a color friendly pipe.
	Stdout = colorable.NewColorableStdout()

	// Stderr is a color friendly pipe.
	Stderr = colorable.NewColorableStderr()
)

// WrapError wraps err with additional context text, returning nil if err is
// nil. Use it to annotate an error returned from a lower-level call before
// passing it back up the call stack.
func WrapError(err error, text string) error {
	if err != nil {
		return fmt.Errorf("%s: %w", text, err)
	}
	return nil
}

// Errorf returns a new error built from the given message.
func Errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// Info prints information.
func Info(format string, args ...interface{}) {
	fmt.Fprintf(Stdout, color.GreenString(format)+"\n", args...)
}
