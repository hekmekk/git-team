package status

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/hekmekk/git-team/src/config"
)

// TODO: define an interface ?!

func PersistEnabled(coauthors []string) error {
	return persist(state{Status: enabled, Coauthors: coauthors})
}

func PersistDisabled() error {
	return persist(state{Status: disabled})
}

func Print() {
	status := fetch()
	statusStr := toString(status)
	fmt.Println(statusStr)
}

func persist(state state) error {
	deps := persistDependencies{
		loadConfig: config.Load,
		writeFile:  ioutil.WriteFile,
	}
	return persistFactory(deps)(state)
}

func fetch() state {
	deps := loadDependencies{
		loadConfig:     config.Load,
		tomlDecodeFile: toml.DecodeFile,
	}
	return loadStatusFromFileFactory(deps)()
}
