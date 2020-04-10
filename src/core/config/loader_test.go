package config

import (
	"fmt"
	"reflect"
	"testing"
)

func TestLoadSucceeds(t *testing.T) {
	home := "/home/some-user"

	deps := dependencies{getEnv: func(variable string) string { return home }}

	expectedCfg := Config{
		Ro: ReadOnlyProperties{
			GitTeamCommitTemplatePath: fmt.Sprintf("%s/.config/git-team/COMMIT_TEMPLATE", home),
			GitTeamHooksPath:          "/usr/local/etc/git-team/hooks",
		},
		Rw: ReadWriteProperties{
			ActivationScope: Global,
		},
	}

	cfg := executorFactory(deps)()

	if !reflect.DeepEqual(expectedCfg, cfg) {
		t.Errorf("expected: %s, received %s", expectedCfg, cfg)
		t.Fail()
	}
}
