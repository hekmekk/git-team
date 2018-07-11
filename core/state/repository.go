package state

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
)

func Save(baseDir, stateFile string, status Status, coauthors ...string) error {
	state := State{Status: status, CoAuthors: coauthors}
	buf := new(bytes.Buffer)

	err := toml.NewEncoder(buf).Encode(state)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fmt.Sprintf("%s/%s", baseDir, stateFile), []byte(buf.String()), 0644)
}

func Print(baseDir, stateFile string) {
	var state State
	if _, err := toml.DecodeFile(fmt.Sprintf("%s/%s", baseDir, stateFile), &state); err != nil {
		color.Yellow("Unable to determine current state. Assuming disabled.")
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
