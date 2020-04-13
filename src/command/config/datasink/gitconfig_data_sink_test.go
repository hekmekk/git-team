package datasink

import (
	"testing"

	activationscope "github.com/hekmekk/git-team/src/command/config/entity/activationscope"
)

func TestSetActivationScopeSucceeds(t *testing.T) {
	err := newGitconfigDataSink().SetActivationScope(activationscope.RepoLocal)

	if err != nil {
		t.Errorf("expected no error, received %s", err)
		t.Fail()
	}
}
