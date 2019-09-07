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
		BaseDir:          fmt.Sprintf("%s/.config/git-team", home),
		GitHooksPath:     "/usr/local/etc/git-team/hooks",
		TemplateFileName: "COMMIT_TEMPLATE",
		StatusFileName:   "status.toml",
	}

	cfg := executorFactory(deps)()

	if !reflect.DeepEqual(expectedCfg, cfg) {
		t.Errorf("expected: %s, received %s", expectedCfg, cfg)
		t.Fail()
	}
}
