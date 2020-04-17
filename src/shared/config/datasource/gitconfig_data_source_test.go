package datasource

import (
	"errors"
	"reflect"
	"testing"

	activationscope "github.com/hekmekk/git-team/src/shared/config/entity/activationscope"
	config "github.com/hekmekk/git-team/src/shared/config/entity/config"
)

type gitconfigRawReaderMock struct {
	get func(key string) (string, error)
}

func (mock gitconfigRawReaderMock) Get(key string) (string, error) {
	return mock.get(key)
}

func TestLoadSucceeds(t *testing.T) {
	t.Parallel()

	cases := []struct {
		rawGitConfigScope string
		expectedConfig    config.Config
	}{
		{"global", config.Config{ActivationScope: activationscope.Global}},
		{"repo-local", config.Config{ActivationScope: activationscope.RepoLocal}},
	}

	for _, caseLoopVar := range cases {
		rawGitConfigScope := caseLoopVar.rawGitConfigScope
		expectedCfg := caseLoopVar.expectedConfig

		t.Run(rawGitConfigScope, func(t *testing.T) {
			t.Parallel()

			rawReader := gitconfigRawReaderMock{
				get: func(key string) (string, error) {
					if key != "team.config.activation-scope" {
						return "", errors.New("wrong key")
					}
					return rawGitConfigScope, nil
				},
			}
			cfg, _ := newGitconfigDataSource(rawReader).Read()

			if !reflect.DeepEqual(expectedCfg, cfg) {
				t.Errorf("expected: %s, received %s", expectedCfg, cfg)
				t.Fail()
			}
		})
	}
}

func TestLoadSucceedsWhenActivationScopeIsUnsetInGitconfig(t *testing.T) {
	t.Parallel()

	expectedCfg := config.Config{ActivationScope: activationscope.Global}

	rawReader := gitconfigRawReaderMock{
		get: func(key string) (string, error) {
			if key != "team.config.activation-scope" {
				return "", errors.New("wrong key")
			}
			return "", errors.New("exit status 1")
		},
	}

	cfg, err := newGitconfigDataSource(rawReader).Read()

	if err != nil {
		t.Errorf("expected no error, received: %s", err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedCfg, cfg) {
		t.Errorf("expected: %s, received %s", expectedCfg, cfg)
		t.Fail()
	}
}

func TestLoadFailsWhenReadingFromGitconfigFails(t *testing.T) {
	t.Parallel()
	rawReader := gitconfigRawReaderMock{
		get: func(key string) (string, error) {
			if key != "team.config.activation-scope" {
				return "", errors.New("wrong key")
			}
			return "", errors.New("reading from gitconfig failed")
		},
	}
	cfg, err := newGitconfigDataSource(rawReader).Read()

	if err == nil {
		t.Errorf("expected error, received %s", cfg)
		t.Fail()
	}
}

func TestLoadFailsWhenActivationScopeFromGitconfigIsUnknown(t *testing.T) {
	t.Parallel()
	rawReader := gitconfigRawReaderMock{
		get: func(key string) (string, error) {
			if key != "team.config.activation-scope" {
				return "", errors.New("wrong key")
			}
			return "unknown", nil
		},
	}
	cfg, err := newGitconfigDataSource(rawReader).Read()

	if err == nil {
		t.Errorf("expected error, received %s", cfg)
		t.Fail()
	}
}
