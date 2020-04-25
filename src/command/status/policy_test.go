package status

import (
	"errors"
	"reflect"
	"testing"

	activationscope "github.com/hekmekk/git-team/src/shared/config/entity/activationscope"
	config "github.com/hekmekk/git-team/src/shared/config/entity/config"
	state "github.com/hekmekk/git-team/src/shared/state/entity"
)

type configReaderMock struct {
	read func() (config.Config, error)
}

func (mock configReaderMock) Read() (config.Config, error) {
	return mock.read()
}

type stateReaderMock struct {
	query func(scope activationscope.ActivationScope) (state.State, error)
}

func (mock stateReaderMock) Query(scope activationscope.ActivationScope) (state.State, error) {
	return mock.query(scope)
}

type activationValidatorMock struct {
	isInsideAGitRepository func() bool
}

func (mock activationValidatorMock) IsInsideAGitRepository() bool {
	return mock.isInsideAGitRepository()
}

func TestStatusShouldBeRetrieved(t *testing.T) {
	t.Parallel()

	currState := state.NewStateDisabled()

	cases := []activationscope.ActivationScope{activationscope.Global, activationscope.RepoLocal}

	for _, caseLoopVar := range cases {
		activationScope := caseLoopVar

		t.Run(activationScope.String(), func(t *testing.T) {
			t.Parallel()

			deps := Dependencies{
				ConfigReader: &configReaderMock{
					read: func() (config.Config, error) {
						return config.Config{ActivationScope: activationScope}, nil
					},
				},
				ActivationValidator: &activationValidatorMock{
					isInsideAGitRepository: func() bool {
						return true
					},
				},
				StateReader: &stateReaderMock{
					query: func(scope activationscope.ActivationScope) (state.State, error) {
						if scope != activationScope {
							return state.State{}, errors.New("wrong scope")
						}
						return currState, nil
					},
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
	}

	expectedEvent := StateRetrievalFailed{Reason: err}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestStatusShouldFailWhenNotInsideAGitRepository(t *testing.T) {
	deps := Dependencies{
		ConfigReader: &configReaderMock{
			read: func() (config.Config, error) {
				return config.Config{ActivationScope: activationscope.RepoLocal}, nil
			},
		},
		ActivationValidator: &activationValidatorMock{
			isInsideAGitRepository: func() bool {
				return false
			},
		},
	}

	expectedErr := errors.New("Failed to get status with scope=repo-local: not inside a git repository")

	expectedEvent := StateRetrievalFailed{Reason: expectedErr}

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
		ActivationValidator: &activationValidatorMock{
			isInsideAGitRepository: func() bool {
				return true
			},
		},
		StateReader: &stateReaderMock{
			query: func(scope activationscope.ActivationScope) (state.State, error) {
				if scope != activationscope.Global {
					return state.State{}, errors.New("wrong scope")
				}
				return state.State{}, err
			},
		},
	}

	expectedEvent := StateRetrievalFailed{Reason: err}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
