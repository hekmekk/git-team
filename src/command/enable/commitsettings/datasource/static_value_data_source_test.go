package datasource

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/command/enable/commitsettings/entity"
)

func TestReadSucceeds(t *testing.T) {
	home := "/home/some-user"

	deps := dependencies{getEnv: func(variable string) string { return home }}

	expectedSettings := entity.CommitSettings{
		TemplatesBaseDir: fmt.Sprintf("%s/.git-team/commit-templates", home),
		HooksDir:         fmt.Sprintf("%s/.git-team/hooks", home),
	}

	settings := newStaticValueDataSource(deps).Read()

	if !reflect.DeepEqual(expectedSettings, settings) {
		t.Errorf("expected: %s, received %s", expectedSettings, settings)
		t.Fail()
	}
}
