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

func fetchFromFileFactory(deps fetchDependencies) func() (state, error) {
	return func() (state, error) {
		cfg, err := deps.loadConfig()
		if err != nil {
			return state{}, err
		}

		// if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
		// return state{Status: disabled, Coauthors: []string{}}, nil
		// }

		var decodedState state
		if _, err := deps.tomlDecodeFile(fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.StatusFileName), &decodedState); err != nil {
			return state{}, err
		}

		return decodedState, nil
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
