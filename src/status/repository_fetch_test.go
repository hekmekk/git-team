package status

import (
	"errors"
	// "os"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/hekmekk/git-team/src/config"
)

var (
	cfg            = config.Config{TemplateFileName: "TEMPLATE_FILE", BaseDir: "BASE_DIR", StatusFileName: "STATUS_FILE"}
	loadConfig     = func() (config.Config, error) { return cfg, nil }
	tomlDecodeFile = func(string, interface{}) (toml.MetaData, error) { return toml.MetaData{}, nil }
)

func TestFetchSucceeds(t *testing.T) {
	deps := fetchDependencies{
		loadConfig:     loadConfig,
		tomlDecodeFile: tomlDecodeFile,
	}

	_, err := fetchFromFileFactory(deps)()

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestFetchSucceedsWithDefaultIfFileNotPresent(t *testing.T) {
	deps := fetchDependencies{
		loadConfig:     loadConfig,
		tomlDecodeFile: tomlDecodeFile,
	}

	expectedState := state{Status: disabled, Coauthors: []string{}}
	state, err := fetchFromFileFactory(deps)()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedState, state) {
		t.Errorf("expected: %s, received %s", expectedState, state)
		t.Fail()
	}
}

func TestFetchFailsDueToConfigLoadError(t *testing.T) {
	loadConfig := func() (config.Config, error) { return config.Config{}, errors.New("failed to load config") }
	deps := fetchDependencies{
		loadConfig:     loadConfig,
		tomlDecodeFile: tomlDecodeFile,
	}

	_, err := fetchFromFileFactory(deps)()

	if err == nil {
		t.Error("expected failure")
		t.Fail()
	}
}

func TestFetchFailsDueToDecodeError(t *testing.T) {
	tomlDecodeFile := func(string, interface{}) (toml.MetaData, error) {
		return toml.MetaData{}, errors.New("failed to decode")
	}
	deps := fetchDependencies{
		loadConfig:     loadConfig,
		tomlDecodeFile: tomlDecodeFile,
	}

	_, err := fetchFromFileFactory(deps)()

	if err == nil {
		t.Error("expected failure")
		t.Fail()
	}
}
