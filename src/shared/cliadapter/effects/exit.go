package effects

import "fmt"

type ExitOk struct {
	message string
}

// Message return the message
func (exitOk ExitOk) Message() string {
	return exitOk.message
}

type ExitWarn struct {
	message string
}

// Message return the message
func (exitWarn ExitWarn) Message() string {
	return exitWarn.message
}

type ExitErr struct {
	message string
}

// Message return the message
func (exitErr ExitErr) Message() string {
	return exitErr.message
}

// NewExitOk exit with success code
func NewExitOk() Effect {
	return ExitOk{
		message: "",
	}
}

// NewExitOkMsg exit with success code and print a message
func NewExitOkMsg(message string) Effect {
	return ExitOk{
		message: message,
	}
}

// NewExitWarn exit with warn code and print a message prefixed with warn:
func NewExitWarn(message string) Effect {
	return ExitWarn{
		message: fmt.Sprintf("warn: %s", message),
	}
}

// NewExitErr exit with error code and print a message prefixed with error:
func NewExitErr(err error) Effect {
	return ExitErr{
		message: fmt.Sprintf("error: %s", err.Error()),
	}
}
