package effects

import (
	"os"
)

// Exit the type definition
type Exit struct {
	code int
}

// NewExitOk the constructor for the ok case
func NewExitOk() Exit {
	return Exit{
		code: 0,
	}
}

// NewExitErr the constructor for the error case
func NewExitErr() Exit {
	return Exit{
		code: -1,
	}
}

// Run exit with the respective status code
func (exit Exit) Run() {
	os.Exit(exit.code)
}
