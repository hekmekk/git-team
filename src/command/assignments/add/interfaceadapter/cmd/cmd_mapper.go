package addcmdadapter

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
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
		ArgsUsage: "<alias> <coauthor>",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "force-override", Value: false, Aliases: []string{"f"}, Usage: "Override an existing assignment"},
		},
		Action: func(c *cli.Context) error {
			forceOverride := c.Bool("force-override")

			if c.NArg() == 0 {
				return commandadapter.RunEffect(handleInputFromStdin(forceOverride))
			}

			if c.NArg() != 2 {
				return commandadapter.RunEffect(effects.NewExitErr(errors.New("exactly 2 arguments expected")))
			}

			args := c.Args()
			alias := args.First()
			coauthor := args.Get(1)
			return commandadapter.Run(policy(&alias, &coauthor, &forceOverride), addeventadapter.MapEventToEffect)
		},
	}
}

func handleInputFromStdin(forceOverride bool) effects.Effect {
	lines, err := readLinesFromStdin()

	if err != nil {
		return effects.NewExitErr(err)
	}

	errOccured := false
	for _, line := range lines {
		argsFromStdin := strings.SplitN(strings.TrimSpace(line), " ", 2)

		var effect effects.Effect
		if len(argsFromStdin) != 2 {
			effect = effects.NewExitErr(errors.New("exactly 2 arguments expected"))
		} else {

			alias := argsFromStdin[0]
			coauthor := argsFromStdin[1]
			effect = commandadapter.ApplyPolicy(policy(&alias, &coauthor, &forceOverride), addeventadapter.MapEventToEffect)
		}
		msg, localErrOccurred := derivePrintMsgAndDetermineIfErrorOccured(effect)
		if localErrOccurred {
			errOccured = true
		}
		commandadapter.RunEffect(effects.NewExitOkMsg(msg))
	}

	if errOccured {
		return effects.NewExitErr(errors.New(""))
	}

	return effects.NewExitOk()
}

func derivePrintMsgAndDetermineIfErrorOccured(effect effects.Effect) (string, bool) {
	switch e := effect.(type) {
	case effects.ExitOk:
		return e.Message(), false
	case effects.ExitWarn:
		return color.YellowString(e.Message()), false
	case effects.ExitErr:
		return color.RedString(e.Message()), true
	default:
		return color.RedString("error: undefined behavior encountered"), true
	}
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
