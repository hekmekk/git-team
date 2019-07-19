package disable

import (
	"errors"
	"testing"

	"github.com/hekmekk/git-team/core/status"
	"github.com/hekmekk/git-team/src/config"
)

var (
	unsetCommitTemplate = func() error { return nil }
	removeCommitSection = func() error { return nil }
	cfg                 = config.Config{TemplateFileName: "TEMPLATE_FILE", BaseDir: "BASE_DIR", StatusFileName: "STATUS_FILE"}
	loadConfig          = func() (config.Config, error) { return cfg, nil }
	removeFile          = func(string) error { return nil }
	saveStatus          = func(status.State, ...string) error { return nil }
)

func TestDisableSucceeds(t *testing.T) {
	deps := dependencies{
		GitUnsetCommitTemplate: unsetCommitTemplate,
		GitRemoveCommitSection: removeCommitSection,
		LoadConfig:             loadConfig,
		RemoveFile:             removeFile,
		SaveStatus:             saveStatus,
	}
	err := executorFactory(deps)()

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestDisableShouldSucceedWhenUnsetCommitTemplateFails(t *testing.T) {
	expectedErr := errors.New("failed to unset commit template")
	deps := dependencies{
		GitUnsetCommitTemplate: func() error { return expectedErr },
	}

	err := executorFactory(deps)()

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestDisableShouldFailWhenRemoveCommitSectionFails(t *testing.T) {
	expectedErr := errors.New("failed to remove commit section")
	deps := dependencies{
		GitUnsetCommitTemplate: unsetCommitTemplate,
		GitRemoveCommitSection: func() error { return expectedErr },
	}

	err := executorFactory(deps)()

	if err == nil || expectedErr != err {
		t.Errorf("expected: %s, received: %s", expectedErr, err)
		t.Fail()
	}
}

func TestDisableShouldFailWhenLoadConfigFails(t *testing.T) {
	expectedErr := errors.New("failed to load config")
	deps := dependencies{
		GitUnsetCommitTemplate: unsetCommitTemplate,
		GitRemoveCommitSection: removeCommitSection,
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
		GitUnsetCommitTemplate: unsetCommitTemplate,
		GitRemoveCommitSection: removeCommitSection,
		LoadConfig:             loadConfig,
		RemoveFile:             func(string) error { return expectedErr },
	}

	err := executorFactory(deps)()

	if err == nil || expectedErr != err {
		t.Errorf("expected: %s, received: %s", expectedErr, err)
		t.Fail()
	}
}

func TestDisableShouldFailWhenSaveStatusFails(t *testing.T) {
	expectedErr := errors.New("failed to save status")
	deps := dependencies{
		GitUnsetCommitTemplate: unsetCommitTemplate,
		GitRemoveCommitSection: removeCommitSection,
		LoadConfig:             loadConfig,
		RemoveFile:             removeFile,
		SaveStatus:             func(status.State, ...string) error { return expectedErr },
	}

	err := executorFactory(deps)()

	if err == nil || expectedErr != err {
		t.Errorf("expected: %s, received: %s", expectedErr, err)
		t.Fail()
	}
}
