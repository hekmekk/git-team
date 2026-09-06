package configeventadapter

/*
Copyright (C) 2026 Rea Sand

This file is part of git-team.

This program is free software: you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation, either version 3
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	configevents "github.com/hekmekk/git-team/v2/src/command/config/events"
	activationscope "github.com/hekmekk/git-team/v2/src/shared/activation/scope"
	"github.com/hekmekk/git-team/v2/src/shared/cli/effects"
	config "github.com/hekmekk/git-team/v2/src/shared/config/entity/config"
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

	expectedEffect := effects.NewExitErrMsg(err)

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

	expectedEffect := effects.NewExitErrMsg(err)

	effect := MapEventToEffect(configevents.SettingModificationFailed{Reason: err})

	if !reflect.DeepEqual(expectedEffect, effect) {
		t.Errorf("expected: %s, got: %s", expectedEffect, effect)
		t.Fail()
	}
}

func TestMapEventToEffectReadingSingleSettingNotYetImplemented(t *testing.T) {
	err := errors.New("Reading a single setting has not yet been implemented")

	expectedEffect := effects.NewExitErrMsg(err)

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
