package enablecmdadapter

/*
Copyright (C) 2026 Rea Sand

This file is part of git-team.

This program is free software: you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation, either version 3
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/hekmekk/git-team/v2/src/command/enable"
	enableeventadapter "github.com/hekmekk/git-team/v2/src/command/enable/cliadapter/event"
	commitsettingsds "github.com/hekmekk/git-team/v2/src/command/enable/commitsettings/datasource"
	statuscmdmapper "github.com/hekmekk/git-team/v2/src/command/status/cliadapter/cmd"
	"github.com/hekmekk/git-team/v2/src/core/validation"
	activation "github.com/hekmekk/git-team/v2/src/shared/activation/impl"
	commandadapter "github.com/hekmekk/git-team/v2/src/shared/cli/commandadapter"
	aliascompletion "github.com/hekmekk/git-team/v2/src/shared/completion"
	configds "github.com/hekmekk/git-team/v2/src/shared/config/datasource"
	gitconfig "github.com/hekmekk/git-team/v2/src/shared/gitconfig/impl"
	state "github.com/hekmekk/git-team/v2/src/shared/state/impl"
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
			return commandadapter.Run(policy(&coauthors, &useAll), enableeventadapter.MapEventToEffectFactory(statuscmdmapper.Policy(false)))
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
