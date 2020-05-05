package cli

import (
	"errors"
	"io"
	"os"
)

type ShowHelpError interface {
	error
	Unwrap() error
	ShowHelp() bool
}

type showHelpError struct{ error }

func (she showHelpError) ShowHelp() bool { return true }
func (she showHelpError) Unwrap() error  { return she.error }

func ErrorShowHelp(err error) error { return &showHelpError{err} }

type ExitStatusError interface {
	error
	ExitStatus() int
}

type exitStatusError struct {
	error
	status int
}

func (ese exitStatusError) ExitStatus() int { return ese.status }
func (ese exitStatusError) Unwrap() error   { return ese.error }

func ErrorWithExitStatus(err error, status int) error {
	return &exitStatusError{error: err, status: status}
}

func PrintErrorIfAnyAndExit(writer io.Writer, err error) {
	if err == nil {
		os.Exit(0)
	}

	_, _ = io.WriteString(writer, err.Error()+"\n")
	if closer, ok := writer.(io.Closer); ok {
		_ = closer.Close()
	}

	var errWithStatus ExitStatusError
	if errors.As(err, &errWithStatus) {
		os.Exit(errWithStatus.ExitStatus())
	}

	os.Exit(1)
}
