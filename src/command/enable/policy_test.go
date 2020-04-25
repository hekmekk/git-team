package enable

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	commitsettings "github.com/hekmekk/git-team/src/command/enable/commitsettings/entity"
	activationscope "github.com/hekmekk/git-team/src/shared/config/entity/activationscope"
	config "github.com/hekmekk/git-team/src/shared/config/entity/config"
	"github.com/hekmekk/git-team/src/shared/gitconfig/scope"
	gitconfigscope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
)

type commitSettingsReaderMock struct {
	read func() commitsettings.CommitSettings
}

func (mock commitSettingsReaderMock) Read() commitsettings.CommitSettings {
	return mock.read()
}

var commitSettings = commitsettings.CommitSettings{TemplatesBaseDir: "/path/to/commit-templates", HooksDir: "/path/to/hooks"}

var commitSettingsReader = commitSettingsReaderMock{
	read: func() commitsettings.CommitSettings {
		return commitSettings
	},
}

type configReaderMock struct {
	read func() (config.Config, error)
}

func (mock configReaderMock) Read() (config.Config, error) {
	return mock.read()
}

type gitConfigWriterMock struct {
	replaceAll func(scope.Scope, string, string) error
}

func (mock gitConfigWriterMock) UnsetAll(scope scope.Scope, key string) error {
	return nil
}

func (mock gitConfigWriterMock) ReplaceAll(scope scope.Scope, key string, value string) error {
	return mock.replaceAll(scope, key, value)
}

func (mock gitConfigWriterMock) Add(scope scope.Scope, key string, value string) error {
	return nil
}

type stateWriterMock struct {
	persistEnabled func(activationscope.ActivationScope, []string) error
}

func (mock stateWriterMock) PersistEnabled(scope activationscope.ActivationScope, coauthors []string) error {
	return mock.persistEnabled(scope, coauthors)
}
func (mock stateWriterMock) PersistDisabled(scope activationscope.ActivationScope) error {
	return nil
}

type activationValidatorMock struct {
	isInsideAGitRepository func() bool
}

func (mock activationValidatorMock) IsInsideAGitRepository() bool {
	return mock.isInsideAGitRepository()
}

// TODO: add parameterized case to ensure diff commit template folders for diff users and repos each

