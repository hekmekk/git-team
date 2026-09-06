package status

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

	activationscope "github.com/hekmekk/git-team/v2/src/shared/activation/scope"
	config "github.com/hekmekk/git-team/v2/src/shared/config/entity/config"
	gitconfigerror "github.com/hekmekk/git-team/v2/src/shared/gitconfig/error"
	state "github.com/hekmekk/git-team/v2/src/shared/state/entity"
)

type configReaderMock struct {
	read func() (config.Config, error)
}

func (mock configReaderMock) Read() (config.Config, error) {
	return mock.read()
}

type stateReaderMock struct {
	query func(scope activationscope.Scope) (state.State, error)
}

func (mock stateReaderMock) Query(scope activationscope.Scope) (state.State, error) {
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

	cases := []activationscope.Scope{activationscope.Global, activationscope.RepoLocal}

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
					query: func(scope activationscope.Scope) (state.State, error) {
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
	deps := Dependencies{
		ConfigReader: &configReaderMock{
			read: func() (config.Config, error) {
				return config.Config{}, gitconfigerror.ErrConfigFileCannotBeWritten
			},
		},
	}

	expectedEvent := StateRetrievalFailed{Reason: fmt.Errorf("failed to read config: %s", gitconfigerror.ErrConfigFileCannotBeWritten)}

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

	expectedErr := errors.New("failed to get status with activation-scope=repo-local: not inside a git repository")

	expectedEvent := StateRetrievalFailed{Reason: expectedErr}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestStatusShouldNotBeRetrievedDueToStateRetrievalError(t *testing.T) {
	err := errors.New("no active co-authors found")

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
			query: func(scope activationscope.Scope) (state.State, error) {
				if scope != activationscope.Global {
					return state.State{}, errors.New("wrong scope")
				}
				return state.State{}, err
			},
		},
	}

	expectedEvent := StateRetrievalFailed{Reason: fmt.Errorf("failed to query current state: %s", err)}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
