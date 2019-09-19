package enablecmdadapter

import (
	"io/ioutil"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/hekmekk/git-team/src/command/adapter"
	"github.com/hekmekk/git-team/src/command/enable"
	"github.com/hekmekk/git-team/src/command/enable/interfaceadapter/event"
	"github.com/hekmekk/git-team/src/core/config"
	"github.com/hekmekk/git-team/src/core/gitconfig"
	"github.com/hekmekk/git-team/src/core/state_repository"
	"github.com/hekmekk/git-team/src/core/validation"
)

// Command the enable command
func Command(root commandadapter.CommandRoot) *kingpin.CmdClause {
	enable := root.Command("enable", "Enables injection of the provided co-authors whenever `git-commit` is used").Default()
	coauthors := enable.Arg("co-authors", "The co-authors for the next commit(s). A co-author must either be an alias or of the shape \"Name <email>\"").Strings()

	enable.Action(commandadapter.Run(policy(coauthors), enableeventadapter.MapEventToEffectsFactory(staterepository.Query)))

	return enable
}

func policy(coauthors *[]string) enable.Policy {
	return enable.Policy{
		Req: enable.Request{
			AliasesAndCoauthors: coauthors,
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
	}
}
