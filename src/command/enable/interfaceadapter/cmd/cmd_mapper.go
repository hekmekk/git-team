package enablecmdadapter

import (
	"io/ioutil"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	commandadapter "github.com/hekmekk/git-team/src/command/adapter"
	"github.com/hekmekk/git-team/src/command/enable"
	commitsettingsds "github.com/hekmekk/git-team/src/command/enable/commitsettings/datasource"
	enableeventadapter "github.com/hekmekk/git-team/src/command/enable/interfaceadapter/event"
	statuscmdmapper "github.com/hekmekk/git-team/src/command/status/interfaceadapter/cmd"
	coregitconfig "github.com/hekmekk/git-team/src/core/gitconfig"
	staterepository "github.com/hekmekk/git-team/src/core/state_repository"
	"github.com/hekmekk/git-team/src/core/validation"
	configds "github.com/hekmekk/git-team/src/shared/config/datasource"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/impl"
)

// Command the enable command
func Command(root commandadapter.CommandRoot) *kingpin.CmdClause {
	enable := root.Command("enable", "Enables injection of the provided co-authors whenever `git-commit` is used").Default()
	coauthors := enable.Arg("co-authors", "The co-authors for the next commit(s). A co-author must either be an alias or of the shape \"Name <email>\"").Strings()

	enable.Action(commandadapter.Run(policy(coauthors), enableeventadapter.MapEventToEffectsFactory(statuscmdmapper.Policy())))

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
			GitSetCommitTemplate:          func(path string) error { return coregitconfig.ReplaceAll("commit.template", path) },
			GitSetHooksPath:               func(path string) error { return coregitconfig.ReplaceAll("core.hooksPath", path) },
			GitResolveAliases:             commandadapter.ResolveAliases,
			StateRepositoryPersistEnabled: staterepository.PersistEnabled,
			CommitSettingsReader:          commitsettingsds.NewStaticValueDataSource(),
			ConfigReader:                  configds.NewGitconfigDataSource(gitconfig.NewDataSource()),
		},
	}
}
