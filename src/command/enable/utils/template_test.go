package enableutils

import (
	"strings"
	"testing"
	"testing/quick"
)

func TestToLine(t *testing.T) {
	t.SkipNow()
	toLineGen := func(coauthor string) bool {
		if coAuthorLine := toLine(coauthor); strings.HasPrefix(coAuthorLine, "Co-authored-by: ") && strings.HasSuffix(coAuthorLine, "\n") {
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

	expectedCoauthorsString := "\n\nCo-authored-by: Mr. Noujz <noujz@mr.se>"

	coauthorsString := PrepareForCommitMessage(coAuthors)

	if expectedCoauthorsString != coauthorsString {
		t.Errorf("expected: [%s], got: [%s]", expectedCoauthorsString, coauthorsString)
		t.Fail()
	}
}

func TestPrepareForCommitMessageMultipleAuthors(t *testing.T) {
	coAuthors := []string{"B <b@x.y>", "A <a@x.y>", "C <c@x.y>"}

	expectedCoauthorsString := "\n\nCo-authored-by: A <a@x.y>\nCo-authored-by: B <b@x.y>\nCo-authored-by: C <c@x.y>"

	coauthorsString := PrepareForCommitMessage(coAuthors)

	if expectedCoauthorsString != coauthorsString {
		t.Errorf("expected: [%s], got: [%s]", expectedCoauthorsString, coauthorsString)
		t.Fail()
	}
}
