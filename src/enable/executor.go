package enable

import (
	"fmt"
	"os"

	"github.com/hekmekk/git-team/core/status"
)

// Command add a <Coauthor> under "team.alias.<Alias>"
type Command struct {
	Coauthors        []string
	BaseDir          string
	TemplateFileName string
}

// Dependencies the real-world dependencies of the ExecutorFactory
type Dependencies struct {
	CreateDir         func(path string, perm os.FileMode) error
	WriteFile         func(path string, data []byte, mode os.FileMode) error
	SetCommitTemplate func(path string) error
	SaveStatus        func(state status.State, coauthors ...string) error
}

// ExecutorFactory provisions a Command Processor
func ExecutorFactory(deps Dependencies) func(cmd Command) error {
	return func(cmd Command) error {
		if len(cmd.Coauthors) == 0 {
			return nil
		}

		coauthorsString := PrepareForCommitMessage(cmd.Coauthors)

		err := deps.CreateDir(cmd.BaseDir, os.ModePerm)
		if err != nil {
			return err
		}

		templatePath := fmt.Sprintf("%s/%s", cmd.BaseDir, cmd.TemplateFileName)

		err = deps.WriteFile(templatePath, []byte(coauthorsString), 0644)
		if err != nil {
			return err
		}
		err = deps.SetCommitTemplate(templatePath)
		if err != nil {
			return err
		}
		err = deps.SaveStatus(status.ENABLED, cmd.Coauthors...)
		if err != nil {
			return err
		}
		return nil
	}
}
