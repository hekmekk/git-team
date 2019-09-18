package effects

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// PrintErr the type definition
type PrintErr struct {
	err error
}

// NewPrintErr the constructor
func NewPrintErr(err error) PrintErr {
	return PrintErr{
		err: err,
	}
}

// Run write a line to STDERR
func (printErr PrintErr) Run() {
	os.Stderr.WriteString(color.RedString(fmt.Sprintf("error: %s\n", printErr.err)))
}
