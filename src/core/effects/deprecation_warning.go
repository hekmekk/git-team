package effects

import (
	"os"

	"github.com/fatih/color"
)

// DeprecationWarning the type definition
type DeprecationWarning struct {
	deprecated string
	suggested  string
}

// NewDeprecationWarning the constructor
func NewDeprecationWarning(deprecated string, suggested string) DeprecationWarning {
	return DeprecationWarning{
		deprecated: deprecated,
		suggested:  suggested,
	}
}

// Run write a line to STDOUT
func (deprecationWarning DeprecationWarning) Run() {
	os.Stdout.WriteString(color.YellowString("warn: '%s' has been deprecated and is going to be removed in a future major release, use '%s' instead\n", deprecationWarning.deprecated, deprecationWarning.suggested))
}
