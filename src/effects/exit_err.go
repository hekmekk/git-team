package effects

import (
	"os"
)

// ExitErr the type definition
type ExitErr struct{}

// NewExitErr the constructor
func NewExitErr() ExitErr {
	return ExitErr{}
}

// Run exit with status code 255
func (ExitErr) Run() {
	os.Exit(-1)
}
