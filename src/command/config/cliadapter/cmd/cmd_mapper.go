package configcmdadapter

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

	"github.com/urfave/cli/v2"

	configeventadapter "github.com/hekmekk/git-team/v2/src/command/config/cliadapter/event"
	configpolicy "github.com/hekmekk/git-team/v2/src/command/config/policy"
	commandadapter "github.com/hekmekk/git-team/v2/src/shared/cli/commandadapter"
	configdatasink "github.com/hekmekk/git-team/v2/src/shared/config/datasink"
	configdatasource "github.com/hekmekk/git-team/v2/src/shared/config/datasource"
	gitconfig "github.com/hekmekk/git-team/v2/src/shared/gitconfig/impl"
)

// Command the config command
func Command() *cli.Command {
	return &cli.Command{
		Name:      "config",
		Usage:     "Display and edit the configuration",
		ArgsUsage: "[<key> <value>]",
		Action: func(c *cli.Context) error {
			args := c.Args()
			key := args.First()
			value := args.Get(1)
			return commandadapter.Run(policy(&key, &value), configeventadapter.MapEventToEffect)
		},
		BashComplete: func(c *cli.Context) {
			options := map[string][]string{
				"activation-scope": []string{"repo-local", "global"},
			}

			args := c.Args()
			argsLen := args.Len()

			if argsLen == 0 {
				for key := range options {
					fmt.Println(key)
				}
			}

			if argsLen == 1 {
				values := options[args.First()]
				for _, value := range values {
					fmt.Println(value)
				}
			}
		},
	}
}

func policy(key *string, value *string) configpolicy.Policy {
	return configpolicy.Policy{
		Req: configpolicy.Request{
			Key:   key,
			Value: value,
		},
		Deps: configpolicy.Dependencies{
			ConfigReader: configdatasource.NewGitconfigDataSource(gitconfig.NewDataSource()),
			ConfigWriter: configdatasink.NewGitconfigDataSink(gitconfig.NewDataSink()),
		},
	}
}
