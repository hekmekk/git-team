package status

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/hekmekk/git-team/src/config"
)

type fetchDependencies struct {
	loadConfig     func() (config.Config, error)
	tomlDecodeFile func(string, interface{}) (toml.MetaData, error)
	statFile       func(string) (os.FileInfo, error)
	isFileNotExist func(error) bool
}

func fetchFromFileFactory(deps fetchDependencies) func() (state, error) {
	return func() (state, error) {
		cfg, err := deps.loadConfig()
		if err != nil {
			return state{}, err
		}

		stateFilePath := fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.StatusFileName)

		if _, err := deps.statFile(stateFilePath); deps.isFileNotExist(err) {
			return state{Status: disabled, Coauthors: []string{}}, nil
		}

		var decodedState state
		if _, err := deps.tomlDecodeFile(stateFilePath, &decodedState); err != nil {
			return state{}, err
		}

		return decodedState, nil
	}
}

type persistDependencies struct {
	loadConfig func() (config.Config, error)
	writeFile  func(string, []byte, os.FileMode) error
	tomlEncode func(interface{}) ([]byte, error)
}

func persistToFileFactory(deps persistDependencies) func(state state) error {
	return func(state state) error {
		cfg, err := deps.loadConfig()
		if err != nil {
			return err
		}

		bytes, err := deps.tomlEncode(state)
		if err != nil {
			return err
		}

		return deps.writeFile(fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.StatusFileName), bytes, 0644)
	}
}
