package enable

import (
	"os"
	"testing"

	"github.com/hekmekk/git-team/core/status"
)

func TestEnableSucceeds(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>"}
	baseDir := "BASEDIR"
	templateFileName := "TEMPLATE_FILE"

	createDir := func(string, os.FileMode) error { return nil }
	writeFile := func(string, []byte, os.FileMode) error { return nil }
	setCommitTemplate := func(string) error { return nil }
	saveStatus := func(status.State, ...string) error { return nil }

	effect := EnableEffect{
		CreateDir:         createDir,
		WriteFile:         writeFile,
		SetCommitTemplate: setCommitTemplate,
		SaveStatus:        saveStatus,
	}
	execEnable := ExecutorFactory(effect)

	cmd := EnableCommand{
		Coauthors:        coAuthors,
		BaseDir:          baseDir,
		TemplateFileName: templateFileName,
	}

	errs := execEnable(cmd)

	if len(errs) > 0 && errs[0] != nil {
		t.Error(errs)
		t.Fail()
	}
}
