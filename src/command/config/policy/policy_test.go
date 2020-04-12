package policy

import (
	"errors"
	"reflect"
	"testing"

	configevents "github.com/hekmekk/git-team/src/command/config/events"
	"github.com/hekmekk/git-team/src/core/config"
)

type configRepositoryMock struct {
	query func() (config.Config, error)
}

func (mock configRepositoryMock) Query() (config.Config, error) {
	return mock.query()
}

func TestConfigShouldBeRetrieved(t *testing.T) {
	cfg := config.Config{
		Ro: config.ReadOnlyProperties{
			GitTeamCommitTemplatePath: "/home/some-user/.config/git-team/COMMIT_TEMPLATE",
			GitTeamHooksPath:          "/usr/local/etc/git-team/hooks",
		},
		Rw: config.ReadWriteProperties{
			ActivationScope: config.Global,
		},
	}

	repo := configRepositoryMock{
		query: func() (config.Config, error) {
			return cfg, nil
		},
	}

	deps := Dependencies{
		ConfigRepository: repo,
	}

	expectedEvent := configevents.RetrievalSucceeded{Config: cfg}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestStatusShouldNotBeRetrieved(t *testing.T) {
	err := errors.New("failed to retrieve config")

	repo := configRepositoryMock{
		query: func() (config.Config, error) {
			return config.Config{}, err
		},
	}

	deps := Dependencies{
		ConfigRepository: repo,
	}

	expectedEvent := configevents.RetrievalFailed{Reason: err}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
