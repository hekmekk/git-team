package handler

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/hekmekk/git-team/core/git"
	"github.com/hekmekk/git-team/core/status"
)

type EnableCommand struct {
	Coauthors        []string
	BaseDir          string
	TemplateFileName string
}

type EnableEffect struct {
	CreateDir         func(path string, perm os.FileMode) error
	WriteFile         func(path string, data []byte, mode os.FileMode) error
	SetCommitTemplate func(path string) error
	SaveStatus        func(state status.State, coauthors ...string) error
}

// -> ExecutorFactory -> execEnable(cmd)
func EnableFactory(effect EnableEffect) func(cmd EnableCommand) []error {
	return func(cmd EnableCommand) []error {
		errs := make([]error, 0)
		if len(cmd.Coauthors) == 0 {
			return errs
		}

		coauthorsString := PrepareForCommitMessage(cmd.Coauthors)

		mkdirErr := effect.CreateDir(cmd.BaseDir, os.ModePerm)
		if mkdirErr != nil {
			return append(errs, mkdirErr)
		}

		templatePath := fmt.Sprintf("%s/%s", cmd.BaseDir, cmd.TemplateFileName)

		mutex := &sync.Mutex{}
		var wg sync.WaitGroup
		wg.Add(3)

		go func() {
			defer wg.Done()
			writeTemplateFileErr := effect.WriteFile(templatePath, []byte(coauthorsString), 0644)
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
			writeStatusFileErr := status.Save(status.ENABLED, cmd.Coauthors...)
			if writeStatusFileErr != nil {
				mutex.Lock()
				errs = append(errs, writeStatusFileErr)
				mutex.Unlock()
			}
		}()

		wg.Wait()
		return errs
	}
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
