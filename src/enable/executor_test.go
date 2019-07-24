package enable

import (
	"errors"
	"os"
	"testing"

	"github.com/hekmekk/git-team/src/config"
)

func TestEnableSucceeds(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>"}

	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(string, []byte, os.FileMode) error { return nil }
	setCommitTemplate := func(string) error { return nil }
	persistEnabled := func([]string) error { return nil }
	cfg := config.Config{TemplateFileName: "TEMPLATE_FILE", BaseDir: "BASE_DIR", StatusFileName: "STATUS_FILE"}
	loadConfig := func() (config.Config, error) { return cfg, nil }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		SetCommitTemplate: setCommitTemplate,
		PersistEnabled:    persistEnabled,
		LoadConfig:        loadConfig,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors: coAuthors,
	}

	err := execEnable(cmd)

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestEnableFailsDueToLoadConfigErr(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("failed to load config")

	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(string, []byte, os.FileMode) error { return nil }
	setCommitTemplate := func(string) error { return nil }
	persistEnabled := func([]string) error { return nil }
	loadConfig := func() (config.Config, error) { return config.Config{}, expectedErr }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		SetCommitTemplate: setCommitTemplate,
		PersistEnabled:    persistEnabled,
		LoadConfig:        loadConfig,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors: coAuthors,
	}

	err := execEnable(cmd)

	if err != expectedErr {
		t.Error(err)
		t.Fail()
	}
}

func TestEnableFailsDueToCreateDirErr(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to create Dir")

	createDir := func(string, os.FileMode) error { return expectedErr }
	writeFile := func(string, []byte, os.FileMode) error { return nil }
	setCommitTemplate := func(string) error { return nil }
	persistEnabled := func([]string) error { return nil }
	cfg := config.Config{TemplateFileName: "TEMPLATE_FILE", BaseDir: "BASE_DIR", StatusFileName: "STATUS_FILE"}
	loadConfig := func() (config.Config, error) { return cfg, nil }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		SetCommitTemplate: setCommitTemplate,
		PersistEnabled:    persistEnabled,
		LoadConfig:        loadConfig,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors: coAuthors,
	}

	err := execEnable(cmd)

	if err != expectedErr {
		t.Error(err)
		t.Fail()
	}
}

func TestEnableFailsDueToWriteFileErr(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to write file")

	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(string, []byte, os.FileMode) error { return expectedErr }
	setCommitTemplate := func(string) error { return nil }
	persistEnabled := func([]string) error { return nil }
	cfg := config.Config{TemplateFileName: "TEMPLATE_FILE", BaseDir: "BASE_DIR", StatusFileName: "STATUS_FILE"}
	loadConfig := func() (config.Config, error) { return cfg, nil }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		SetCommitTemplate: setCommitTemplate,
		PersistEnabled:    persistEnabled,
		LoadConfig:        loadConfig,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors: coAuthors,
	}

	err := execEnable(cmd)

	if err != expectedErr {
		t.Error(err)
		t.Fail()
	}
}

func TestEnableFailsDueToSetCommitTemplateErr(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to set commit template")

	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(string, []byte, os.FileMode) error { return nil }
	setCommitTemplate := func(string) error { return expectedErr }
	persistEnabled := func([]string) error { return nil }
	cfg := config.Config{TemplateFileName: "TEMPLATE_FILE", BaseDir: "BASE_DIR", StatusFileName: "STATUS_FILE"}
	loadConfig := func() (config.Config, error) { return cfg, nil }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		SetCommitTemplate: setCommitTemplate,
		PersistEnabled:    persistEnabled,
		LoadConfig:        loadConfig,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors: coAuthors,
	}

	err := execEnable(cmd)

	if err != expectedErr {
		t.Error(err)
		t.Fail()
	}
}

func TestEnableFailsDueToSaveStatusErr(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to set status")

	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(string, []byte, os.FileMode) error { return nil }
	setCommitTemplate := func(string) error { return nil }
	persistEnabled := func([]string) error { return expectedErr }
	cfg := config.Config{TemplateFileName: "TEMPLATE_FILE", BaseDir: "BASE_DIR", StatusFileName: "STATUS_FILE"}
	loadConfig := func() (config.Config, error) { return cfg, nil }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		SetCommitTemplate: setCommitTemplate,
		PersistEnabled:    persistEnabled,
		LoadConfig:        loadConfig,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors: coAuthors,
	}

	err := execEnable(cmd)

	if err != expectedErr {
		t.Error(err)
		t.Fail()
	}
}
