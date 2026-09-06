package commandadapter

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
	"errors"
	"reflect"
	"testing"
)

func TestShouldReturnNoErrors(t *testing.T) {
	aliases := []string{"mrs", "mr"}
	expectedCoauthors := []string{"Mrs. Noujz <noujz@mrs.se>", "Mr. Noujz <noujz@mr.se>"}

	coauthorMapping := make(map[string]string)
	for index, alias := range aliases {
		coauthorMapping[alias] = expectedCoauthors[index]
	}

	resolveAlias := func(alias string) (string, error) { return coauthorMapping[alias], nil }

	coauthors, errs := resolveAliases(resolveAlias)(aliases)

	if len(errs) > 0 {
		t.Errorf("unexpected errors: %s", errs)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedCoauthors, coauthors) {
		t.Errorf("expected: %s, got: %s", expectedCoauthors, coauthors)
		t.Fail()
	}
}

type resolveresult struct {
	coauthor string
	err      error
}

func TestShouldAccumulateErrs(t *testing.T) {
	aliases := []string{"mrs", "mr"}
	coauthorMapping := map[string]resolveresult{"mrs": resolveresult{coauthor: "Mrs. Noujz <noujz@mrs.se>", err: nil}, "mr": resolveresult{coauthor: "", err: errors.New("failed to resolve alias mr")}}

	resolveAlias := func(alias string) (string, error) { return coauthorMapping[alias].coauthor, coauthorMapping[alias].err }

	_, errs := resolveAliases(resolveAlias)(aliases)

	if len(errs) != 1 || errs[0].Error() != "failed to resolve alias mr" {
		t.Errorf("unexpected amount of errors: %s", errs)
		t.Fail()
	}
}
