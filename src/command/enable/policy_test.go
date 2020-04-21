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

	setHooksPath := func(string) error { return nil }
	CreateTemplateDir := func(string, os.FileMode) error { return nil }
	WriteTemplateFile := func(_ string, data []byte, _ os.FileMode) error {
		if expectedCommitTemplateCoauthors != string(data) {
			t.Errorf("expected: %s, got: %s", expectedCommitTemplateCoauthors, string(data))
			t.Fail()
		}
		return nil
	}
	GitSetCommitTemplate := func(string) error { return nil }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	StateRepositoryPersistEnabled := func(coauthors []string) error {
		if !reflect.DeepEqual(expectedStateRepositoryPersistEnabledCoauthors, coauthors) {
			t.Errorf("expected: %s, got: %s", expectedStateRepositoryPersistEnabledCoauthors, coauthors)
			t.Fail()
		}
		return nil
	}

	cases := []struct {
		activationScope activationscope.ActivationScope
		templateDir     string
		gitconfigScope  gitconfigscope.Scope
	}{
		{activationscope.Global, fmt.Sprintf("%s/global", commitSettings.TemplatesBaseDir), gitconfigscope.Global},
		{activationscope.RepoLocal, fmt.Sprintf("%s/repo-local/<hash>", commitSettings.TemplatesBaseDir), gitconfigscope.Local},
	}

	deps := Dependencies{
		SanityCheckCoauthors:          func(coauthors []string) []error { return []error{} },
		CreateTemplateDir:             CreateTemplateDir,
		WriteTemplateFile:             WriteTemplateFile,
		GitSetHooksPath:               setHooksPath,
		GitSetCommitTemplate:          GitSetCommitTemplate,
		GitResolveAliases:             resolveAliases,
		StateRepositoryPersistEnabled: StateRepositoryPersistEnabled,
		CommitSettingsReader:          commitSettingsReader,
	}

	for _, caseLoopVar := range cases {
		activationScope := caseLoopVar.activationScope
		// templateDir := caseLoopVar.templateDir
		// gitconfigScope := caseLoopVar.gitconfigScope

		t.Run(activationScope.String(), func(t *testing.T) {
			t.Parallel()

			deps.ConfigReader = &configReaderMock{
				read: func() (config.Config, error) {
					return config.Config{ActivationScope: activationScope}, nil
				},
			}

			// TODO: add assertions for functions depending on global or repo-local (see acceptance test assertions)
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

	setHooksPath := func(string) error { return nil }
	CreateTemplateDir := func(string, os.FileMode) error { return nil }
	WriteTemplateFile := func(_ string, data []byte, _ os.FileMode) error {
		if expectedCommitTemplateCoauthors != string(data) {
			t.Errorf("expected: %s, got: %s", expectedCommitTemplateCoauthors, string(data))
			t.Fail()
		}
		return nil
	}
	GitSetCommitTemplate := func(string) error { return nil }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	StateRepositoryPersistEnabled := func(coauthors []string) error {
		if !reflect.DeepEqual(expectedStateRepositoryPersistEnabledCoauthors, coauthors) {
			t.Errorf("expected: %s, got: %s", expectedStateRepositoryPersistEnabledCoauthors, coauthors)
			t.Fail()
		}
		return nil
	}
	configReader := &configReaderMock{
		read: func() (config.Config, error) {
			return config.Config{ActivationScope: activationscope.Global}, nil
		},
	}

	deps := Dependencies{
		SanityCheckCoauthors:          func(coauthors []string) []error { return []error{} },
		CreateTemplateDir:             CreateTemplateDir,
		WriteTemplateFile:             WriteTemplateFile,
		GitSetHooksPath:               setHooksPath,
		GitSetCommitTemplate:          GitSetCommitTemplate,
		GitResolveAliases:             resolveAliases,
		StateRepositoryPersistEnabled: StateRepositoryPersistEnabled,
		CommitSettingsReader:          commitSettingsReader,
		ConfigReader:                  configReader,
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

	deps := Dependencies{
		SanityCheckCoauthors: sanityCheck,
		CreateTemplateDir:    CreateTemplateDir,
		GitResolveAliases:    resolveAliases,
		CommitSettingsReader: commitSettingsReader,
		ConfigReader:         configReader,
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

	deps := Dependencies{
		SanityCheckCoauthors: sanityCheck,
		CreateTemplateDir:    CreateTemplateDir,
		WriteTemplateFile:    WriteTemplateFile,
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

func TestEnableFailsDueToGitSetCommitTemplateErr(t *testing.T) {
	coauthors := &[]string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to set commit template")

	sanityCheck := func([]string) []error { return []error{} }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	CreateTemplateDir := func(string, os.FileMode) error { return nil }
	WriteTemplateFile := func(string, []byte, os.FileMode) error { return nil }
	GitSetCommitTemplate := func(string) error { return expectedErr }
	configReader := &configReaderMock{
		read: func() (config.Config, error) {
			return config.Config{ActivationScope: activationscope.Global}, nil
		},
	}

	deps := Dependencies{
		SanityCheckCoauthors: sanityCheck,
		CreateTemplateDir:    CreateTemplateDir,
		WriteTemplateFile:    WriteTemplateFile,
		GitSetCommitTemplate: GitSetCommitTemplate,
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

func TestEnableFailsDueToSetHooksPathErr(t *testing.T) {
	coauthors := &[]string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to set hooks path")

	sanityCheck := func([]string) []error { return []error{} }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	CreateTemplateDir := func(string, os.FileMode) error { return nil }
	WriteTemplateFile := func(string, []byte, os.FileMode) error { return nil }
	GitSetCommitTemplate := func(string) error { return nil }
	setHooksPath := func(string) error { return expectedErr }
	configReader := &configReaderMock{
		read: func() (config.Config, error) {
			return config.Config{ActivationScope: activationscope.Global}, nil
		},
	}

	deps := Dependencies{
		SanityCheckCoauthors: sanityCheck,
		CreateTemplateDir:    CreateTemplateDir,
		WriteTemplateFile:    WriteTemplateFile,
		GitSetHooksPath:      setHooksPath,
		GitSetCommitTemplate: GitSetCommitTemplate,
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

func TestEnableFailsDueToSaveStatusErr(t *testing.T) {
	coauthors := &[]string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to set status")

	sanityCheck := func([]string) []error { return []error{} }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	CreateTemplateDir := func(string, os.FileMode) error { return nil }
	WriteTemplateFile := func(string, []byte, os.FileMode) error { return nil }
	GitSetCommitTemplate := func(string) error { return nil }
	setHooksPath := func(string) error { return nil }
	configReader := &configReaderMock{
		read: func() (config.Config, error) {
			return config.Config{ActivationScope: activationscope.Global}, nil
		},
	}

	deps := Dependencies{
		SanityCheckCoauthors:          sanityCheck,
		CreateTemplateDir:             CreateTemplateDir,
		WriteTemplateFile:             WriteTemplateFile,
		GitSetHooksPath:               setHooksPath,
		GitSetCommitTemplate:          GitSetCommitTemplate,
		GitResolveAliases:             resolveAliases,
		CommitSettingsReader:          commitSettingsReader,
		ConfigReader:                  configReader,
		StateRepositoryPersistEnabled: func([]string) error { return expectedErr },
	}

	req := Request{AliasesAndCoauthors: coauthors}

	expectedEvent := Failed{Reason: []error{expectedErr}}

	event := Policy{deps, req}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
