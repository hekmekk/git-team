package handler

import (
	"errors"
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

func TestFoldErrors(t *testing.T) {
	err_prefix := errors.New("_prefix_")
	err_suffix := errors.New("_suffix_")

	// Note: It is more than twice as slow with this predicate approach... Maybe revert to direct inline calls
	isNotNil := func(err error) bool { return err != nil }
	hasProperPrefix := func(err error) bool { return strings.HasPrefix(err.Error(), err_prefix.Error()) }
	hasProperSuffix := func(err error) bool { return strings.HasSuffix(err.Error(), err_suffix.Error()) }

	errorsGen := func(msg string) bool {
		generated_err := errors.New(msg)
		errs := []error{err_prefix, generated_err, err_suffix}

		if folded_err := FoldErrors(errs); isNotNil(folded_err) && hasProperPrefix(folded_err) && hasProperSuffix(folded_err) {
			return true
		}
		return false
	}

	if err := quick.Check(errorsGen, nil); err != nil {
		t.Error(err)
	}
}
