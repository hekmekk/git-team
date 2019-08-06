package disable

import (
	"fmt"
	"github.com/hekmekk/git-team/src/config"
	git "github.com/hekmekk/git-team/src/gitconfig"
	statusApi "github.com/hekmekk/git-team/src/status"
	"os"
)

// TODO: Move Exec piece with dependency declaration (NOT definition) into diff file or even module

// Exec disable team mode. Remove commit section and unset commit template.
func Exec() error {
	deps := dependencies{
		GitUnsetCommitTemplate: git.UnsetCommitTemplate,
		GitUnsetHooksPath:      git.UnsetHooksPath,
		GitRemoveCommitSection: git.RemoveCommitSection,
		LoadConfig:             config.Load,
		RemoveFile:             os.Remove,
		PersistDisabled:        statusApi.PersistDisabled,
	}
	return executorFactory(deps)()
}

type dependencies struct {
	GitUnsetCommitTemplate func() error
	GitUnsetHooksPath      func() error
	GitRemoveCommitSection func() error
	LoadConfig             func() (config.Config, error)
	RemoveFile             func(string) error
	PersistDisabled        func() error
}

func executorFactory(deps dependencies) func() error {
	return func() error {
		if err := deps.GitUnsetHooksPath(); err != nil {
			return err
		}

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

		if err := deps.PersistDisabled(); err != nil {
			return err
		}

		return nil
	}
}
