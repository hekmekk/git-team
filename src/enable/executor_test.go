package enable

import (
	"errors"
	"os"
	"testing"
)

func TestEnableSucceeds(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>"}
	baseDir := "BASEDIR"
	templateFileName := "TEMPLATE_FILE"

	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(string, []byte, os.FileMode) error { return nil }
	setCommitTemplate := func(string) error { return nil }
	persistEnabled := func([]string) error { return nil }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		SetCommitTemplate: setCommitTemplate,
		PersistEnabled:    persistEnabled,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors:        coAuthors,
		BaseDir:          baseDir,
		TemplateFileName: templateFileName,
	}

	err := execEnable(cmd)

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestEnableFailsDueToCreateDirErr(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>"}
	baseDir := "BASEDIR"
	templateFileName := "TEMPLATE_FILE"

	expectedErr := errors.New("Failed to create Dir")

	createDir := func(string, os.FileMode) error { return expectedErr }
	writeFile := func(string, []byte, os.FileMode) error { return nil }
	setCommitTemplate := func(string) error { return nil }
	persistEnabled := func([]string) error { return nil }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		SetCommitTemplate: setCommitTemplate,
		PersistEnabled:    persistEnabled,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors:        coAuthors,
		BaseDir:          baseDir,
		TemplateFileName: templateFileName,
	}

	err := execEnable(cmd)

	if err != expectedErr {
		t.Error(err)
		t.Fail()
	}
}

func TestEnableFailsDueToWriteFileErr(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>"}
	baseDir := "BASEDIR"
	templateFileName := "TEMPLATE_FILE"

	expectedErr := errors.New("Failed to write file")

	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(string, []byte, os.FileMode) error { return expectedErr }
	setCommitTemplate := func(string) error { return nil }
	persistEnabled := func([]string) error { return nil }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		SetCommitTemplate: setCommitTemplate,
		PersistEnabled:    persistEnabled,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors:        coAuthors,
		BaseDir:          baseDir,
		TemplateFileName: templateFileName,
	}

	err := execEnable(cmd)

	if err != expectedErr {
		t.Error(err)
		t.Fail()
	}
}

func TestEnableFailsDueToSetCommitTemplateErr(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>"}
	baseDir := "BASEDIR"
	templateFileName := "TEMPLATE_FILE"

	expectedErr := errors.New("Failed to set commit template")

	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(string, []byte, os.FileMode) error { return nil }
	setCommitTemplate := func(string) error { return expectedErr }
	persistEnabled := func([]string) error { return nil }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		SetCommitTemplate: setCommitTemplate,
		PersistEnabled:    persistEnabled,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors:        coAuthors,
		BaseDir:          baseDir,
		TemplateFileName: templateFileName,
	}

	err := execEnable(cmd)

	if err != expectedErr {
		t.Error(err)
		t.Fail()
	}
}

func TestEnableFailsDueToSaveStatusErr(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>"}
	baseDir := "BASEDIR"
	templateFileName := "TEMPLATE_FILE"

	expectedErr := errors.New("Failed to set status")

	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(string, []byte, os.FileMode) error { return nil }
	setCommitTemplate := func(string) error { return nil }
	persistEnabled := func([]string) error { return expectedErr }

	deps := Dependencies{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		SetCommitTemplate: setCommitTemplate,
		PersistEnabled:    persistEnabled,
	}
	execEnable := ExecutorFactory(deps)

	cmd := Command{
		Coauthors:        coAuthors,
		BaseDir:          baseDir,
		TemplateFileName: templateFileName,
	}

	err := execEnable(cmd)

	if err != expectedErr {
		t.Error(err)
		t.Fail()
	}
}
