package enableutils

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
