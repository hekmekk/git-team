package handler

import (
	"fmt"
	statusRepository "github.com/hekmekk/git-team/core/status"
	"github.com/hekmekk/git-team/src/config"
	git "github.com/hekmekk/git-team/src/gitconfig"
	"os"
)

type Dependencies struct {
	LoadConfig             func() (config.Config, error)
	SaveStatus             func(state statusRepository.State, coauthors ...string) error
	GitUnsetCommitTemplate func() error
	GitRemoveCommitSection func() error
}

func Disable() error {
	deps := Dependencies{
		LoadConfig:             config.Load,
		SaveStatus:             statusRepository.Save,
		GitUnsetCommitTemplate: git.UnsetCommitTemplate,
		GitRemoveCommitSection: git.RemoveCommitSection,
	}
	return executorFactory(deps)()
}

func executorFactory(deps Dependencies) func() error {
	return func() error {
		if err := deps.GitUnsetCommitTemplate(); err != nil {
			return err
		}

		if err := deps.GitRemoveCommitSection(); err != nil {
			return err
		}

		cfg, err := deps.LoadConfig()
		if err != nil {
			return err
		}
		os.Remove(fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.TemplateFileName))

		if err := deps.SaveStatus(statusRepository.DISABLED); err != nil {
			return err
		}

		return nil
	}
}
