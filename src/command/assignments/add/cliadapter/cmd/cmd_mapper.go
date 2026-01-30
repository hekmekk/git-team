package addcmdadapter

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/hekmekk/git-team/v2/src/command/assignments/add"
	addeventadapter "github.com/hekmekk/git-team/v2/src/command/assignments/add/cliadapter/event"
	"github.com/hekmekk/git-team/v2/src/core/validation"
	commandadapter "github.com/hekmekk/git-team/v2/src/shared/cli/commandadapter"
	"github.com/hekmekk/git-team/v2/src/shared/cli/effects"
	gitconfiglegacy "github.com/hekmekk/git-team/v2/src/shared/gitconfig/impl/legacy"
	gitconfigscope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
)

// Command the add command
func Command() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "Add a new or override an existing alias to co-author assignment",
		ArgsUsage: "<alias> <coauthor>",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "force-override", Value: false, Aliases: []string{"f"}, Usage: "Override an existing assignment"},
			&cli.BoolFlag{Name: "keep-existing", Value: false, Aliases: []string{"k"}, Usage: "Keep existing assignment"},
		},
		Action: func(c *cli.Context) error {
			forceOverride := c.Bool("force-override")
			keepExisting := c.Bool("keep-existing")

			if c.NArg() == 0 {
				return handleInputFromStdin(forceOverride, keepExisting).Run()
			}

			if c.NArg() != 2 {
				return effects.NewExitErrMsg(errors.New("exactly 2 arguments expected")).Run()
			}

			args := c.Args()
			alias := args.First()
			coauthor := args.Get(1)
			return commandadapter.Run(policy(&alias, &coauthor, &forceOverride, &keepExisting), addeventadapter.MapEventToEffect)
		},
	}
}

func handleInputFromStdin(forceOverride bool, keepExisting bool) effects.Effect {
	lines, err := readLinesFromStdin()

	if err != nil {
		return effects.NewExitErrMsg(err)
	}

	errOccured := false
	for _, line := range lines {
		argsFromStdin := strings.SplitN(strings.TrimSpace(line), " ", 2)

		var effect effects.Effect
		if len(argsFromStdin) != 2 {
			effect = effects.NewExitErrMsg(errors.New("exactly 2 arguments expected"))
		} else {

			alias := argsFromStdin[0]
			coauthor := argsFromStdin[1]
			effect = commandadapter.ApplyPolicy(policy(&alias, &coauthor, &forceOverride, &keepExisting), addeventadapter.MapEventToEffect)
		}
		err := effect.Run()
		if err != nil {
			errOccured = true
		}
	}

	if errOccured {
		return effects.NewExitErr()
	}

	return effects.NewExitOk()
}

func readLinesFromStdin() ([]string, error) {
	s := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for s.Scan() {
		line := s.Text()
		lines = append(lines, line)
	}

	if s.Err() != nil {
		return []string{}, s.Err()
	}

	return lines, nil
}

func policy(alias *string, coauthor *string, forceOverride *bool, keepExisting *bool) add.Policy {
	return add.Policy{
		Req: add.AssignmentRequest{
			Alias:         alias,
			Coauthor:      coauthor,
			ForceOverride: forceOverride,
			KeepExisting:  keepExisting,
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
