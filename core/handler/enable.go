package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/hekmekk/git-team/core/config"
	"github.com/hekmekk/git-team/core/git"
	"github.com/hekmekk/git-team/core/status"
)

// should config loading happen here?! or just cut entirely different?!
// config -> load, io -> createDir, io -> writeFile, git -> setCommitTemplate, status -> Save, status -> Print

// Command Pattern
// TODO: functions are not serializable... is this a problem here?
type EnableCommand struct {
	coauthors        []string
	baseDir          string
	templateFileName string
}

type EnableEffects struct {
	createDir         func(path string, perm os.FileMode) error
	writeFile         func(path string, data []byte, mode int) error
	setCommitTemplate func(path string) error
	saveStatus        func(state status.State, coauthors ...string) error
}

// func X(effects EnableEffects) func(cmd EnableCommand) []error {
// }

func Enable(coauthors []string, cfg config.Config) []error {
	errs := make([]error, 0)
	if len(coauthors) == 0 {
		return errs
	}

	coauthorsString := PrepareForCommitMessage(coauthors)

	mkdirErr := os.MkdirAll(cfg.BaseDir, os.ModePerm)
	if mkdirErr != nil {
		return append(errs, mkdirErr)
	}

	templatePath := fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.TemplateFileName)

	mutex := &sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		writeTemplateFileErr := ioutil.WriteFile(templatePath, []byte(coauthorsString), 0644)
		if writeTemplateFileErr != nil {
			mutex.Lock()
			errs = append(errs, writeTemplateFileErr)
			mutex.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		gitConfigErr := git.SetCommitTemplate(templatePath)
		if gitConfigErr != nil {
			mutex.Lock()
			errs = append(errs, gitConfigErr)
			mutex.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		writeStatusFileErr := status.Save(status.ENABLED, coauthors...)
		if writeStatusFileErr != nil {
			mutex.Lock()
			errs = append(errs, writeStatusFileErr)
			mutex.Unlock()
		}
	}()

	wg.Wait()
	return errs
}

func ToLine(coauthor string) string {
	return fmt.Sprintf("Co-authored-by: %s\n", coauthor)
}

/*
import (
	"errors"
	"strings"
	"testing"
	"testing/quick"
)

*/

/*
func TestToLine(t *testing.T) {
	t.SkipNow()
	toLineGen := func(coauthor string) bool {
		if coAuthorLine := ToLine(coauthor); strings.HasPrefix(coAuthorLine, "Co-authored-by: ") && strings.HasSuffix(coAuthorLine, "\n") {
			return true
		}
		return false
	}
	if err := quick.Check(toLineGen, nil); err != nil {
		t.Error(err)
	}
}
*/

func PrepareForCommitMessage(coauthors []string) string {
	if len(coauthors) == 0 {
		return ""
	}

	var buffer bytes.Buffer
	buffer.WriteString("\n\n")
	for _, coauthor := range coauthors {
		buffer.WriteString(ToLine(coauthor))
	}
	return strings.TrimRight(buffer.String(), "\n")
}

/*
func TestPrepareForCommitMessageNoAuthors(t *testing.T) {
	coAuthors := []string{}

	coauthorsString := PrepareForCommitMessage(coAuthors)

	if coauthorsString != "" {
		t.Fail()
	}
}

func TestPrepareForCommitMessageOneAuthor(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>"}

	coauthorsString := PrepareForCommitMessage(coAuthors)

	if !strings.HasPrefix(coauthorsString, "\n\n") || strings.HasSuffix(coauthorsString, "\n") {
		t.Fail()
	}
}

func TestPrepareForCommitMessageMultipleAuthors(t *testing.T) {
	coAuthors := []string{"Mr. Noujz <noujz@mr.se>", "Mr. Noujz <noujz@mr.se>", "Mr. Noujz <noujz@mr.se>"}

	coauthorsString := PrepareForCommitMessage(coAuthors)

	if !strings.HasPrefix(coauthorsString, "\n\n") || strings.HasSuffix(coauthorsString, "\n") {
		t.Fail()
	}
}
*/
