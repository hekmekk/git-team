package status

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/hekmekk/git-team/src/core/config"
)

// Fetch read the current state from file
func Fetch() (State, error) {
	deps := fetchDependencies{
		loadConfig:     config.Load,
		tomlDecodeFile: toml.DecodeFile,
		statFile:       os.Stat,
		isFileNotExist: os.IsNotExist,
	}
	return fetchFromFileFactory(deps)()
}

// PersistEnabled persist the current state to file as enabled
func PersistEnabled(coauthors []string) error {
	return persist(State{Status: enabled, Coauthors: coauthors})
}

// PersistDisabled persist the current state to file as disabled
func PersistDisabled() error {
	return persist(State{Status: disabled})
}

func persist(state State) error {
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