func TestEnableAborted(t *testing.T) {
	deps := Dependencies{}
	req := Request{AliasesAndCoauthors: &[]string{}}

	expectedEvent := Aborted{}

	event := Policy{deps, req}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestEnableSucceeds(t *testing.T) {
	t.Parallel()

	coauthors := &[]string{"Mr. Noujz <noujz@mr.se>", "mrs"}
	expectedStateRepositoryPersistEnabledCoauthors := []string{"Mr. Noujz <noujz@mr.se>", "Mrs. Noujz <noujz@mrs.se>"}
	expectedCommitTemplateCoauthors := "\n\nCo-authored-by: Mr. Noujz <noujz@mr.se>\nCo-authored-by: Mrs. Noujz <noujz@mrs.se>"

	WriteTemplateFile := func(_ string, data []byte, _ os.FileMode) error {
		if expectedCommitTemplateCoauthors != string(data) {
			t.Errorf("expected: %s, got: %s", expectedCommitTemplateCoauthors, string(data))
			t.Fail()
		}
		return nil
	}
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }

	user := "someone"
	pathToRepo := "/path/to/repo"
	repoChecksum := "a05b508ef212dc640aced4037a66021d" // echo -n <user>:<pathToRepo> | md5sum | awk '{ print $1 }'

	cases := []struct {
		activationScope activationscope.ActivationScope
		templateDir     string
		gitconfigScope  gitconfigscope.Scope
	}{
		{activationscope.Global, fmt.Sprintf("%s/global", commitSettings.TemplatesBaseDir), gitconfigscope.Global},
		{activationscope.RepoLocal, fmt.Sprintf("%s/repo-local/%s", commitSettings.TemplatesBaseDir, repoChecksum), gitconfigscope.Local},
	}

	deps := Dependencies{
		SanityCheckCoauthors: func(coauthors []string) []error { return []error{} },
		WriteTemplateFile:    WriteTemplateFile,
		GitResolveAliases:    resolveAliases,
		CommitSettingsReader: commitSettingsReader,
		GetEnv: func(string) string {
			return user
		},
		GetWd: func() (string, error) {
			return pathToRepo, nil
		},
		ActivationValidator: &activationValidatorMock{
			isInsideAGitRepository: func() bool {
				return true
			},
		},
	}

	for _, caseLoopVar := range cases {
		activationScope := caseLoopVar.activationScope
		expectedTemplateDir := caseLoopVar.templateDir
		expectedGitConfigScope := caseLoopVar.gitconfigScope

		t.Run(activationScope.String(), func(t *testing.T) {
			t.Parallel()

			deps.ConfigReader = &configReaderMock{
				read: func() (config.Config, error) {
					return config.Config{ActivationScope: activationScope}, nil
				},
			}

			deps.GitConfigWriter = &gitConfigWriterMock{
				replaceAll: func(scope scope.Scope, key string, value string) error {
					if key != "commit.template" && key != "core.hooksPath" {
						t.Errorf("wrong key: %s", key)
						t.Fail()
					}
					if scope != expectedGitConfigScope {
						t.Errorf("wrong scope, expected: %s, got: %s", expectedGitConfigScope, scope)
						t.Fail()
					}
					return nil
				},
			}

			deps.StateWriter = &stateWriterMock{
				persistEnabled: func(scope activationscope.ActivationScope, coauthors []string) error {
					if scope != activationScope {
						t.Errorf("wrong scope, expected: %s, got: %s", activationScope, scope)
						t.Fail()
					}
					if !reflect.DeepEqual(expectedStateRepositoryPersistEnabledCoauthors, coauthors) {
						t.Errorf("expected: %s, got: %s", expectedStateRepositoryPersistEnabledCoauthors, coauthors)
						t.Fail()
					}
					return nil
				},
			}

			deps.CreateTemplateDir = func(path string, _ os.FileMode) error {
				if path != expectedTemplateDir {
					t.Errorf("wrong path to template dir, expected: %s, got: %s", expectedTemplateDir, path)
					t.Fail()
				}
				return nil
			}

			req := Request{AliasesAndCoauthors: coauthors}

			expectedEvent := Succeeded{}

			event := Policy{deps, req}.Apply()

			if !reflect.DeepEqual(expectedEvent, event) {
				t.Errorf("expected: %s, got: %s", expectedEvent, event)
				t.Fail()
			}
		})
	}
}

