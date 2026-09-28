package script

// Copyright (C) 2026 Rea Sand
// Licensed under the EUPL

import (
	_ "embed"
)

var (
	//go:embed bash.sh
	Bash string
	//go:embed zsh.sh
	Zsh string
)
