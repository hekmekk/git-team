package effects

import (
	"os"
)

// ExitOk the type definition
type ExitOk struct{}

// NewExitOk the constructor
func NewExitOk() ExitOk {
	return ExitOk{}
}

// Run exit with status code 0
func (ExitOk) Run() {
	os.Exit(0)
}
