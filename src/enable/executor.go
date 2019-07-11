package enable

import (
	"fmt"
	"os"

	"github.com/hekmekk/git-team/core/status"
)

type EnableCommand struct {
	Coauthors        []string
	BaseDir          string
	TemplateFileName string
}

type EnableEffect struct {
	CreateDir         func(path string, perm os.FileMode) error
	WriteFile         func(path string, data []byte, mode os.FileMode) error
	SetCommitTemplate func(path string) error
	SaveStatus        func(state status.State, coauthors ...string) error
}

func ExecutorFactory(effect EnableEffect) func(cmd EnableCommand) error {
	return func(cmd EnableCommand) error {
		if len(cmd.Coauthors) == 0 {
			return nil
		}

		coauthorsString := PrepareForCommitMessage(cmd.Coauthors)

		err := effect.CreateDir(cmd.BaseDir, os.ModePerm)
		if err != nil {
			return err
		}

		templatePath := fmt.Sprintf("%s/%s", cmd.BaseDir, cmd.TemplateFileName)

		err = effect.WriteFile(templatePath, []byte(coauthorsString), 0644)
		if err != nil {
			return err
		}
		err = effect.SetCommitTemplate(templatePath)
		if err != nil {
			return err
		}
		err = effect.SaveStatus(status.ENABLED, cmd.Coauthors...)
		if err != nil {
			return err
		}
		return nil
	}
}
