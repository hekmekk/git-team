package disable

import (
	"errors"
	"testing"

	"github.com/hekmekk/git-team/src/config"
)

var (
	unsetCommitTemplate = func() error { return nil }
	unsetHooksPath      = func() error { return nil }
	removeCommitSection = func() error { return nil }
	cfg                 = config.Config{TemplateFileName: "TEMPLATE_FILE", BaseDir: "BASE_DIR", StatusFileName: "STATUS_FILE"}
	loadConfig          = func() (config.Config, error) { return cfg, nil }
	removeFile          = func(string) error { return nil }
	persistDisabled     = func() error { return nil }
)

func TestDisableSucceeds(t *testing.T) {
	deps := dependencies{
		GitUnsetHooksPath:      unsetHooksPath,
		GitUnsetCommitTemplate: unsetCommitTemplate,
		LoadConfig:             loadConfig,
		RemoveFile:             removeFile,
		PersistDisabled:        persistDisabled,
	}
	err := executorFactory(deps)()

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestDisableShouldFailWhenUnsetHooksPathFails(t *testing.T) {
	expectedErr := errors.New("failed to unset hooks path")
	deps := dependencies{
		GitUnsetHooksPath: func() error { return expectedErr },
	}

	err := executorFactory(deps)()

	if err == nil || expectedErr != err {
		t.Errorf("expected: %s, received: %s", expectedErr, err)
		t.Fail()
	}
}

func TestDisableShouldSucceedWhenUnsetCommitTemplateFails(t *testing.T) {
	expectedErr := errors.New("failed to unset commit template")
	deps := dependencies{
		GitUnsetHooksPath:      unsetHooksPath,
		GitUnsetCommitTemplate: func() error { return expectedErr },
	}

	err := executorFactory(deps)()

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestDisableShouldFailWhenLoadConfigFails(t *testing.T) {
	expectedErr := errors.New("failed to load config")
	deps := dependencies{
		GitUnsetHooksPath:      unsetHooksPath,
		GitUnsetCommitTemplate: unsetCommitTemplate,
		LoadConfig:             func() (config.Config, error) { return config.Config{}, expectedErr },
	}

	err := executorFactory(deps)()

	if err == nil || expectedErr != err {
		t.Errorf("expected: %s, received: %s", expectedErr, err)
		t.Fail()
	}
}

func TestDisableShouldFailWhenRemoveFileFails(t *testing.T) {
	expectedErr := errors.New("failed to remove file")
	deps := dependencies{
		GitUnsetHooksPath:      unsetHooksPath,
		GitUnsetCommitTemplate: unsetCommitTemplate,
		LoadConfig:             loadConfig,
		RemoveFile:             func(string) error { return expectedErr },
	}

	err := executorFactory(deps)()

	if err == nil || expectedErr != err {
		t.Errorf("expected: %s, received: %s", expectedErr, err)
		t.Fail()
	}
}

func TestDisableShouldFailWhenpersistDisabledFails(t *testing.T) {
	expectedErr := errors.New("failed to save status")
	deps := dependencies{
		GitUnsetHooksPath:      unsetHooksPath,
		GitUnsetCommitTemplate: unsetCommitTemplate,
		LoadConfig:             loadConfig,
		RemoveFile:             removeFile,
		PersistDisabled:        func() error { return expectedErr },
	}

	err := executorFactory(deps)()

	if err == nil || expectedErr != err {
		t.Errorf("expected: %s, received: %s", expectedErr, err)
		t.Fail()
	}
}
