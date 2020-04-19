package status

import (
	"errors"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/core/state"
	activationscope "github.com/hekmekk/git-team/src/shared/config/entity/activationscope"
	config "github.com/hekmekk/git-team/src/shared/config/entity/config"
)

type configReaderMock struct {
	read func() (config.Config, error)
}

func (mock configReaderMock) Read() (config.Config, error) {
	return mock.read()
}

func TestStatusShouldBeRetrieved(t *testing.T) {
	t.Parallel()

	currState := state.NewStateDisabled()

	deps := Dependencies{
		StateRepositoryQuery: func() (state.State, error) { return currState, nil },
	}

	cases := []struct {
		activationScope activationscope.ActivationScope
		gitconfigScope  string // TODO: introduce a type for this
	}{
		{activationscope.Global, "global"},
		{activationscope.RepoLocal, "local"},
	}

	for _, caseLoopVar := range cases {
		activationScope := caseLoopVar.activationScope
		// gitconfigScope := caseLoopVar.gitconfigScope

		t.Run(activationScope.String(), func(t *testing.T) {
			t.Parallel()

			deps.ConfigReader = &configReaderMock{
				read: func() (config.Config, error) {
					return config.Config{ActivationScope: activationScope}, nil
				},
			}

			expectedEvent := StateRetrievalSucceeded{State: currState}

			event := Policy{deps}.Apply()

			if !reflect.DeepEqual(expectedEvent, event) {
				t.Errorf("expected: %s, got: %s", expectedEvent, event)
				t.Fail()
			}
		})
	}
}

func TestStatusShouldNotBeRetrievedDueToConfigReaderError(t *testing.T) {
	err := errors.New("failed to read config")

	deps := Dependencies{
		ConfigReader: &configReaderMock{
			read: func() (config.Config, error) {
				return config.Config{}, err
			},
		},
		StateRepositoryQuery: func() (state.State, error) { return state.NewStateDisabled(), nil },
	}

	expectedEvent := StateRetrievalFailed{Reason: err}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestStatusShouldNotBeRetrievedDueToStateRetrievalError(t *testing.T) {
	err := errors.New("failed to retrieve state")

	deps := Dependencies{
		ConfigReader: &configReaderMock{
			read: func() (config.Config, error) {
				return config.Config{ActivationScope: activationscope.Global}, nil
			},
		},
		StateRepositoryQuery: func() (state.State, error) { return state.State{}, err },
	}

	expectedEvent := StateRetrievalFailed{Reason: err}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
