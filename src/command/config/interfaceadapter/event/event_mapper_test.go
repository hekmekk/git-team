package configeventadapter

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	configevents "github.com/hekmekk/git-team/src/command/config/events"
	"github.com/hekmekk/git-team/src/core/effects"
	activationscope "github.com/hekmekk/git-team/src/shared/activation/scope"
	config "github.com/hekmekk/git-team/src/shared/config/entity/config"
)

func TestMapEventToEffectRetrievalSucceeded(t *testing.T) {
	msg := "config\n─ activation-scope: global"

	cfg := config.Config{
		ActivationScope: activationscope.Global,
	}

	expectedEffect := effects.NewExitOkMsg(msg)

	effect := MapEventToEffect(configevents.RetrievalSucceeded{Config: cfg})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectRetrievalFailed(t *testing.T) {
	err := errors.New("failed to retrieve config")

	expectedEffect := effects.NewExitErr(err)

	effect := MapEventToEffect(configevents.RetrievalFailed{Reason: err})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectSettingModificationSucceeded(t *testing.T) {
	key := "activation-scope"
	value := "repo-local"

	msg := fmt.Sprintf("Configuration updated: '%s' → '%s'", key, value)

	expectedEffect := effects.NewExitOkMsg(msg)

	effect := MapEventToEffect(configevents.SettingModificationSucceeded{Key: key, Value: value})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectSettingModificationFailed(t *testing.T) {
	err := errors.New("failed to modify setting")

	expectedEffect := effects.NewExitErr(err)

	effect := MapEventToEffect(configevents.SettingModificationFailed{Reason: err})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectReadingSingleSettingNotYetImplemented(t *testing.T) {
	err := errors.New("Reading a single setting has not yet been implemented")

	expectedEffect := effects.NewExitErr(err)

	effect := MapEventToEffect(configevents.ReadingSingleSettingNotYetImplemented{})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectUnknownEvent(t *testing.T) {
	expectedEffect := effects.NewExitOk()

	effect := MapEventToEffect("UNKNOWN_EVENT")

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}
