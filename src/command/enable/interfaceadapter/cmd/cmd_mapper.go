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
	"github.com/hekmekk/git-team/src/core/validation"
	activation "github.com/hekmekk/git-team/src/shared/activation/impl"
	configds "github.com/hekmekk/git-team/src/shared/config/datasource"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/impl"
	state "github.com/hekmekk/git-team/src/shared/state/impl"
)

// Command the enable command
func Command(root commandadapter.CommandRoot) *kingpin.CmdClause {
	enable := root.Command("enable", "Enables injection of the provided co-authors whenever `git-commit` is used")
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
			SanityCheckCoauthors: validation.SanityCheckCoauthors,
			CreateTemplateDir:    os.MkdirAll,
			WriteTemplateFile:    ioutil.WriteFile,
			GitConfigWriter:      gitconfig.NewDataSink(),
			GitResolveAliases:    commandadapter.ResolveAliases,
			CommitSettingsReader: commitsettingsds.NewStaticValueDataSource(),
			ConfigReader:         configds.NewGitconfigDataSource(gitconfig.NewDataSource()),
			StateWriter:          state.NewGitConfigDataSink(gitconfig.NewDataSink()),
			GetEnv:               os.Getenv,
			GetWd:                os.Getwd,
			ActivationValidator:  activation.NewGitConfigDataSource(gitconfig.NewDataSource()),
		},
	}
}
