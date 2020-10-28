package enable

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	commitsettings "github.com/hekmekk/git-team/src/command/enable/commitsettings/entity"
	activationscope "github.com/hekmekk/git-team/src/shared/activation/scope"
	config "github.com/hekmekk/git-team/src/shared/config/entity/config"
	gitconfigscope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
)

type gitConfigReaderMock struct {
	getRegexp func(gitconfigscope.Scope, string) (map[string]string, error)
}

func (mock gitConfigReaderMock) Get(scope gitconfigscope.Scope, key string) (string, error) {
	return "", nil
}

func (mock gitConfigReaderMock) GetAll(scope gitconfigscope.Scope, key string) ([]string, error) {
	return []string{}, nil
}

func (mock gitConfigReaderMock) GetRegexp(scope gitconfigscope.Scope, pattern string) (map[string]string, error) {
	return mock.getRegexp(scope, pattern)
}

func (mock gitConfigReaderMock) List(scope gitconfigscope.Scope) (map[string]string, error) {
	return nil, nil
}

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
	replaceAll func(gitconfigscope.Scope, string, string) error
}

func (mock gitConfigWriterMock) UnsetAll(scope gitconfigscope.Scope, key string) error {
	return nil
}

func (mock gitConfigWriterMock) ReplaceAll(scope gitconfigscope.Scope, key string, value string) error {
	return mock.replaceAll(scope, key, value)
}

func (mock gitConfigWriterMock) Add(scope gitconfigscope.Scope, key string, value string) error {
	return nil
}

type stateWriterMock struct {
	persistEnabled func(activationscope.Scope, []string) error
}

func (mock stateWriterMock) PersistEnabled(scope activationscope.Scope, coauthors []string) error {
	return mock.persistEnabled(scope, coauthors)
}
func (mock stateWriterMock) PersistDisabled(scope activationscope.Scope) error {
	return nil
}

type activationValidatorMock struct {
	isInsideAGitRepository func() bool
}

func (mock activationValidatorMock) IsInsideAGitRepository() bool {
	return mock.isInsideAGitRepository()
}

func TestEnableAborted(t *testing.T) {
	deps := Dependencies{}
	req := Request{AliasesAndCoauthors: &[]string{}, UseAll: &[]bool{false}[0]}

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
		activationScope activationscope.Scope
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
				replaceAll: func(scope gitconfigscope.Scope, key string, value string) error {
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
				persistEnabled: func(scope activationscope.Scope, coauthors []string) error {
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

			req := Request{AliasesAndCoauthors: coauthors, UseAll: &[]bool{false}[0]}

			expectedEvent := Succeeded{}

			event := Policy{deps, req}.Apply()

			if !reflect.DeepEqual(expectedEvent, event) {
				t.Errorf("expected: %s, got: %s", expectedEvent, event)
				t.Fail()
			}
		})
	}
}

