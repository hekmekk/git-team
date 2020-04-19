package disable

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"
)

type gitConfigReaderMock struct {
	get func(key string) (string, error)
}

func (mock gitConfigReaderMock) Get(key string) (string, error) {
	return mock.get(key)
}

type gitConfigWriterMock struct {
	unsetAll func(key string) error
}

func (mock gitConfigWriterMock) UnsetAll(key string) error {
	return mock.unsetAll(key)
}

func (mock gitConfigWriterMock) ReplaceAll(key string, value string) error {
	return nil
}

var (
	fileInfo        os.FileInfo
	statFile        = func(string) (os.FileInfo, error) { return fileInfo, nil }
	removeFile      = func(string) error { return nil }
	persistDisabled = func() error { return nil }
)

func TestDisableSucceeds(t *testing.T) {
	templatePath := "/path/to/template"

	gitConfigReader := &gitConfigReaderMock{
		get: func(key string) (string, error) {
			return templatePath, nil
		},
	}

	gitConfigWriter := &gitConfigWriterMock{
		unsetAll: func(key string) error {
			switch key {
			case "commit.template":
				return nil
			case "core.hooksPath":
				return nil
			default:
				return fmt.Errorf("wrong key: %s", key)
			}
		},
	}

	deps := Dependencies{
		GitConfigReader: gitConfigReader,
		GitConfigWriter: gitConfigWriter,
		StatFile:        statFile,
		RemoveFile: func(path string) error {
			if path != templatePath {
				t.Errorf("trying to delete the wrong commit template file at path: %s", path)
				t.Fail()
			}
			return nil
		},
		PersistDisabled: persistDisabled,
	}

	expectedEvent := Succeeded{}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestDisableShouldSucceedWhenUnsetHooksPathFailsBecauseTheOptionDoesntExist(t *testing.T) {
	gitConfigReader := &gitConfigReaderMock{
		get: func(key string) (string, error) {
			return "/path/to/template", nil
		},
	}

	expectedErr := errors.New("exit status 5")
	gitConfigWriter := &gitConfigWriterMock{
		unsetAll: func(key string) error {
			switch key {
			case "commit.template":
				return nil
			case "core.hooksPath":
				return expectedErr
			default:
				return fmt.Errorf("wrong key: %s", key)
			}
		},
	}

	deps := Dependencies{
		GitConfigReader: gitConfigReader,
		GitConfigWriter: gitConfigWriter,
		StatFile:        statFile,
		RemoveFile:      removeFile,
		PersistDisabled: persistDisabled,
	}

	expectedEvent := Succeeded{}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestDisableShouldFailWhenUnsetHooksPathFails(t *testing.T) {
	expectedErr := errors.New("failed to unset hooks path")
	gitConfigWriter := &gitConfigWriterMock{
		unsetAll: func(key string) error {
			switch key {
			case "commit.template":
				return nil
			case "core.hooksPath":
				return expectedErr
			default:
				return fmt.Errorf("wrong key: %s", key)
			}
		},
	}

	deps := Dependencies{
		GitConfigWriter: gitConfigWriter,
	}

	expectedEvent := Failed{Reason: expectedErr}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestDisableShouldSucceedWhenReadingCommitTemplateFailsBecauseItHasBeenUnsetAlready(t *testing.T) {
	expectedErr := errors.New("exit status 1")
	gitConfigReader := &gitConfigReaderMock{
		get: func(key string) (string, error) {
			return "", expectedErr
		},
	}

	gitConfigWriter := &gitConfigWriterMock{
		unsetAll: func(key string) error {
			switch key {
			case "commit.template":
				return nil
			case "core.hooksPath":
				return nil
			default:
				return fmt.Errorf("wrong key: %s", key)
			}
		},
	}

	deps := Dependencies{
		GitConfigReader: gitConfigReader,
		GitConfigWriter: gitConfigWriter,
		StatFile:        statFile,
		RemoveFile:      removeFile,
		PersistDisabled: persistDisabled,
	}

	expectedEvent := Succeeded{}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestDisableShouldFailWhenReadingCommitTemplatePathFails(t *testing.T) {
	expectedErr := errors.New("failed to get commit.template")
	gitConfigReader := &gitConfigReaderMock{
		get: func(key string) (string, error) {
			return "", expectedErr
		},
	}

	gitConfigWriter := &gitConfigWriterMock{
		unsetAll: func(key string) error {
			switch key {
			case "commit.template":
				return nil
			case "core.hooksPath":
				return nil
			default:
				return fmt.Errorf("wrong key: %s", key)
			}
		},
	}

	deps := Dependencies{
		GitConfigReader: gitConfigReader,
		GitConfigWriter: gitConfigWriter,
	}

	expectedEvent := Failed{Reason: expectedErr}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestDisableShouldSucceedWhenUnsetCommitTemplateFailsBecauseItWasUnsetAlready(t *testing.T) {
	gitConfigReader := &gitConfigReaderMock{
		get: func(key string) (string, error) {
			return "/path/to/template", nil
		},
	}

	expectedErr := errors.New("exit status 5")
	gitConfigWriter := &gitConfigWriterMock{
		unsetAll: func(key string) error {
			switch key {
			case "commit.template":
				return expectedErr
			case "core.hooksPath":
				return nil
			default:
				return fmt.Errorf("wrong key: %s", key)
			}
		},
	}

	deps := Dependencies{
		GitConfigReader: gitConfigReader,
		GitConfigWriter: gitConfigWriter,
		StatFile:        statFile,
		RemoveFile:      removeFile,
		PersistDisabled: persistDisabled,
	}

	expectedEvent := Succeeded{}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestDisableShouldFailWhenUnsetCommitTemplateFails(t *testing.T) {
	gitConfigReader := &gitConfigReaderMock{
		get: func(key string) (string, error) {
			return "/path/to/template", nil
		},
	}

	expectedErr := errors.New("failed to unset commit template")
	gitConfigWriter := &gitConfigWriterMock{
		unsetAll: func(key string) error {
			switch key {
			case "commit.template":
				return expectedErr
			case "core.hooksPath":
				return nil
			default:
				return fmt.Errorf("wrong key: %s", key)
			}
		},
	}

	deps := Dependencies{
		GitConfigReader: gitConfigReader,
		GitConfigWriter: gitConfigWriter,
	}

	expectedEvent := Failed{Reason: expectedErr}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestDisableShouldFailWhenRemoveFileFails(t *testing.T) {
	gitConfigReader := &gitConfigReaderMock{
		get: func(key string) (string, error) {
			return "/path/to/template", nil
		},
	}

	gitConfigWriter := &gitConfigWriterMock{
		unsetAll: func(key string) error {
			switch key {
			case "commit.template":
				return nil
			case "core.hooksPath":
				return nil
			default:
				return fmt.Errorf("wrong key: %s", key)
			}
		},
	}

	err := errors.New("failed to remove file")

	deps := Dependencies{
		GitConfigReader: gitConfigReader,
		GitConfigWriter: gitConfigWriter,
		StatFile:        statFile,
		RemoveFile:      func(string) error { return err },
	}

	expectedEvent := Failed{Reason: err}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestDisableShouldSucceedButNotTryToRemoveTheCommitTemplateFileWhenStatFileFails(t *testing.T) {
	gitConfigReader := &gitConfigReaderMock{
		get: func(key string) (string, error) {
			return "/path/to/template", nil
		},
	}

	gitConfigWriter := &gitConfigWriterMock{
		unsetAll: func(key string) error {
			switch key {
			case "commit.template":
				return nil
			case "core.hooksPath":
				return nil
			default:
				return fmt.Errorf("wrong key: %s", key)
			}
		},
	}
	deps := Dependencies{
		GitConfigReader: gitConfigReader,
		GitConfigWriter: gitConfigWriter,
		StatFile:        func(string) (os.FileInfo, error) { return fileInfo, errors.New("failed to stat file") },
		PersistDisabled: persistDisabled,
	}

	expectedEvent := Succeeded{}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}

func TestDisableShouldFailWhenpersistDisabledFails(t *testing.T) {
	gitConfigReader := &gitConfigReaderMock{
		get: func(key string) (string, error) {
			return "/path/to/template", nil
		},
	}

	gitConfigWriter := &gitConfigWriterMock{
		unsetAll: func(key string) error {
			switch key {
			case "commit.template":
				return nil
			case "core.hooksPath":
				return nil
			default:
				return fmt.Errorf("wrong key: %s", key)
			}
		},
	}

	expectedErr := errors.New("failed to save status")
	deps := Dependencies{
		GitConfigReader: gitConfigReader,
		GitConfigWriter: gitConfigWriter,
		StatFile:        statFile,
		RemoveFile:      removeFile,
		PersistDisabled: func() error { return expectedErr },
	}

	expectedEvent := Failed{Reason: expectedErr}

	event := Policy{deps}.Apply()

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected: %s, got: %s", expectedEvent, event)
		t.Fail()
	}
}
