package datasource

import (
	"reflect"
	"testing"

	activationscope "github.com/hekmekk/git-team/src/command/config/entity/activationscope"
	config "github.com/hekmekk/git-team/src/command/config/entity/config"
)

func TestLoadSucceeds(t *testing.T) {
	expectedCfg := config.Config{
		ActivationScope: activationscope.Global,
	}

	cfg, _ := newGitconfigDataSource().Read()

	if !reflect.DeepEqual(expectedCfg, cfg) {
		t.Errorf("expected: %s, received %s", expectedCfg, cfg)
		t.Fail()
	}
}
