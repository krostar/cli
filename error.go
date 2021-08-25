package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
)

// ShowHelpError defines a new type of error that can show the CLI help.
type ShowHelpError interface {
	error
	ShowHelp() bool
}

type showHelpError struct{ err error }

func (she showHelpError) ShowHelp() bool { return true }
func (she showHelpError) Unwrap() error  { return she.err }
func (she showHelpError) Error() string {
	if she.err != nil {
		return "help requested: " + she.err.Error()
	}
	return ""
}

// ErrorShowHelp wraps provided error and tells the CLI to show usage help.
func ErrorShowHelp(err error) error { return &showHelpError{err: err} }

// ExitStatusError defines a new type of error that can change the CLI exit status.
type ExitStatusError interface {
	error
	ExitStatus() uint8
}

type exitStatusError struct {
	error
	status uint8
}

func (ese exitStatusError) ExitStatus() uint8 { return ese.status }
func (ese exitStatusError) Unwrap() error     { return ese.error }

// ErrorWithExitStatus wraps the provided error and tells the CLI to exit with provided code.
func ErrorWithExitStatus(err error, status uint8) error {
	return &exitStatusError{error: err, status: status}
}

// Exit exits the program and uses provided error to define program success or failure.
func Exit(ctx context.Context, err error) {
	var (
		msg    string
		status uint8
	)

	if err != nil {
		var errWithStatus ExitStatusError
		if errors.As(err, &errWithStatus) {
			status = errWithStatus.ExitStatus()
		} else {
			status = 125
		}
		msg = err.Error()
	}

	if msg != "" {
		writer := getExitLogger(ctx)
		if _, err := io.WriteString(writer, msg+"\n"); err != nil {
			fmt.Printf("unable to write program exit message: %v", err)
		}
		if err := writer.Close(); err != nil {
			fmt.Printf("unable to close writer: %v", err)
		}
	}

	os.Exit(int(status))
}
