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

func EnableCommand(coauthors []string) {
	cfg, _ := config.Load()

	defer status.Print()

	if len(coauthors) == 0 {
		return
	}

	coauthorsString := PrepareForCommitMessage(coauthors)

	mkdirErr := os.MkdirAll(cfg.BaseDir, os.ModePerm)
	if mkdirErr != nil {
		os.Stderr.WriteString(fmt.Sprintf("error: %s\n", mkdirErr.Error()))
		os.Exit(-1)
	}

	templatePath := fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.TemplateFileName)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		writeTemplateFileErr := ioutil.WriteFile(templatePath, []byte(coauthorsString), 0644)
		if writeTemplateFileErr != nil {
			os.Stderr.WriteString(fmt.Sprintf("error: %s\n", writeTemplateFileErr.Error()))
			os.Exit(-1)
		}
	}()

	go func() {
		defer wg.Done()
		gitConfigErr := git.SetCommitTemplate(templatePath)
		if gitConfigErr != nil {
			os.Stderr.WriteString(fmt.Sprintf("error: %s\n", gitConfigErr.Error()))
			os.Exit(-1)
		}
	}()

	go func() {
		defer wg.Done()
		writeStatusFileErr := status.Save(status.ENABLED, coauthors...)
		if writeStatusFileErr != nil {
			os.Stderr.WriteString(fmt.Sprintf("error: %s\n", writeStatusFileErr.Error()))
			os.Exit(-1)
		}
	}()

	wg.Wait()
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
