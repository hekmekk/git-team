package enable

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/hekmekk/git-team/src/config"
)

func TestEnableSucceeds(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>", "mrs"}
	expectedPersistEnabledCoauthors := []string{"Mr. Noujz <noujz@mr.se>", "Mrs. Noujz <noujz@mrs.se>"}
	expectedCommitTemplateCoauthors := "\n\nCo-authored-by: Mr. Noujz <noujz@mr.se>\nCo-authored-by: Mrs. Noujz <noujz@mrs.se>"

	setHooksPath := func(string) error { return nil }
	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(_ string, data []byte, _ os.FileMode) error {
		if expectedCommitTemplateCoauthors != string(data) {
			t.Errorf("expected: %s, got: %s", expectedCommitTemplateCoauthors, string(data))
			t.Fail()
		}
		return nil
	}
	setCommitTemplate := func(string) error { return nil }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	persistEnabled := func(coauthors []string) error {
		if !reflect.DeepEqual(expectedPersistEnabledCoauthors, coauthors) {
			t.Errorf("expected: %s, got: %s", expectedPersistEnabledCoauthors, coauthors)
			t.Fail()
		}
		return nil
	}
	cfg := config.Config{TemplateFileName: "TEMPLATE_FILE", BaseDir: "BASE_DIR", StatusFileName: "STATUS_FILE"}
	loadConfig := func() (config.Config, error) { return cfg, nil }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		GitSetHooksPath:   setHooksPath,
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

func TestEnableDropsDuplicateEntries(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>", "mrs", "Mr. Noujz <noujz@mr.se>", "mrs", "Mrs. Noujz <noujz@mrs.se>"}
	expectedPersistEnabledCoauthors := []string{"Mr. Noujz <noujz@mr.se>", "Mrs. Noujz <noujz@mrs.se>"}
	expectedCommitTemplateCoauthors := "\n\nCo-authored-by: Mr. Noujz <noujz@mr.se>\nCo-authored-by: Mrs. Noujz <noujz@mrs.se>"

	setHooksPath := func(string) error { return nil }
	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(_ string, data []byte, _ os.FileMode) error {
		if expectedCommitTemplateCoauthors != string(data) {
			t.Errorf("expected: %s, got: %s", expectedCommitTemplateCoauthors, string(data))
			t.Fail()
		}
		return nil
	}
	setCommitTemplate := func(string) error { return nil }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	persistEnabled := func(coauthors []string) error {
		if !reflect.DeepEqual(expectedPersistEnabledCoauthors, coauthors) {
			t.Errorf("expected: %s, got: %s", expectedPersistEnabledCoauthors, coauthors)
			t.Fail()
		}
		return nil
	}
	cfg := config.Config{TemplateFileName: "TEMPLATE_FILE", BaseDir: "BASE_DIR", StatusFileName: "STATUS_FILE"}
	loadConfig := func() (config.Config, error) { return cfg, nil }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		GitSetHooksPath:   setHooksPath,
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

	setHooksPath := func(string) error { return nil }
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
		GitSetHooksPath:   setHooksPath,
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

	setHooksPath := func(string) error { return nil }
	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(string, []byte, os.FileMode) error { return nil }
	setCommitTemplate := func(string) error { return nil }
	resolveAliases := func([]string) ([]string, []error) { return []string{}, []error{expectedErr} }
	persistEnabled := func([]string) error { return nil }
	loadConfig := func() (config.Config, error) { return config.Config{}, expectedErr }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		GitSetHooksPath:   setHooksPath,
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

	setHooksPath := func(string) error { return nil }
	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(string, []byte, os.FileMode) error { return nil }
	setCommitTemplate := func(string) error { return nil }
	resolveAliases := func([]string) ([]string, []error) { return []string{"Mrs. Noujz <noujz@mrs.se>"}, []error{} }
	persistEnabled := func([]string) error { return nil }
	loadConfig := func() (config.Config, error) { return config.Config{}, expectedErr }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		GitSetHooksPath:   setHooksPath,
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

	setHooksPath := func(string) error { return nil }
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
		GitSetHooksPath:   setHooksPath,
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

	setHooksPath := func(string) error { return nil }
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
		GitSetHooksPath:   setHooksPath,
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
	coAuthorsAndAliases := []string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to set commit template")

	setHooksPath := func(string) error { return nil }
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
		GitSetHooksPath:   setHooksPath,
		SetCommitTemplate: setCommitTemplate,
		GitResolveAliases: resolveAliases,
		PersistEnabled:    persistEnabled,
		LoadConfig:        loadConfig,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors: coAuthorsAndAliases,
	}

	errs := execEnable(cmd)
	assertEqualsErr(t, expectedErr, errs)
}

func TestEnableFailsDueToSetHooksPathErr(t *testing.T) {
	coAuthorsAndAliases := []string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to set hooks path")

	setHooksPath := func(string) error { return expectedErr }
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
		GitSetHooksPath:   setHooksPath,
		SetCommitTemplate: setCommitTemplate,
		GitResolveAliases: resolveAliases,
		PersistEnabled:    persistEnabled,
		LoadConfig:        loadConfig,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors: coAuthorsAndAliases,
	}

	errs := execEnable(cmd)
	assertEqualsErr(t, expectedErr, errs)
}

func TestEnableFailsDueToSaveStatusErr(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>"}

	expectedErr := errors.New("Failed to set status")

	setHooksPath := func(string) error { return nil }
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
		GitSetHooksPath:   setHooksPath,
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