func TestEnableAllShouldSucceed(t *testing.T) {
	coauthors := &[]string{}
	expectedStateRepositoryPersistEnabledCoauthors := []string{"Mr. Noujz <noujz@mr.se>", "Mrs. Noujz <noujz@mrs.se>"}

	WriteTemplateFile := func(_ string, data []byte, _ os.FileMode) error {
		return nil
	}

	deps := Dependencies{
		GitConfigReader: &gitConfigReaderMock{
			getRegexp: func(_ gitconfigscope.Scope, pattern string) (map[string]string, error) {
				return map[string]string{
					"team.alias.alias1": "Mr. Noujz <noujz@mr.se>",
					"team.alias.alias2": "Mrs. Noujz <noujz@mrs.se>",
				}, nil
			},
		},

		WriteTemplateFile:    WriteTemplateFile,
		CommitSettingsReader: commitSettingsReader,
		GetEnv: func(string) string {
			return "someone"
		},
		GetWd: func() (string, error) {
			return "/path/to/repo", nil
		},
		ActivationValidator: &activationValidatorMock{
			isInsideAGitRepository: func() bool {
				return true
			},
		},
	}

	deps.ConfigReader = &configReaderMock{
		read: func() (config.Config, error) {
			return config.Config{ActivationScope: activationscope.Global}, nil
		},
	}

	deps.GitConfigWriter = &gitConfigWriterMock{
		replaceAll: func(scope gitconfigscope.Scope, key string, value string) error {
			return nil
		},
	}

	deps.StateWriter = &stateWriterMock{
		persistEnabled: func(scope activationscope.Scope, coauthors []string) error {
			if !reflect.DeepEqual(expectedStateRepositoryPersistEnabledCoauthors, coauthors) {
				t.Errorf("expected: %s, got: %s", expectedStateRepositoryPersistEnabledCoauthors, coauthors)
				t.Fail()
			}
			return nil
		},
	}

	deps.CreateTemplateDir = func(path string, _ os.FileMode) error {
		return nil
	}

	req := Request{AliasesAndCoauthors: coauthors, UseAll: &[]bool{true}[0]}

	expectedEvent := Succeeded{}

	event := Policy{deps, req}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func testEnableAllShouldFailWhenLookingUpCoauthorsReturnsAnError(t *testing.T) {
	err := errors.New("exit status 1")

	deps := Dependencies{
		GitConfigReader: &gitConfigReaderMock{
			getRegexp: func(_ gitconfigscope.Scope, pattern string) (map[string]string, error) {
				return map[string]string{}, err
			},
		},
	}

	req := Request{AliasesAndCoauthors: &[]string{}, UseAll: &[]bool{true}[0]}

	expectedEvent := Failed{Reason: []error{err}}

	event := Policy{deps, req}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestEnableAllShouldAbortWhenNoCoauthorsCouldBeFound(t *testing.T) {
	deps := Dependencies{
		GitConfigReader: &gitConfigReaderMock{
			getRegexp: func(_ gitconfigscope.Scope, pattern string) (map[string]string, error) {
				return map[string]string{}, nil
			},
		},
	}

	req := Request{AliasesAndCoauthors: &[]string{}, UseAll: &[]bool{true}[0]}

	expectedEvent := Aborted{}

	event := Policy{deps, req}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestEnableKeepsSeparateCommitTemplatesPerUserAndRepository(t *testing.T) {
	t.Parallel()

	deps := Dependencies{
		SanityCheckCoauthors: func(coauthors []string) []error { return []error{} },
		WriteTemplateFile: func(_ string, data []byte, _ os.FileMode) error {
			return nil
		},
		GitResolveAliases:    func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} },
		CommitSettingsReader: commitSettingsReader,
		ConfigReader: &configReaderMock{
			read: func() (config.Config, error) {
				return config.Config{ActivationScope: activationscope.RepoLocal}, nil
			},
		},
		ActivationValidator: &activationValidatorMock{
			isInsideAGitRepository: func() bool {
				return true
			},
		},
		GitConfigWriter: &gitConfigWriterMock{
			replaceAll: func(gitconfigscope.Scope, string, string) error {
				return nil
			},
		},
		StateWriter: &stateWriterMock{
			persistEnabled: func(activationscope.Scope, []string) error {
				return nil
			},
		},
	}

	properties := []struct {
		user             string
		pathToRepo       string
		expectedChecksum string // echo -n <user>:<pathToRepo> | md5sum | awk '{ print $1 }'
	}{
		{"someone", "/path/to/repo", "a05b508ef212dc640aced4037a66021d"},
		{"someone", "/path/to/another/repo", "0e272adf29c65b982f85ba4c034b5e3a"},
		{"someone-else", "/path/to/repo", "35e7efa80c9f6db4637267a1b03e2131"},
		{"someone-else", "/path/to/another/repo", "70271fa3dc8450ff89b0cc66ccf76411"},
	}

	for _, caseLoopVar := range properties {
		user := caseLoopVar.user
		pathToRepo := caseLoopVar.pathToRepo
		expectedChecksum := caseLoopVar.expectedChecksum

		t.Run(fmt.Sprintf("%s:%s", user, pathToRepo), func(t *testing.T) {
			t.Parallel()

			deps.GetEnv = func(string) string {
				return user
			}
			deps.GetWd = func() (string, error) {
				return pathToRepo, nil
			}

			expectedTemplateDir := fmt.Sprintf("/path/to/commit-templates/repo-local/%s", expectedChecksum)

			deps.CreateTemplateDir = func(path string, _ os.FileMode) error {
				if path != expectedTemplateDir {
					t.Errorf("wrong path to template dir, expected: %s, got: %s", expectedTemplateDir, path)
					t.Fail()
				}
				return nil
			}

			req := Request{AliasesAndCoauthors: &[]string{"Mr. Noujz <noujz@mr.se>", "mrs"}, UseAll: &[]bool{false}[0]}

			Policy{deps, req}.Apply()
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
		persistEnabled: func(_ activationscope.Scope, coauthors []string) error {
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
	req := Request{AliasesAndCoauthors: &coauthors, UseAll: &[]bool{false}[0]}

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
	req := Request{AliasesAndCoauthors: &coauthors, UseAll: &[]bool{false}[0]}

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
	req := Request{AliasesAndCoauthors: &coauthors, UseAll: &[]bool{false}[0]}

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

	req := Request{AliasesAndCoauthors: coauthors, UseAll: &[]bool{false}[0]}

	expectedEvent := Failed{Reason: []error{expectedErr}}

	event := Policy{deps, req}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestEnableFailsWhenNotInsideAGitRepository(t *testing.T) {
	coauthors := &[]string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to enable with activation-scope=repo-local: not inside a git repository")

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
	req := Request{AliasesAndCoauthors: coauthors, UseAll: &[]bool{false}[0]}

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
	req := Request{AliasesAndCoauthors: &coauthors, UseAll: &[]bool{false}[0]}

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
	req := Request{AliasesAndCoauthors: &coauthors, UseAll: &[]bool{false}[0]}

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
	req := Request{AliasesAndCoauthors: coauthors, UseAll: &[]bool{false}[0]}

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
	req := Request{AliasesAndCoauthors: coauthors, UseAll: &[]bool{false}[0]}

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
	req := Request{AliasesAndCoauthors: coauthors, UseAll: &[]bool{false}[0]}

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
		persistEnabled: func(_ activationscope.Scope, _ []string) error {
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

	req := Request{AliasesAndCoauthors: coauthors, UseAll: &[]bool{false}[0]}

	expectedEvent := Failed{Reason: []error{expectedErr}}

	event := Policy{deps, req}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
