package datasource

/*
Copyright (C) 2026 Rea Sand

This file is part of git-team.

This program is free software: you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation, either version 3
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"errors"
	"fmt"

	activationscope "github.com/hekmekk/git-team/v2/src/shared/activation/scope"
	config "github.com/hekmekk/git-team/v2/src/shared/config/entity/config"
	giterror "github.com/hekmekk/git-team/v2/src/shared/gitconfig/error"
	gitconfig "github.com/hekmekk/git-team/v2/src/shared/gitconfig/interface"
	gitconfigscope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
)

// GitconfigDataSource reads configuration from git config
type GitconfigDataSource struct {
	GitConfigReader gitconfig.Reader
}

// NewGitconfigDataSource constructs new GitconfigDataSource
func NewGitconfigDataSource(gitSettingsReader gitconfig.Reader) GitconfigDataSource {
	return GitconfigDataSource{gitSettingsReader}
}

func (ds GitconfigDataSource) Read() (config.Config, error) {
	rawScope, err := ds.GitConfigReader.Get(gitconfigscope.Global, "team.config.activation-scope")

	if err != nil && errors.Is(err, giterror.ErrSectionOrKeyIsInvalid) {
		return config.Config{ActivationScope: activationscope.Global}, nil
	}

	if err != nil {
		return config.Config{}, fmt.Errorf("failed to get team.config.activation-scope: %s", err)
	}

	scope := activationscope.FromString(rawScope)
	if scope == activationscope.Unknown {
		return config.Config{}, fmt.Errorf("unknown activation-scope '%s' found in config. Did you edit it manually?", rawScope)
	}

	cfg := config.Config{
		ActivationScope: scope,
	}

	return cfg, nil
}
