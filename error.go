package cli

// ShowHelpError defines a new type of error that defines whenever the error should display the command help before the error message.
type ShowHelpError interface {
	error
	ShowHelp() bool
}

// NewErrorWithHelp wraps provided error and tells the CLI to show usage help.
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

// ExitStatusError defines a new type of error that allow the customization of the CLI exit status.
type ExitStatusError interface {
	error
	ExitStatus() uint8
}

// NewErrorWithExitStatus wraps the provided error and tells the CLI to exit with provided code.
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