func TestEnableDropsDuplicateEntries(t *testing.T) {
	coauthors := []string{"Mr. Noujz <noujz@mr.se>", "mrs", "Mr. Noujz <noujz@mr.se>", "mrs", "Mrs. Noujz <noujz@mrs.se>"}
	expectedStateRepositoryPersistEnabledCoauthors := []string{"Mr. Noujz <noujz@mr.se>", "Mrs. Noujz <noujz@mrs.se>"}
	expectedCommitTemplateCoauthors := "\n\nCo-authored-by: Mr. Noujz <noujz@mr.se>\nCo-authored-by: Mrs. Noujz <noujz@mrs.se>"

	CreateTemplateDir := func(string, os.FileMode) error { return nil }
	WriteTemplateFile := func(_ string, data []byte, _ os.FileMode) error {
		if expectedCommitTemplateCoauthors != string(data) {
			t.Errorf("expected: %s, got: %s", expectedCommitTemplateCoauthors, string(data))
			t.Fail()
		}
		return nil
	}
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }

	configReader := &configReaderMock{
		read: func() (config.Config, error) {
			return config.Config{ActivationScope: activationscope.Global}, nil
		},
	}

	gitConfigWriter := &gitConfigWriterMock{
		replaceAll: func(gitconfigscope.Scope, string, string) error {
			return nil
		},
	}

	stateWriter := &stateWriterMock{
		persistEnabled: func(_ activationscope.ActivationScope, coauthors []string) error {
			if !reflect.DeepEqual(expectedStateRepositoryPersistEnabledCoauthors, coauthors) {
				t.Errorf("expected: %s, got: %s", expectedStateRepositoryPersistEnabledCoauthors, coauthors)
				t.Fail()
			}
			return nil
		},
	}

	activationValidator := &activationValidatorMock{
		isInsideAGitRepository: func() bool {
			return true
		},
	}

	deps := Dependencies{
		SanityCheckCoauthors: func(coauthors []string) []error { return []error{} },
		CreateTemplateDir:    CreateTemplateDir,
		WriteTemplateFile:    WriteTemplateFile,
		GitResolveAliases:    resolveAliases,
		StateWriter:          stateWriter,
		CommitSettingsReader: commitSettingsReader,
		ConfigReader:         configReader,
		GitConfigWriter:      gitConfigWriter,
		ActivationValidator:  activationValidator,
	}
	req := Request{AliasesAndCoauthors: &coauthors}

	expectedEvent := Succeeded{}

	event := Policy{deps, req}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestEnableFailsDueToSanityCheckErr(t *testing.T) {
	coauthors := []string{"INVALID COAUTHOR"}

	expectedErr := errors.New("Not a valid coauthor: INVALID COAUTHOR")

	deps := Dependencies{
		SanityCheckCoauthors: func(coauthors []string) []error { return []error{expectedErr} },
	}
	req := Request{AliasesAndCoauthors: &coauthors}

	expectedEvent := Failed{Reason: []error{expectedErr}}

	event := Policy{deps, req}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestEnableFailsDueToResolveAliasesErr(t *testing.T) {
	coauthors := []string{"Mr. Noujz <noujz@mr.se>", "mrs"}

	expectedErr := errors.New("failed to resolve alias mrs")

	deps := Dependencies{
		SanityCheckCoauthors: func(coauthors []string) []error { return []error{} },
		GitResolveAliases:    func([]string) ([]string, []error) { return []string{}, []error{expectedErr} },
	}
	req := Request{AliasesAndCoauthors: &coauthors}

	expectedEvent := Failed{Reason: []error{expectedErr}}

	event := Policy{deps, req}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestEnableFailsDueToConfigReaderError(t *testing.T) {
	coauthors := &[]string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to read from config")

	sanityCheck := func([]string) []error { return []error{} }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	configReader := &configReaderMock{
		read: func() (config.Config, error) {
			return config.Config{}, expectedErr
		},
	}

	deps := Dependencies{
		SanityCheckCoauthors: sanityCheck,
		GitResolveAliases:    resolveAliases,
		CommitSettingsReader: commitSettingsReader,
		ConfigReader:         configReader,
	}

	req := Request{AliasesAndCoauthors: coauthors}

	expectedEvent := Failed{Reason: []error{expectedErr}}

	event := Policy{deps, req}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestEnableFailsWhileTryingToEnableOutsideOfAGitRepository(t *testing.T) {
	coauthors := &[]string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to enable with scope=repo-local: not inside a git repository")

	sanityCheck := func([]string) []error { return []error{} }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	CreateTemplateDir := func(string, os.FileMode) error { return nil }
	WriteTemplateFile := func(string, []byte, os.FileMode) error { return nil }

	configReader := &configReaderMock{
		read: func() (config.Config, error) {
			return config.Config{ActivationScope: activationscope.RepoLocal}, nil
		},
	}

	gitConfigWriter := &gitConfigWriterMock{
		replaceAll: func(_ gitconfigscope.Scope, key string, _ string) error {
			return nil
		},
	}

	activationValidator := &activationValidatorMock{
		isInsideAGitRepository: func() bool {
			return false
		},
	}

	deps := Dependencies{
		SanityCheckCoauthors: sanityCheck,
		CreateTemplateDir:    CreateTemplateDir,
		WriteTemplateFile:    WriteTemplateFile,
		GitResolveAliases:    resolveAliases,
		CommitSettingsReader: commitSettingsReader,
		ConfigReader:         configReader,
		GitConfigWriter:      gitConfigWriter,
		ActivationValidator:  activationValidator,
	}
	req := Request{AliasesAndCoauthors: coauthors}

	expectedEvent := Failed{Reason: []error{expectedErr}}

	event := Policy{deps, req}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestEnableFailsDueToGetWdErr(t *testing.T) {
	coauthors := []string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to get working dir")

	sanityCheck := func([]string) []error { return []error{} }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	CreateTemplateDir := func(string, os.FileMode) error { return nil }
	configReader := &configReaderMock{
		read: func() (config.Config, error) {
			return config.Config{ActivationScope: activationscope.RepoLocal}, nil
		},
	}
	activationValidator := &activationValidatorMock{
		isInsideAGitRepository: func() bool {
			return true
		},
	}

	deps := Dependencies{
		SanityCheckCoauthors: sanityCheck,
		CreateTemplateDir:    CreateTemplateDir,
		GitResolveAliases:    resolveAliases,
		CommitSettingsReader: commitSettingsReader,
		ConfigReader:         configReader,
		GetEnv: func(string) string {
			return "USER"
		},
		GetWd: func() (string, error) {
			return "", expectedErr
		},
		ActivationValidator: activationValidator,
	}
	req := Request{AliasesAndCoauthors: &coauthors}

	expectedEvent := Failed{Reason: []error{expectedErr}}

	event := Policy{deps, req}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestEnableFailsDueToCreateTemplateDirErr(t *testing.T) {
	coauthors := []string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to create Dir")

	sanityCheck := func([]string) []error { return []error{} }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	CreateTemplateDir := func(string, os.FileMode) error { return expectedErr }
	configReader := &configReaderMock{
		read: func() (config.Config, error) {
			return config.Config{ActivationScope: activationscope.Global}, nil
		},
	}
	activationValidator := &activationValidatorMock{
		isInsideAGitRepository: func() bool {
			return true
		},
	}

	deps := Dependencies{
		SanityCheckCoauthors: sanityCheck,
		CreateTemplateDir:    CreateTemplateDir,
		GitResolveAliases:    resolveAliases,
		CommitSettingsReader: commitSettingsReader,
		ConfigReader:         configReader,
		ActivationValidator:  activationValidator,
	}
	req := Request{AliasesAndCoauthors: &coauthors}

	expectedEvent := Failed{Reason: []error{expectedErr}}

	event := Policy{deps, req}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestEnableFailsDueToWriteTemplateFileErr(t *testing.T) {
	coauthors := &[]string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to write file")

	sanityCheck := func([]string) []error { return []error{} }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	CreateTemplateDir := func(string, os.FileMode) error { return nil }
	WriteTemplateFile := func(string, []byte, os.FileMode) error { return expectedErr }
	configReader := &configReaderMock{
		read: func() (config.Config, error) {
			return config.Config{ActivationScope: activationscope.Global}, nil
		},
	}
	activationValidator := &activationValidatorMock{
		isInsideAGitRepository: func() bool {
			return true
		},
	}

	deps := Dependencies{
		SanityCheckCoauthors: sanityCheck,
		CreateTemplateDir:    CreateTemplateDir,
		WriteTemplateFile:    WriteTemplateFile,
		GitResolveAliases:    resolveAliases,
		CommitSettingsReader: commitSettingsReader,
		ConfigReader:         configReader,
		ActivationValidator:  activationValidator,
	}
	req := Request{AliasesAndCoauthors: coauthors}

	expectedEvent := Failed{Reason: []error{expectedErr}}

	event := Policy{deps, req}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestEnableFailsDueToGitSetCommitTemplateErr(t *testing.T) {
	coauthors := &[]string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to set commit template")

	sanityCheck := func([]string) []error { return []error{} }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	CreateTemplateDir := func(string, os.FileMode) error { return nil }
	WriteTemplateFile := func(string, []byte, os.FileMode) error { return nil }

	configReader := &configReaderMock{
		read: func() (config.Config, error) {
			return config.Config{ActivationScope: activationscope.Global}, nil
		},
	}

	gitConfigWriter := &gitConfigWriterMock{
		replaceAll: func(_ gitconfigscope.Scope, key string, _ string) error {
			if key == "commit.template" {
				return expectedErr
			}
			return nil
		},
	}

	activationValidator := &activationValidatorMock{
		isInsideAGitRepository: func() bool {
			return true
		},
	}

	deps := Dependencies{
		SanityCheckCoauthors: sanityCheck,
		CreateTemplateDir:    CreateTemplateDir,
		WriteTemplateFile:    WriteTemplateFile,
		GitResolveAliases:    resolveAliases,
		CommitSettingsReader: commitSettingsReader,
		ConfigReader:         configReader,
		GitConfigWriter:      gitConfigWriter,
		ActivationValidator:  activationValidator,
	}
	req := Request{AliasesAndCoauthors: coauthors}

	expectedEvent := Failed{Reason: []error{expectedErr}}

	event := Policy{deps, req}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestEnableFailsDueToSetHooksPathErr(t *testing.T) {
	coauthors := &[]string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to set hooks path")

	sanityCheck := func([]string) []error { return []error{} }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	CreateTemplateDir := func(string, os.FileMode) error { return nil }
	WriteTemplateFile := func(string, []byte, os.FileMode) error { return nil }
	configReader := &configReaderMock{
		read: func() (config.Config, error) {
			return config.Config{ActivationScope: activationscope.Global}, nil
		},
	}

	gitConfigWriter := &gitConfigWriterMock{
		replaceAll: func(_ gitconfigscope.Scope, key string, _ string) error {
			if key == "core.hooksPath" {
				return expectedErr
			}
			return nil
		},
	}

	activationValidator := &activationValidatorMock{
		isInsideAGitRepository: func() bool {
			return true
		},
	}

	deps := Dependencies{
		SanityCheckCoauthors: sanityCheck,
		CreateTemplateDir:    CreateTemplateDir,
		WriteTemplateFile:    WriteTemplateFile,
		GitResolveAliases:    resolveAliases,
		CommitSettingsReader: commitSettingsReader,
		ConfigReader:         configReader,
		GitConfigWriter:      gitConfigWriter,
		ActivationValidator:  activationValidator,
	}
	req := Request{AliasesAndCoauthors: coauthors}

	expectedEvent := Failed{Reason: []error{expectedErr}}

	event := Policy{deps, req}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestEnableFailsDueToSaveStatusErr(t *testing.T) {
	coauthors := &[]string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to set status")

	sanityCheck := func([]string) []error { return []error{} }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	CreateTemplateDir := func(string, os.FileMode) error { return nil }
	WriteTemplateFile := func(string, []byte, os.FileMode) error { return nil }

	configReader := &configReaderMock{
		read: func() (config.Config, error) {
			return config.Config{ActivationScope: activationscope.Global}, nil
		},
	}

	gitConfigWriter := &gitConfigWriterMock{
		replaceAll: func(_ gitconfigscope.Scope, key string, _ string) error {
			return nil
		},
	}

	stateWriter := &stateWriterMock{
		persistEnabled: func(_ activationscope.ActivationScope, _ []string) error {
			return expectedErr
		},
	}

	activationValidator := &activationValidatorMock{
		isInsideAGitRepository: func() bool {
			return true
		},
	}

	deps := Dependencies{
		SanityCheckCoauthors: sanityCheck,
		CreateTemplateDir:    CreateTemplateDir,
		WriteTemplateFile:    WriteTemplateFile,
		GitResolveAliases:    resolveAliases,
		CommitSettingsReader: commitSettingsReader,
		ConfigReader:         configReader,
		GitConfigWriter:      gitConfigWriter,
		StateWriter:          stateWriter,
		ActivationValidator:  activationValidator,
	}

	req := Request{AliasesAndCoauthors: coauthors}

	expectedEvent := Failed{Reason: []error{expectedErr}}

	event := Policy{deps, req}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
