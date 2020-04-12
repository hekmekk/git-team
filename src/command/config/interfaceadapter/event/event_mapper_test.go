package configeventadapter

import (
	"errors"
	"reflect"
	"testing"

	configevents "github.com/hekmekk/git-team/src/command/config/events"
	config "github.com/hekmekk/git-team/src/core/config"
	"github.com/hekmekk/git-team/src/core/effects"
)

func TestMapEventToEffectsRetrievalSucceeded(t *testing.T) {
	msg := "config\n─ [ro] commit-template-path: /home/some-user/.config/git-team/COMMIT_TEMPLATE\n─ [ro] hooks-path: /usr/local/etc/git-team/hooks\n─ [rw] activation-scope: global"

	cfg := config.Config{
		Ro: config.ReadOnlyProperties{
			GitTeamCommitTemplatePath: "/home/some-user/.config/git-team/COMMIT_TEMPLATE",
			GitTeamHooksPath:          "/usr/local/etc/git-team/hooks",
		},
		Rw: config.ReadWriteProperties{
			ActivationScope: config.Global,
		},
	}

	expectedEffects := []effects.Effect{
		effects.NewPrintMessage(msg),
		effects.NewExitOk(),
	}

	effects := MapEventToEffects(configevents.RetrievalSucceeded{Config: cfg})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsRetrievalFailed(t *testing.T) {
	err := errors.New("failure")

	expectedEffects := []effects.Effect{
		effects.NewPrintErr(err),
		effects.NewExitErr(),
	}

	effects := MapEventToEffects(configevents.RetrievalFailed{Reason: err})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsUnknownEvent(t *testing.T) {
	expectedEffects := []effects.Effect{}

	effects := MapEventToEffects("UNKNOWN_EVENT")

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}
