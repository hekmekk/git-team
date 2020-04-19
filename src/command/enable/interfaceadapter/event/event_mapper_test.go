package enableeventadapter

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"testing/quick"

	"github.com/hekmekk/git-team/src/command/enable"
	status "github.com/hekmekk/git-team/src/command/status"
	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/events"
	"github.com/hekmekk/git-team/src/core/state"
)

type statusPolicyMock struct {
	apply func() events.Event
}

func (mock statusPolicyMock) Apply() events.Event {
	return mock.apply()
}

func TestFoldErrors(t *testing.T) {
	errPrefix := errors.New("_prefix_")
	errSuffix := errors.New("_suffix_")

	// Note: It is more than twice as slow with this predicate approach... Maybe revert to direct inline calls
	isNotNil := func(err error) bool { return err != nil }
	hasProperPrefix := func(err error) bool { return strings.HasPrefix(err.Error(), errPrefix.Error()) }
	hasProperSuffix := func(err error) bool { return strings.HasSuffix(err.Error(), errSuffix.Error()) }

	errorsGen := func(msg string) bool {
		generatedErr := errors.New(msg)
		errs := []error{errPrefix, generatedErr, errSuffix}

		if foldedErr := foldErrors(errs); isNotNil(foldedErr) && hasProperPrefix(foldedErr) && hasProperSuffix(foldedErr) {
			return true
		}
		return false
	}

	if err := quick.Check(errorsGen, nil); err != nil {
		t.Error(err)
	}
}

func TestMapEventToEffectsSucceeded(t *testing.T) {
	coauthors := []string{"A Coauthor"}
	expectedEffects := []effects.Effect{
		effects.NewPrintMessage(fmt.Sprintf("git-team enabled\n\nco-authors\n─ A Coauthor")),
		effects.NewExitOk(),
	}

	statusPolicy := &statusPolicyMock{
		apply: func() events.Event {
			return status.StateRetrievalSucceeded{State: state.NewStateEnabled(coauthors)}
		},
	}

	effects := MapEventToEffectsFactory(statusPolicy)(enable.Succeeded{})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsAborted(t *testing.T) {
	expectedEffects := []effects.Effect{
		effects.NewPrintMessage("git-team disabled"),
		effects.NewExitOk(),
	}

	statusPolicy := &statusPolicyMock{
		apply: func() events.Event {
			return status.StateRetrievalSucceeded{State: state.NewStateDisabled()}
		},
	}

	effects := MapEventToEffectsFactory(statusPolicy)(enable.Aborted{})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsSucceededButQueryStatusFailed(t *testing.T) {
	err := errors.New("query status failure")

	expectedEffects := []effects.Effect{
		effects.NewPrintErr(err),
		effects.NewExitErr(),
	}

	statusPolicy := &statusPolicyMock{
		apply: func() events.Event {
			return status.StateRetrievalFailed{Reason: err}
		},
	}

	effects := MapEventToEffectsFactory(statusPolicy)(enable.Succeeded{})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsFailed(t *testing.T) {
	err := errors.New("enable failure")

	expectedEffects := []effects.Effect{
		effects.NewPrintErr(err),
		effects.NewExitErr(),
	}

	effects := MapEventToEffectsFactory(nil)(enable.Failed{Reason: []error{err}})

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}

func TestMapEventToEffectsUnknownEvent(t *testing.T) {
	expectedEffects := []effects.Effect{}

	effects := MapEventToEffectsFactory(nil)("UNKNOWN_EVENT")

	if !reflect.DeepEqual(expectedEffects, effects) {
		t.Errorf("expected: %s, got: %s", expectedEffects, effects)
		t.Fail()
	}
}
