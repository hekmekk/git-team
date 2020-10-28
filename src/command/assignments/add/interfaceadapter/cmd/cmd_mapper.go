package addcmdadapter

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	commandadapter "github.com/hekmekk/git-team/src/command/adapter"
	"github.com/hekmekk/git-team/src/command/assignments/add"
	addeventadapter "github.com/hekmekk/git-team/src/command/assignments/add/interfaceadapter/event"
	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/validation"
	gitconfiglegacy "github.com/hekmekk/git-team/src/shared/gitconfig/impl/legacy"
	gitconfigscope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
)

// Command the add command
func Command() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "Add a new or override an existing alias to co-author assignment",
		ArgsUsage: "alias - The alias to assign a co-author to\ncoauthor - The co-author",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "force-override", Value: false, Aliases: []string{"f"}, Usage: "Override an existing assignment"},
		},
		Action: func(c *cli.Context) error {
			args := c.Args()
			if args.Len() != 2 {
				effects.NewPrintErr(errors.New("exactly 2 arguments expected")).Run()
				return nil
			}

			alias := args.First()
			coauthor := args.Get(1)
			forceOverride := c.Bool("force-override")
			return commandadapter.RunUrFave(policy(&alias, &coauthor, &forceOverride), addeventadapter.MapEventToEffects)(c)
		},
	}
}

func policy(alias *string, coauthor *string, forceOverride *bool) add.Policy {
	return add.Policy{
		Req: add.AssignmentRequest{
			Alias:         alias,
			Coauthor:      coauthor,
			ForceOverride: forceOverride,
		},
		Deps: add.Dependencies{
			SanityCheckCoauthor: validation.SanityCheckCoauthor,
			GitResolveAlias:     commandadapter.ResolveAlias,
			GitAddAlias: func(alias string, coauthor string) error {
				return gitconfiglegacy.ReplaceAll(gitconfigscope.Global, fmt.Sprintf("team.alias.%s", alias), coauthor)
			},
			GetAnswerFromUser: func(question string) (string, error) {
				_, err := os.Stdout.WriteString(question)
				if err != nil {
					return "", err
				}
				return bufio.NewReader(os.Stdin).ReadString('\n')
			},
		},
	}
}
