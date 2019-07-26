package enable

import (
	"errors"
	"os"
	"testing"

	"github.com/hekmekk/git-team/src/config"
)

func TestEnableSucceeds(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>", "mrs"}

	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(string, []byte, os.FileMode) error { return nil }
	setCommitTemplate := func(string) error { return nil }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	persistEnabled := func([]string) error { return nil }
	cfg := config.Config{TemplateFileName: "TEMPLATE_FILE", BaseDir: "BASE_DIR", StatusFileName: "STATUS_FILE"}
	loadConfig := func() (config.Config, error) { return cfg, nil }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		SetCommitTemplate: setCommitTemplate,
		GitResolveAliases: resolveAliases,
		PersistEnabled:    persistEnabled,
		LoadConfig:        loadConfig,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors: coAuthors,
	}

	errs := execEnable(cmd)

	if len(errs) > 0 {
		t.Error(errs[0])
		t.Fail()
	}
}

func TestEnableFailsDueToSanityCheckErr(t *testing.T) {
	coAuthors := []string{"INVALID COAUTHOR"}

	expectedErr := errors.New("Not a valid coauthor: INVALID COAUTHOR")

	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(string, []byte, os.FileMode) error { return nil }
	setCommitTemplate := func(string) error { return nil }
	resolveAliases := func([]string) ([]string, []error) { return []string{}, []error{} }
	persistEnabled := func([]string) error { return nil }
	cfg := config.Config{TemplateFileName: "TEMPLATE_FILE", BaseDir: "BASE_DIR", StatusFileName: "STATUS_FILE"}
	loadConfig := func() (config.Config, error) { return cfg, nil }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		SetCommitTemplate: setCommitTemplate,
		GitResolveAliases: resolveAliases,
		PersistEnabled:    persistEnabled,
		LoadConfig:        loadConfig,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors: coAuthors,
	}

	errs := execEnable(cmd)
	assertEqualsErr(t, expectedErr, errs)
}

func TestEnableFailsDueToResolveAliasesErr(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>", "mrs"}

	expectedErr := errors.New("failed to resolve alias mrs")

	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(string, []byte, os.FileMode) error { return nil }
	setCommitTemplate := func(string) error { return nil }
	resolveAliases := func([]string) ([]string, []error) { return []string{}, []error{expectedErr} }
	persistEnabled := func([]string) error { return nil }
	loadConfig := func() (config.Config, error) { return config.Config{}, expectedErr }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		SetCommitTemplate: setCommitTemplate,
		GitResolveAliases: resolveAliases,
		PersistEnabled:    persistEnabled,
		LoadConfig:        loadConfig,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors: coAuthors,
	}

	errs := execEnable(cmd)
	assertEqualsErr(t, expectedErr, errs)
}

func TestEnableFailsDueToLoadConfigErr(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("failed to load config")

	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(string, []byte, os.FileMode) error { return nil }
	setCommitTemplate := func(string) error { return nil }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	persistEnabled := func([]string) error { return nil }
	loadConfig := func() (config.Config, error) { return config.Config{}, expectedErr }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		SetCommitTemplate: setCommitTemplate,
		GitResolveAliases: resolveAliases,
		PersistEnabled:    persistEnabled,
		LoadConfig:        loadConfig,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors: coAuthors,
	}

	errs := execEnable(cmd)
	assertEqualsErr(t, expectedErr, errs)
}

func TestEnableFailsDueToCreateDirErr(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to create Dir")

	createDir := func(string, os.FileMode) error { return expectedErr }
	writeFile := func(string, []byte, os.FileMode) error { return nil }
	setCommitTemplate := func(string) error { return nil }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	persistEnabled := func([]string) error { return nil }
	cfg := config.Config{TemplateFileName: "TEMPLATE_FILE", BaseDir: "BASE_DIR", StatusFileName: "STATUS_FILE"}
	loadConfig := func() (config.Config, error) { return cfg, nil }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		SetCommitTemplate: setCommitTemplate,
		GitResolveAliases: resolveAliases,
		PersistEnabled:    persistEnabled,
		LoadConfig:        loadConfig,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors: coAuthors,
	}

	errs := execEnable(cmd)
	assertEqualsErr(t, expectedErr, errs)
}

func TestEnableFailsDueToWriteFileErr(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to write file")

	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(string, []byte, os.FileMode) error { return expectedErr }
	setCommitTemplate := func(string) error { return nil }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	persistEnabled := func([]string) error { return nil }
	cfg := config.Config{TemplateFileName: "TEMPLATE_FILE", BaseDir: "BASE_DIR", StatusFileName: "STATUS_FILE"}
	loadConfig := func() (config.Config, error) { return cfg, nil }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		SetCommitTemplate: setCommitTemplate,
		GitResolveAliases: resolveAliases,
		PersistEnabled:    persistEnabled,
		LoadConfig:        loadConfig,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors: coAuthors,
	}

	errs := execEnable(cmd)
	assertEqualsErr(t, expectedErr, errs)
}

func TestEnableFailsDueToSetCommitTemplateErr(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to set commit template")

	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(string, []byte, os.FileMode) error { return nil }
	setCommitTemplate := func(string) error { return expectedErr }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	persistEnabled := func([]string) error { return nil }
	cfg := config.Config{TemplateFileName: "TEMPLATE_FILE", BaseDir: "BASE_DIR", StatusFileName: "STATUS_FILE"}
	loadConfig := func() (config.Config, error) { return cfg, nil }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		SetCommitTemplate: setCommitTemplate,
		GitResolveAliases: resolveAliases,
		PersistEnabled:    persistEnabled,
		LoadConfig:        loadConfig,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors: coAuthors,
	}

	errs := execEnable(cmd)
	assertEqualsErr(t, expectedErr, errs)
}

func TestEnableFailsDueToSaveStatusErr(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to set status")

	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(string, []byte, os.FileMode) error { return nil }
	setCommitTemplate := func(string) error { return nil }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	persistEnabled := func([]string) error { return expectedErr }
	cfg := config.Config{TemplateFileName: "TEMPLATE_FILE", BaseDir: "BASE_DIR", StatusFileName: "STATUS_FILE"}
	loadConfig := func() (config.Config, error) { return cfg, nil }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		SetCommitTemplate: setCommitTemplate,
		GitResolveAliases: resolveAliases,
		PersistEnabled:    persistEnabled,
		LoadConfig:        loadConfig,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors: coAuthors,
	}

	errs := execEnable(cmd)
	assertEqualsErr(t, expectedErr, errs)
}

func assertEqualsErr(t *testing.T, expectedErr error, errs []error) {
	if len(errs) != 1 {
		t.Errorf("got unexpected errs: %s", errs)
		t.Fail()
		return
	}
	if errs[0].Error() != expectedErr.Error() {
		t.Errorf("expexted: %s, got: %s", expectedErr.Error(), errs[0].Error())
		t.Fail()
		return
	}
}
