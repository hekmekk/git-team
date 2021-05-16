package enablecmdadapter

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli/v2"

	commandadapter "github.com/hekmekk/git-team/src/command/adapter"
	"github.com/hekmekk/git-team/src/command/enable"
	commitsettingsds "github.com/hekmekk/git-team/src/command/enable/commitsettings/datasource"
	enableeventadapter "github.com/hekmekk/git-team/src/command/enable/interfaceadapter/event"
	statuscmdmapper "github.com/hekmekk/git-team/src/command/status/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/core/validation"
	activation "github.com/hekmekk/git-team/src/shared/activation/impl"
	aliascompletion "github.com/hekmekk/git-team/src/shared/completion"
	configds "github.com/hekmekk/git-team/src/shared/config/datasource"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/impl"
	state "github.com/hekmekk/git-team/src/shared/state/impl"
)

// Command the enable command
func Command() *cli.Command {
	return &cli.Command{
		Name:      "enable",
		Usage:     "Enables injection of the provided co-authors whenever `git-commit` is used",
		ArgsUsage: "<co-authors> (A co-author must either be an alias or of the shape \"Name <email>\")",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "all", Value: false, Aliases: []string{"A"}, Usage: "Use all known co-authors"},
		},
		Action: func(c *cli.Context) error {
			coauthors := c.Args().Slice()
			useAll := c.Bool("all")
			return commandadapter.RunUrFave(policy(&coauthors, &useAll), enableeventadapter.MapEventToEffectsFactory(statuscmdmapper.Policy()))(c)
		},
		BashComplete: func(c *cli.Context) {
			remainingAliases := aliascompletion.NewAliasShellCompletion(gitconfig.NewDataSource()).Complete(c.Args().Slice())
			for _, alias := range remainingAliases {
				fmt.Println(alias)
			}
		},
	}
}

func policy(coauthors *[]string, useAll *bool) enable.Policy {
	return enable.Policy{
		Req: enable.Request{
			AliasesAndCoauthors: coauthors,
			UseAll:              useAll,
		},
		Deps: enable.Dependencies{
			SanityCheckCoauthors: validation.SanityCheckCoauthors,
			CreateTemplateDir:    os.MkdirAll,
			WriteTemplateFile:    ioutil.WriteFile,
			CreateHooksDir:       os.MkdirAll,
			WriteHookFile:        ioutil.WriteFile,
			Lstat:                os.Lstat,
			Remove:               os.Remove,
			Symlink:              os.Symlink,
			GitConfigWriter:      gitconfig.NewDataSink(),
			GitConfigReader:      gitconfig.NewDataSource(),
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
