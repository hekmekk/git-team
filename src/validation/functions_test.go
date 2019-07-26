package validation

import (
	"testing"
)

var (
	validCoauthors      = []string{"Mr. Noujz <noujz@mr.se>", "Foo <foo@bar.baz>"}
	invalidCoauthors    = []string{"INVALID", "Foo Bar", "A B <a@b.com", "= <>", "foo", "<bar@baz.foo>"} // TODO: Make this more exhaustive...
	bothValidAndInvalid = []string{"Mrs. Noujz <foo@mrs.se>", "foo", "bar", "INVALID"}
)

func TestSanityCheckCoAuthorsValidAuthors(t *testing.T) {
	for _, validCoauthor := range validCoauthors {
		if validationErr := SanityCheckCoauthor(validCoauthor); validationErr != nil {
			t.Errorf("Failed for %s", validCoauthor)
			t.Fail()
		}
	}
}

func TestSanityCheckCoAuthorsInValidAuthors(t *testing.T) {
	for _, invalidCoauthor := range invalidCoauthors {
		if validationErr := SanityCheckCoauthor(invalidCoauthor); validationErr == nil {
			t.Errorf("Failed for %s", invalidCoauthor)
			t.Fail()
		}
	}
}

func TestSanityCheckCoAuthorsShouldReportAllErrors(t *testing.T) {
	errs := SanityCheckCoauthors(bothValidAndInvalid)

	if len(errs) != 3 {
		t.Errorf("expected 2 errors, got: %s", errs)
		t.Fail()
	}
}
