package datasource

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/shared/internalconfig/entity"
)

func TestReadSucceeds(t *testing.T) {
	home := "/home/some-user"

	deps := dependencies{getEnv: func(variable string) string { return home }}

	expectedCfg := entity.InternalConfig{
		GitTeamCommitTemplatePath: fmt.Sprintf("%s/.config/git-team/COMMIT_TEMPLATE", home),
		GitTeamHooksPath:          "/usr/local/etc/git-team/hooks",
	}

	cfg := newStaticValueDataSource(deps).Read()

	if !reflect.DeepEqual(expectedCfg, cfg) {
		t.Errorf("expected: %s, received %s", expectedCfg, cfg)
		t.Fail()
	}
}
