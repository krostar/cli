package cli

import (
	"context"
	"errors"
	"io"
	"os"
)

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

func ErrorShowHelp(err error) error { return &showHelpError{err: err} }

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

func Exit(ctx context.Context, err error) {
	var (
		msg    string
		status int
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
		_, _ = io.WriteString(writer, msg+"\n")
		if closer, ok := writer.(io.Closer); ok {
			_ = closer.Close()
		}
	}

	os.Exit(status)
}
