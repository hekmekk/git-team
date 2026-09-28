package hookscript

// Copyright (C) 2026 Rea Sand
// Licensed under the EUPL

import (
	_ "embed"
)

var (
	//go:embed proxy.sh
	Proxy string

	//go:embed prepare-commit-msg.sh
	PrepareCommitMsg string

	//go:embed prepare-commit-msg-git-team.sh
	PrepareCommitMsgGitTeam string
)
