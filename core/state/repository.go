package state

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/hekmekk/git-team/core/config"
	"io/ioutil"
	"os"
)

func Save(status Status, coauthors ...string) error {
	cfg, _ := config.Load()

	state := State{Status: status, CoAuthors: coauthors}
	buf := new(bytes.Buffer)

	err := toml.NewEncoder(buf).Encode(state)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.StatusFileName), []byte(buf.String()), 0644)
}

func Print() {
	cfg, _ := config.Load()

	var state State
	if _, err := toml.DecodeFile(fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.StatusFileName), &state); err != nil {
		color.Red("Team mode disabled.")
		os.Exit(0)
	}

	switch state.Status {
	case ENABLED:
		color.Green("Team mode enabled.")
		coauthors := state.CoAuthors
		if len(coauthors) > 0 {
			blackBold := color.New(color.FgBlack).Add(color.Bold)
			fmt.Println()
			blackBold.Println("Co-authors:")
			blackBold.Println("-----------")
			for _, coauthor := range coauthors {
				color.Magenta(coauthor)
			}
		}
	default:
		color.Red("Team mode disabled.")
	}
}
