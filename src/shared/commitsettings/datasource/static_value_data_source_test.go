package datasource

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/shared/commitsettings/entity"
)

func TestReadSucceeds(t *testing.T) {
	home := "/home/some-user"

	deps := dependencies{getEnv: func(variable string) string { return home }}

	expectedSettings := entity.CommitSettings{
		GitTeamCommitTemplatePath: fmt.Sprintf("%s/.config/git-team/COMMIT_TEMPLATE", home),
		GitTeamHooksPath:          "/usr/local/etc/git-team/hooks",
	}

	settings := newStaticValueDataSource(deps).Read()

	if !reflect.DeepEqual(expectedSettings, settings) {
		t.Errorf("expected: %s, received %s", expectedSettings, settings)
		t.Fail()
	}
}
