package effects

import (
	"fmt"
	"os"
)

// PrintMessage the type definition
type PrintMessage struct {
	message string
}

// NewPrintMessage the constructor
func NewPrintMessage(msg string) PrintMessage {
	return PrintMessage{
		message: msg,
	}
}

// Run write a line to STDOUT
func (printMsg PrintMessage) Run() {
	os.Stdout.WriteString(fmt.Sprintf("%s\n", printMsg.message))
}
