package staterepository

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/hekmekk/git-team/src/core/config"
	"github.com/hekmekk/git-team/src/core/state"
)

// Query read the current state from file
func Query() (state.State, error) {
	deps := queryDependencies{
		loadConfig:     config.Load,
		tomlDecodeFile: toml.DecodeFile,
		statFile:       os.Stat,
		isFileNotExist: os.IsNotExist,
	}
	return queryFileFactory(deps)()
}

// PersistEnabled persist the current state to file as enabled
func PersistEnabled(coauthors []string) error {
	return persist(state.NewStateEnabled(coauthors))
}

// PersistDisabled persist the current state to file as disabled
func PersistDisabled() error {
	return persist(state.NewStateDisabled())
}

func persist(state state.State) error {
	deps := persistDependencies{
		loadConfig: config.Load,
		writeFile:  ioutil.WriteFile,
		tomlEncode: func(state interface{}) ([]byte, error) {
			buf := new(bytes.Buffer)
			err := toml.NewEncoder(buf).Encode(state)
			return buf.Bytes(), err
		},
	}
	return persistToFileFactory(deps)(state)
}

type queryDependencies struct {
	loadConfig     func() config.Config
	tomlDecodeFile func(string, interface{}) (toml.MetaData, error)
	statFile       func(string) (os.FileInfo, error)
	isFileNotExist func(error) bool
}

func queryFileFactory(deps queryDependencies) func() (state.State, error) {
	return func() (state.State, error) {
		cfg := deps.loadConfig()

		stateFilePath := fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.StatusFileName)

		if _, err := deps.statFile(stateFilePath); deps.isFileNotExist(err) {
			return state.NewStateDisabled(), nil
		}

		var decodedState state.State
		if _, err := deps.tomlDecodeFile(stateFilePath, &decodedState); err != nil {
			return state.State{}, err
		}

		return decodedState, nil
	}
}

type persistDependencies struct {
	loadConfig func() config.Config
	writeFile  func(string, []byte, os.FileMode) error
	tomlEncode func(interface{}) ([]byte, error)
}

func persistToFileFactory(deps persistDependencies) func(state state.State) error {
	return func(state state.State) error {
		cfg := deps.loadConfig()

		bytes, err := deps.tomlEncode(state)
		if err != nil {
			return err
		}

		return deps.writeFile(fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.StatusFileName), bytes, 0644)
	}
}
