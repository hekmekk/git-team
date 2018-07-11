package git

import (
	"testing"
)

var (
	validCoauthors   = []string{"Mr. Noujz <noujz@mr.se>", "Foo <foo@bar.baz>"}            // TODO: Make this more exhaustive...
	invalidCoauthors = []string{"Foo Bar", "A B <a@b.com", "= <>", "foo", "<bar@baz.foo>"} // TODO: Make this more exhaustive...
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
