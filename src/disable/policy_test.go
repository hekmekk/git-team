package disable

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/core/config"
)

var (
	unsetCommitTemplate = func() error { return nil }
	unsetHooksPath      = func() error { return nil }
	cfg                 = config.Config{TemplateFileName: "TEMPLATE_FILE", BaseDir: "BASE_DIR", GitHooksPath: "GIT_HOOKSPATH"}
	loadConfig          = func() config.Config { return cfg }
	fileInfo            os.FileInfo
	statFile            = func(string) (os.FileInfo, error) { return fileInfo, nil }
	removeFile          = func(string) error { return nil }
	persistDisabled     = func() error { return nil }
)

func TestDisableSucceeds(t *testing.T) {
	deps := Dependencies{
		GitUnsetHooksPath:      unsetHooksPath,
		GitUnsetCommitTemplate: unsetCommitTemplate,
		LoadConfig:             loadConfig,
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
		LoadConfig:             loadConfig,
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
		LoadConfig:             loadConfig,
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
		LoadConfig:             loadConfig,
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
		LoadConfig:             loadConfig,
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
		LoadConfig:             loadConfig,
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
