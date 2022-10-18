package script

import (
	_ "embed"
)

var (
	//go:embed bash.sh
	Bash string
	//go:embed zsh.sh
	Zsh string
)
