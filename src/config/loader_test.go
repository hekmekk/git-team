package config

import (
	"errors"
	"reflect"
	"testing"
)

func TestLoadSucceeds(t *testing.T) {
	expectedBaseDir := "BASEDIR"

	deps := dependencies{expandHomedir: func(path string) (string, error) { return expectedBaseDir, nil }}

	expectedCfg := Config{BaseDir: expectedBaseDir, TemplateFileName: "COMMIT_TEMPLATE", StatusFileName: "status.toml"}

	cfg, err := executorFactory(deps)()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedCfg, cfg) {
		t.Errorf("expected: %s, received %s", expectedCfg, cfg)
		t.Fail()
	}
}

func TestLoadFailsBecauseHomeDirExpansionFails(t *testing.T) {
	deps := dependencies{expandHomedir: func(path string) (string, error) { return "", errors.New("failed to expand dir") }}

	expectedCfg := Config{}

	cfg, err := executorFactory(deps)()

	if err == nil {
		t.Error("expected failure")
		t.Fail()
	}

	if !reflect.DeepEqual(expectedCfg, cfg) {
		t.Errorf("expected: %s, received %s", expectedCfg, cfg)
		t.Fail()
	}
}
