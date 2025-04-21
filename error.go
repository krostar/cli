package cli

// ShowHelpError is an interface that allows commands to signal whether
// the help message should be displayed along with the error.
type ShowHelpError interface {
	error
	ShowHelp() bool
}

// NewErrorWithHelp creates a new error that implements the ShowHelpError
// interface, indicating that the help message should be shown.
func NewErrorWithHelp(err error) error {
	return &showHelpError{err: err}
}

type showHelpError struct{ err error }

func (showHelpError) ShowHelp() bool  { return true }
func (e showHelpError) Unwrap() error { return e.err }
func (e showHelpError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return ""
}

// ExitStatusError allows commands to specify a custom exit
// status code for the CLI application.
type ExitStatusError interface {
	error
	ExitStatus() uint8
}

// NewErrorWithExitStatus creates a new error that implements the
// ExitStatusError interface, allowing a custom exit status to be set.
func NewErrorWithExitStatus(err error, status uint8) error {
	return &exitStatusError{err: err, status: status}
}

type exitStatusError struct {
	err    error
	status uint8
}

func (e exitStatusError) ExitStatus() uint8 { return e.status }
func (e exitStatusError) Unwrap() error     { return e.err }
func (e exitStatusError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return ""
}
