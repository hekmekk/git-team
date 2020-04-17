package configeventadapter

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	configevents "github.com/hekmekk/git-team/src/command/config/events"
	"github.com/hekmekk/git-team/src/core/effects"
	activationscope "github.com/hekmekk/git-team/src/shared/config/entity/activationscope"
	config "github.com/hekmekk/git-team/src/shared/config/entity/config"
)

func TestMapEventToEffectsRetrievalSucceeded(t *testing.T) {
	msg := "config\n─ activation-scope: global"

	cfg := config.Config{
		ActivationScope: activationscope.Global,
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
	err := errors.New("failed to retrieve config")

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

func TestMapEventToEffectsSettingModificationSucceeded(t *testing.T) {
	key := "activation-scope"
	value := "repo-local"

	msg := fmt.Sprintf("Configuration updated: '%s' → '%s'", key, value)

	expectedEffects := []effects.Effect{
		effects.NewPrintMessage(msg),
		effects.NewExitOk(),
	}

	effects := MapEventToEffects(configevents.SettingModificationSucceeded{Key: key, Value: value})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsSettingModificationFailed(t *testing.T) {
	err := errors.New("failed to modify setting")

	expectedEffects := []effects.Effect{
		effects.NewPrintErr(err),
		effects.NewExitErr(),
	}

	effects := MapEventToEffects(configevents.SettingModificationFailed{Reason: err})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsReadingSingleSettingNotYetImplemented(t *testing.T) {
	err := errors.New("Reading a single setting has not yet been implemented")

	expectedEffects := []effects.Effect{
		effects.NewPrintErr(err),
		effects.NewExitErr(),
	}

	effects := MapEventToEffects(configevents.ReadingSingleSettingNotYetImplemented{})

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
