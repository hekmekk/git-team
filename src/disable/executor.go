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
		LoadConfig:             config.Load,
		StatFile:               os.Stat,
		RemoveFile:             os.Remove,
		PersistDisabled:        statusApi.PersistDisabled,
	}
	return executorFactory(deps)()
}

type dependencies struct {
	GitUnsetCommitTemplate func() error
	LoadConfig             func() (config.Config, error)
	StatFile               func(string) (os.FileInfo, error)
	RemoveFile             func(string) error
	PersistDisabled        func() error
}

func executorFactory(deps dependencies) func() error {
	return func() error {
		if err := deps.GitUnsetCommitTemplate(); err != nil && err.Error() != "exit status 5" {
			return err
		}

		cfg, err := deps.LoadConfig()
		if err != nil {
			return err
		}

		templateFilePath := fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.TemplateFileName)

		if _, err := deps.StatFile(templateFilePath); err == nil {
			if err := deps.RemoveFile(templateFilePath); err != nil {
				return err
			}
		}

		if err := deps.PersistDisabled(); err != nil {
			return err
		}

		return nil
	}
}
