package datasource

import (
	"fmt"
	"reflect"
	"testing"

	activationscope "github.com/hekmekk/git-team/v2/src/shared/activation/scope"
	config "github.com/hekmekk/git-team/v2/src/shared/config/entity/config"
	gitconfigerror "github.com/hekmekk/git-team/v2/src/shared/gitconfig/error"
	gitconfigscope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
)

type gitConfigReaderMock struct {
	get func(gitconfigscope.Scope, string) (string, error)
}

func (mock gitConfigReaderMock) Get(scope gitconfigscope.Scope, key string) (string, error) {
	return mock.get(scope, key)
}

func (mock gitConfigReaderMock) GetAll(scope gitconfigscope.Scope, key string) ([]string, error) {
	return []string{}, nil
}

func (mock gitConfigReaderMock) GetRegexp(scope gitconfigscope.Scope, pattern string) (map[string]string, error) {
	return nil, nil
}

func (mock gitConfigReaderMock) List(scope gitconfigscope.Scope) (map[string]string, error) {
	return nil, nil
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

			gitConfigReader := gitConfigReaderMock{
				get: func(scope gitconfigscope.Scope, key string) (string, error) {
					if scope != gitconfigscope.Global {
						return "", fmt.Errorf("wrong scope: %s", scope)
					}
					if key != "team.config.activation-scope" {
						return "", fmt.Errorf("wrong key: %s", key)
					}
					return rawGitConfigScope, nil
				},
			}
			cfg, _ := NewGitconfigDataSource(gitConfigReader).Read()

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

	gitConfigReader := gitConfigReaderMock{
		get: func(scope gitconfigscope.Scope, key string) (string, error) {
			if scope != gitconfigscope.Global {
				return "", fmt.Errorf("wrong scope: %s", scope)
			}
			if key != "team.config.activation-scope" {
				return "", fmt.Errorf("wrong key: %s", key)
			}
			return "", gitconfigerror.ErrSectionOrKeyIsInvalid
		},
	}

	cfg, err := NewGitconfigDataSource(gitConfigReader).Read()

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
	gitConfigReader := gitConfigReaderMock{
		get: func(scope gitconfigscope.Scope, key string) (string, error) {
			if scope != gitconfigscope.Global {
				return "", fmt.Errorf("wrong scope: %s", scope)
			}
			if key != "team.config.activation-scope" {
				return "", fmt.Errorf("wrong key: %s", key)
			}
			return "", gitconfigerror.ErrConfigFileCannotBeWritten
		},
	}
	cfg, err := NewGitconfigDataSource(gitConfigReader).Read()

	if err == nil {
		t.Errorf("expected error, received %s", cfg)
		t.Fail()
	}
}

func TestLoadFailsWhenActivationScopeFromGitconfigIsUnknown(t *testing.T) {
	t.Parallel()
	gitConfigReader := gitConfigReaderMock{
		get: func(scope gitconfigscope.Scope, key string) (string, error) {
			if scope != gitconfigscope.Global {
				return "", fmt.Errorf("wrong scope: %s", scope)
			}
			if key != "team.config.activation-scope" {
				return "", fmt.Errorf("wrong key: %s", key)
			}
			return "unknown", nil
		},
	}
	cfg, err := NewGitconfigDataSource(gitConfigReader).Read()

	if err == nil {
		t.Errorf("expected error, received %s", cfg)
		t.Fail()
	}
}
