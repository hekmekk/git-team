package disable

import (
	"fmt"
	statusRepository "github.com/hekmekk/git-team/core/status"
	"github.com/hekmekk/git-team/src/config"
	git "github.com/hekmekk/git-team/src/gitconfig"
	"os"
)

// TODO: Move Exec piece with dependency declaration (NOT definition) into diff file or even module

// Exec disable team mode. Remove commit section and unset commit template.
func Exec() error {
	deps := dependencies{
		GitUnsetCommitTemplate: git.UnsetCommitTemplate,
		GitRemoveCommitSection: git.RemoveCommitSection,
		LoadConfig:             config.Load,
		RemoveFile:             os.Remove,
		SaveStatus:             statusRepository.Save,
	}
	return executorFactory(deps)()
}

type dependencies struct {
	GitUnsetCommitTemplate func() error
	GitRemoveCommitSection func() error
	LoadConfig             func() (config.Config, error)
	RemoveFile             func(string) error
	SaveStatus             func(state statusRepository.State, coauthors ...string) error
}

func executorFactory(deps dependencies) func() error {
	return func() error {
		if err := deps.GitUnsetCommitTemplate(); err != nil {
			return nil
		}

		if err := deps.GitRemoveCommitSection(); err != nil {
			return err
		}

		cfg, err := deps.LoadConfig()
		if err != nil {
			return err
		}

		if err := deps.RemoveFile(fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.TemplateFileName)); err != nil {
			return err
		}

		if err := deps.SaveStatus(statusRepository.DISABLED); err != nil {
			return err
		}

		return nil
	}
}
