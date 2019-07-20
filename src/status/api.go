package status

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/hekmekk/git-team/src/config"
)

// Fetch read the current state from file
func Fetch() (state, error) {
	deps := fetchDependencies{
		loadConfig:     config.Load,
		tomlDecodeFile: toml.DecodeFile,
	}
	return fetchFromFileFactory(deps)()
}

// PersistEnabled persist the current state to file as enabled
func PersistEnabled(coauthors []string) error {
	return persist(state{Status: enabled, Coauthors: coauthors})
}

// PersistDisabled persist the current state to file as disabled
func PersistDisabled() error {
	return persist(state{Status: disabled})
}

func persist(state state) error {
	deps := persistDependencies{
		loadConfig: config.Load,
		writeFile:  ioutil.WriteFile,
	}
	return persistToFileFactory(deps)(state)
}
