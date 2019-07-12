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

		if err := deps.CreateDir(cmd.BaseDir, os.ModePerm); err != nil {
			return err
		}

		templatePath := fmt.Sprintf("%s/%s", cmd.BaseDir, cmd.TemplateFileName)

		if err := deps.WriteFile(templatePath, []byte(PrepareForCommitMessage(cmd.Coauthors)), 0644); err != nil {
			return err
		}
		if err := deps.SetCommitTemplate(templatePath); err != nil {
			return err
		}
		if err := deps.SaveStatus(status.ENABLED, cmd.Coauthors...); err != nil {
			return err
		}
		return nil
	}
}
