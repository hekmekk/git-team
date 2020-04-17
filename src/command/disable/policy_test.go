package disable

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/shared/commitsettings/entity"
	commitsettingsentity "github.com/hekmekk/git-team/src/shared/commitsettings/entity"
)

type commitSettingsReaderMock struct {
	read func() entity.CommitSettings
}

func (mock commitSettingsReaderMock) Read() entity.CommitSettings {
	return mock.read()
}

var (
	unsetCommitTemplate  = func() error { return nil }
	unsetHooksPath       = func() error { return nil }
	fileInfo             os.FileInfo
	statFile             = func(string) (os.FileInfo, error) { return fileInfo, nil }
	removeFile           = func(string) error { return nil }
	persistDisabled      = func() error { return nil }
	commitSettings       = commitsettingsentity.CommitSettings{TemplatesBaseDir: "/path/to/commit-templates", HooksDir: "/path/to/hooks"}
	commitSettingsReader = commitSettingsReaderMock{
		read: func() entity.CommitSettings {
			return commitSettings
		},
	}
)

func TestDisableSucceeds(t *testing.T) {
	deps := Dependencies{
		GitUnsetHooksPath:      unsetHooksPath,
		GitUnsetCommitTemplate: unsetCommitTemplate,
		CommitSettingsReader:   commitSettingsReader,
		StatFile:               statFile,
		RemoveFile:             removeFile,
		PersistDisabled:        persistDisabled,
	}

	expectedEvent := Succeeded{}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestDisableShouldSucceedWhenUnsetHooksPathFailsBecauseTheOptionDoesntExist(t *testing.T) {
	expectedErr := errors.New("exit status 5")
	deps := Dependencies{
		GitUnsetHooksPath:      func() error { return expectedErr },
		GitUnsetCommitTemplate: unsetCommitTemplate,
		CommitSettingsReader:   commitSettingsReader,
		StatFile:               statFile,
		RemoveFile:             removeFile,
		PersistDisabled:        persistDisabled,
	}

	expectedEvent := Succeeded{}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestDisableShouldFailWhenUnsetHooksPathFails(t *testing.T) {
	err := errors.New("failed to unset hooks path")

	deps := Dependencies{
		GitUnsetHooksPath: func() error { return err },
	}

	expectedEvent := Failed{Reason: err}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestDisableShouldSucceedWhenUnsetCommitTemplateFailsBecauseItWasUnsetAlready(t *testing.T) {
	err := errors.New("exit status 5")

	deps := Dependencies{
		GitUnsetHooksPath:      unsetHooksPath,
		GitUnsetCommitTemplate: func() error { return err },
		CommitSettingsReader:   commitSettingsReader,
		StatFile:               statFile,
		RemoveFile:             removeFile,
		PersistDisabled:        persistDisabled,
	}

	expectedEvent := Succeeded{}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestDisableShouldFailWhenUnsetCommitTemplateFails(t *testing.T) {
	err := errors.New("failed to unset commit template")

	deps := Dependencies{
		GitUnsetHooksPath:      unsetHooksPath,
		GitUnsetCommitTemplate: func() error { return err },
	}

	expectedEvent := Failed{Reason: err}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestDisableShouldFailWhenRemoveFileFails(t *testing.T) {
	err := errors.New("failed to remove file")

	deps := Dependencies{
		GitUnsetHooksPath:      unsetHooksPath,
		GitUnsetCommitTemplate: unsetCommitTemplate,
		CommitSettingsReader:   commitSettingsReader,
		StatFile:               statFile,
		RemoveFile:             func(string) error { return err },
	}

	expectedEvent := Failed{Reason: err}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestDisableShouldSucceedButNotTryToRemoveTheCommitTemplateFileWhenStatFileFails(t *testing.T) {
	deps := Dependencies{
		GitUnsetHooksPath:      unsetHooksPath,
		GitUnsetCommitTemplate: unsetCommitTemplate,
		CommitSettingsReader:   commitSettingsReader,
		StatFile:               func(string) (os.FileInfo, error) { return fileInfo, errors.New("failed to stat file") },
		PersistDisabled:        persistDisabled,
	}

	expectedEvent := Succeeded{}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestDisableShouldFailWhenpersistDisabledFails(t *testing.T) {
	err := errors.New("failed to save status")

	deps := Dependencies{
		GitUnsetHooksPath:      unsetHooksPath,
		GitUnsetCommitTemplate: unsetCommitTemplate,
		CommitSettingsReader:   commitSettingsReader,
		StatFile:               statFile,
		RemoveFile:             removeFile,
		PersistDisabled:        func() error { return err },
	}

	expectedEvent := Failed{Reason: err}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
