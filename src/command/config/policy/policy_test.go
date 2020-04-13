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

type configWriterMock struct {
	setActivationScope func(scope activationscope.ActivationScope) error
}

func (mock configWriterMock) SetActivationScope(scope activationscope.ActivationScope) error {
	return mock.setActivationScope(scope)
}

func TestConfigShouldBeRetrieved(t *testing.T) {
	expectedEvent := configevents.RetrievalSucceeded{Config: cfg}

	configReader := configReaderMock{
		read: func() (config.Config, error) {
			return cfg, nil
		},
	}

	event := Policy{Req: Request{nil, nil}, Deps: Dependencies{configReader, nil}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestConfigShouldNotBeRetrieved(t *testing.T) {
	err := errors.New("failed to retrieve config")

	expectedEvent := configevents.RetrievalFailed{Reason: err}

	configReader := configReaderMock{
		read: func() (config.Config, error) {
			return config.Config{}, err
		},
	}

	emptyString := ""

	eventNil := Policy{Req: Request{nil, nil}, Deps: Dependencies{configReader, nil}}.Apply()
	eventEmptyString := Policy{Req: Request{&emptyString, &emptyString}, Deps: Dependencies{configReader, nil}}.Apply()

	if !reflect.DeepEqual(expectedEvent, eventNil) {
		t.Errorf("expected: %s, got: %s", expectedEvent, eventNil)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedEvent, eventEmptyString) {
		t.Errorf("expected: %s, got: %s", expectedEvent, eventEmptyString)
		t.Fail()
	}
}

func TestFailIfOnlyOneArgumentIsProvided(t *testing.T) {
	expectedEvent := configevents.ReadingSingleSettingNotYetImplemented{}

	key := "A"
	value := "B"
	emptyString := ""

	eventKeyOnlyNil := Policy{Req: Request{&key, nil}, Deps: Dependencies{nil, nil}}.Apply()
	eventKeyOnlyEmptyString := Policy{Req: Request{&key, &emptyString}, Deps: Dependencies{nil, nil}}.Apply()
	eventValueOnlyNil := Policy{Req: Request{nil, &value}, Deps: Dependencies{nil, nil}}.Apply()
	eventValueOnlyEmptyString := Policy{Req: Request{&emptyString, &value}, Deps: Dependencies{nil, nil}}.Apply()

	if !reflect.DeepEqual(expectedEvent, eventKeyOnlyNil) {
		t.Errorf("expected: %s, got: %s", expectedEvent, eventKeyOnlyNil)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedEvent, eventKeyOnlyEmptyString) {
		t.Errorf("expected: %s, got: %s", expectedEvent, eventKeyOnlyEmptyString)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedEvent, eventValueOnlyNil) {
		t.Errorf("expected: %s, got: %s", expectedEvent, eventValueOnlyNil)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedEvent, eventValueOnlyEmptyString) {
		t.Errorf("expected: %s, got: %s", expectedEvent, eventValueOnlyEmptyString)
		t.Fail()
	}
}

func TestFailOnUnknownSetting(t *testing.T) {
	err := errors.New("unknown setting 'A'")

	expectedEvent := configevents.SettingModificationFailed{Reason: err}

	key := "A"
	value := "B"
	event := Policy{Req: Request{&key, &value}, Deps: Dependencies{nil, nil}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestFailOnUnknownActivationScope(t *testing.T) {
	err := errors.New("unknown activation-scope 'A'")

	expectedEvent := configevents.SettingModificationFailed{Reason: err}

	configWriter := &configWriterMock{
		setActivationScope: func(scope activationscope.ActivationScope) error {
			return err
		},
	}

	key := "activation-scope"
	value := "A"
	event := Policy{Req: Request{&key, &value}, Deps: Dependencies{nil, configWriter}}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
