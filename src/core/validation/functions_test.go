package validation

/*
Copyright (C) 2026 Rea Sand

This file is part of git-team.

This program is free software: you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation, either version 3
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

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
