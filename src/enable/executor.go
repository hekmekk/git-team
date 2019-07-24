package enable

import (
	"fmt"
	"os"

	"github.com/hekmekk/git-team/src/config"
	utils "github.com/hekmekk/git-team/src/enable/utils"
)

// Command add a <Coauthor> under "team.alias.<Alias>"
type Command struct {
	Coauthors []string
}

// TODO: make interface?!
// Dependencies the real-world dependencies of the ExecutorFactory
type Dependencies struct {
	LoadConfig        func() (config.Config, error)
	CreateDir         func(path string, perm os.FileMode) error
	WriteFile         func(path string, data []byte, mode os.FileMode) error
	SetCommitTemplate func(path string) error
	PersistEnabled    func(coauthors []string) error
}

// ExecutorFactory provisions a Command Processor
func ExecutorFactory(deps Dependencies) func(cmd Command) error {
	return func(cmd Command) error {
		if len(cmd.Coauthors) == 0 {
			return nil
		}

		cfg, err := deps.LoadConfig()
		if err != nil {
			return err
		}

		if err := deps.CreateDir(cfg.BaseDir, os.ModePerm); err != nil {
			return err
		}

		templatePath := fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.TemplateFileName)

		if err := deps.WriteFile(templatePath, []byte(utils.PrepareForCommitMessage(cmd.Coauthors)), 0644); err != nil {
			return err
		}
		if err := deps.SetCommitTemplate(templatePath); err != nil {
			return err
		}
		if err := deps.PersistEnabled(cmd.Coauthors); err != nil {
			return err
		}
		return nil
	}
}
