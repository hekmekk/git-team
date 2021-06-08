package effects

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

type exitType int

const (
	Ok exitType = iota
	Error
)

type ExitWithoutMsg struct {
	kind exitType
}

type ExitWithMsg struct {
	kind    exitType
	message string
}

func (exit ExitWithoutMsg) Run() error {
	switch exit.kind {
	case Ok:
		return nil
	case Error:
		return cli.Exit("", 1)
	}
	return nil
}

func (exit ExitWithMsg) Run() error {
	switch exit.kind {
	case Ok:
		fmt.Println(exit.message)
		return nil
	case Error:
		return cli.Exit(exit.message, 1)
	default:
		return NewExitErrMsg(errors.New("unexpected behavior encountered")).Run()
	}
}

// NewExitOk exit with success code
func NewExitOk() Effect {
	return ExitWithoutMsg{
		kind: Ok,
	}
}

// NewExitOkMsg exit with success code and print a message
func NewExitOkMsg(message string) Effect {
	return ExitWithMsg{
		kind:    Ok,
		message: message,
	}
}

// NewExitErr exit with error code
func NewExitErr() Effect {
	return ExitWithoutMsg{
		kind: Error,
	}
}

// NewExitErrMsg exit with error code and print a red colored message with error:
func NewExitErrMsg(err error) Effect {
	return ExitWithMsg{
		kind:    Error,
		message: color.RedString(fmt.Sprintf("error: %s", err)),
	}
}
