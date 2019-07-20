package status

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/hekmekk/git-team/src/config"
)

type fetchDependencies struct {
	loadConfig     func() (config.Config, error)
	tomlDecodeFile func(string, interface{}) (toml.MetaData, error)
}

func fetchFromFileFactory(deps fetchDependencies) func() state {
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

func persistToFileFactory(deps persistDependencies) func(state state) error {
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
