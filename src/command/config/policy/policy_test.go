package policy

import (
	"errors"
	"reflect"
	"testing"

	activationscope "github.com/hekmekk/git-team/src/command/config/entity/activationscope"
	config "github.com/hekmekk/git-team/src/command/config/entity/config"
	configevents "github.com/hekmekk/git-team/src/command/config/events"
)

type configReaderMock struct {
	read func() (config.Config, error)
}

func (mock configReaderMock) Read() (config.Config, error) {
	return mock.read()
}

var cfg = config.Config{ActivationScope: activationscope.Global}

func TestConfigShouldBeRetrieved(t *testing.T) {
	expectedEvent := configevents.RetrievalSucceeded{Config: cfg}

	configReader := configReaderMock{
		read: func() (config.Config, error) {
			return cfg, nil
		},
	}

	event := Policy{Dependencies{configReader}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestStatusShouldNotBeRetrieved(t *testing.T) {
	err := errors.New("failed to retrieve config")

	configReader := configReaderMock{
		read: func() (config.Config, error) {
			return config.Config{}, err
		},
	}

	expectedEvent := configevents.RetrievalFailed{Reason: err}

	event := Policy{Dependencies{configReader}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
