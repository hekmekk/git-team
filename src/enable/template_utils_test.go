package enable

import (
	"strings"
	"testing"
	"testing/quick"
)

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
