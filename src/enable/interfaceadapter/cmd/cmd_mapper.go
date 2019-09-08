package enablecmdadapter

import (
	"io/ioutil"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/hekmekk/git-team/src/core/config"
	"github.com/hekmekk/git-team/src/core/gitconfig"
	"github.com/hekmekk/git-team/src/core/state_repository"
	"github.com/hekmekk/git-team/src/core/validation"
	"github.com/hekmekk/git-team/src/enable"
)

// Definition definition of the add command
type Definition struct {
	CommandName string
	Policy      enable.Policy
}

// NewDefinition the constructor for Definition
func NewDefinition(app *kingpin.Application) Definition {
	command := app.Command("enable", "Enables injection of the provided co-authors whenever `git-commit` is used").Default()

	return Definition{
		CommandName: command.FullCommand(),
		Policy: enable.Policy{
			Req: enable.Request{
				AliasesAndCoauthors: command.Arg("coauthors", "The co-authors for the next commit(s). A co-author must either be an alias or of the shape \"Name <email>\"").Strings(),
			},
			Deps: enable.Dependencies{
				SanityCheckCoauthors:          validation.SanityCheckCoauthors,
				CreateTemplateDir:             os.MkdirAll,
				WriteTemplateFile:             ioutil.WriteFile,
				GitSetCommitTemplate:          gitconfig.SetCommitTemplate,
				GitSetHooksPath:               func(path string) error { return gitconfig.ReplaceAll("core.hooksPath", path) },
				GitResolveAliases:             gitconfig.ResolveAliases,
				StateRepositoryPersistEnabled: staterepository.PersistEnabled,
				LoadConfig:                    config.Load,
			},
		},
	}
}
