package staterepository

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/hekmekk/git-team/src/core/config"
	"github.com/hekmekk/git-team/src/core/state"
)

// TODO: add assertions on arguments to mocked functions
var (
	cfg            = config.Config{TemplateFileName: "TEMPLATE_FILE", BaseDir: "BASE_DIR", StatusFileName: "STATUS_FILE"}
	loadConfig     = func() (config.Config, error) { return cfg, nil }
	tomlDecodeFile = func(string, interface{}) (toml.MetaData, error) { return toml.MetaData{}, nil }
	fileInfo       os.FileInfo
	writeFile      = func(string, []byte, os.FileMode) error { return nil }
	statFile       = func(string) (os.FileInfo, error) { return fileInfo, nil }
	isFileNotExist = func(error) bool { return false }
	tomlEncode     = func(interface{}) ([]byte, error) { return nil, nil }
)

func TestQuerySucceeds(t *testing.T) {
	deps := queryDependencies{
		loadConfig:     loadConfig,
		tomlDecodeFile: tomlDecodeFile,
		statFile:       statFile,
		isFileNotExist: isFileNotExist,
	}

	_, err := queryFileFactory(deps)()

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestQuerySucceedsWithDefaultIfFileNotPresent(t *testing.T) {
	statFile := func(string) (os.FileInfo, error) { return fileInfo, errors.New("failed to stat path") }
	isFileNotExist := func(error) bool { return true }
	deps := queryDependencies{
		loadConfig:     loadConfig,
		tomlDecodeFile: tomlDecodeFile,
		statFile:       statFile,
		isFileNotExist: isFileNotExist,
	}

	expectedState := state.NewStateDisabled()
	state, err := queryFileFactory(deps)()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedState, state) {
		t.Errorf("expected: %s, received %s", expectedState, state)
		t.Fail()
	}
}

func TestQueryFailsDueToConfigLoadError(t *testing.T) {
	loadConfig := func() (config.Config, error) { return config.Config{}, errors.New("failed to load config") }
	deps := queryDependencies{
		loadConfig:     loadConfig,
		tomlDecodeFile: tomlDecodeFile,
		statFile:       statFile,
		isFileNotExist: isFileNotExist,
	}

	_, err := queryFileFactory(deps)()

	if err == nil {
		t.Error("expected failure")
		t.Fail()
	}
}

func TestQueryFailsDueToDecodeError(t *testing.T) {
	tomlDecodeFile := func(string, interface{}) (toml.MetaData, error) {
		return toml.MetaData{}, errors.New("failed to decode")
	}
	deps := queryDependencies{
		loadConfig:     loadConfig,
		tomlDecodeFile: tomlDecodeFile,
		statFile:       statFile,
		isFileNotExist: isFileNotExist,
	}

	_, err := queryFileFactory(deps)()

	if err == nil {
		t.Error("expected failure")
		t.Fail()
	}
}

func TestPersistSucceeds(t *testing.T) {
	deps := persistDependencies{
		loadConfig: loadConfig,
		writeFile:  writeFile,
		tomlEncode: tomlEncode,
	}

	// TODO: assert that this is passed to tomlEncode
	state := state.NewStateEnabled([]string{"CO-AUTHOR"})

	err := persistToFileFactory(deps)(state)

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestPersistFailsDueToConfigLoadError(t *testing.T) {
	loadConfig := func() (config.Config, error) { return config.Config{}, errors.New("failed to load config") }
	deps := persistDependencies{
		loadConfig: loadConfig,
		writeFile:  writeFile,
		tomlEncode: tomlEncode,
	}

	err := persistToFileFactory(deps)(state.NewStateDisabled())

	if err == nil {
		t.Error("expected failure")
		t.Fail()
	}
}

func TestPersistFailsDueToTomlEncodeError(t *testing.T) {
	tomlEncode := func(interface{}) ([]byte, error) {
		return nil, errors.New("failed to encode struct with toml encoder")
	}
	deps := persistDependencies{
		loadConfig: loadConfig,
		writeFile:  writeFile,
		tomlEncode: tomlEncode,
	}

	err := persistToFileFactory(deps)(state.NewStateDisabled())

	if err == nil {
		t.Error("expected failure")
		t.Fail()
	}
}

func TestPersistFailsDueToWriteFileError(t *testing.T) {
	writeFile := func(string, []byte, os.FileMode) error { return errors.New("failed to write file") }
	deps := persistDependencies{
		loadConfig: loadConfig,
		writeFile:  writeFile,
		tomlEncode: tomlEncode,
	}

	err := persistToFileFactory(deps)(state.NewStateDisabled())

	if err == nil {
		t.Error("expected failure")
		t.Fail()
	}
}
