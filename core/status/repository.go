package status

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/hekmekk/git-team/src/config"
)

type teamStatus string

const (
	enabled  teamStatus = "enabled"
	disabled teamStatus = "disabled"
)

type state struct {
	Status    teamStatus `toml:"status"`
	Coauthors []string   `toml:"co-authors"`
}

type loadDependencies struct {
	loadConfig     func() (config.Config, error)
	tomlDecodeFile func(string, interface{}) (toml.MetaData, error)
}

func loadStatusFromFileFactory(deps loadDependencies) func() state {
	return func() state {
		cfg, _ := deps.loadConfig()

		var status state
		if _, err := deps.tomlDecodeFile(fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.StatusFileName), &status); err != nil {
			return state{Status: disabled, Coauthors: []string{}}
		}

		return status
	}
}

type persistDependencies struct {
	loadConfig func() (config.Config, error)
	writeFile  func(string, []byte, os.FileMode) error
}

func persistFactory(deps persistDependencies) func(state state) error {
	return func(state state) error {
		cfg, _ := config.Load()

		buf := new(bytes.Buffer)

		err := toml.NewEncoder(buf).Encode(state)
		if err != nil {
			return err
		}

		return ioutil.WriteFile(fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.StatusFileName), []byte(buf.String()), 0644)
	}
}

// TODO: should this really be here?
func toString(state state) string {
	var buffer bytes.Buffer
	switch state.Status {
	case enabled:
		buffer.WriteString(color.GreenString("git-team %s.", enabled))
		coauthors := state.Coauthors
		if len(coauthors) > 0 {
			buffer.WriteString("\n\n")
			blackBold := color.New(color.FgBlack).Add(color.Bold)
			buffer.WriteString(blackBold.Sprintln("Co-authors:"))
			buffer.WriteString(blackBold.Sprint("-----------"))
			for _, coauthor := range coauthors {
				buffer.WriteString(color.MagentaString("\n%s", coauthor))
			}
		}
	default:
		buffer.WriteString(color.RedString("git-team %s", disabled))
	}

	return buffer.String()
}
