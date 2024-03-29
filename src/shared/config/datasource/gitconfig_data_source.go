package datasource

import (
	"errors"
	"fmt"

	activationscope "github.com/hekmekk/git-team/src/shared/activation/scope"
	config "github.com/hekmekk/git-team/src/shared/config/entity/config"
	giterror "github.com/hekmekk/git-team/src/shared/gitconfig/error"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/interface"
	gitconfigscope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
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
