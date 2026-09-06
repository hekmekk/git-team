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
	"testing"

	gitconfigscope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
)

func TestShouldReturnTheAssignedCoAuthor(t *testing.T) {
	mr := "mr"
	mrNoujz := "Mr. Noujz <noujz@mr.se>"

	gitconfigGet := func(scope gitconfigscope.Scope, args string) (string, error) {
		if scope != gitconfigscope.Global {
			return "", errors.New("wrong scope")
		}
		return mrNoujz, nil
	}

	coauthor, err := resolveAlias(gitconfigGet)(mr)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if coauthor != mrNoujz {
		t.Errorf("expected: %s, received: %s", mrNoujz, coauthor)
		t.Fail()
	}
}

func TestShouldReturnErrorIfNoAssignmentIsFound(t *testing.T) {
	mr := "mr"

	gitconfigGet := func(scope gitconfigscope.Scope, alias string) (string, error) {
		if scope != gitconfigscope.Global {
			return "", errors.New("wrong scope")
		}
		return "", nil
	}

	coauthor, err := resolveAlias(gitconfigGet)(mr)

	if err == nil {
		t.Error("expected an error")
		t.Fail()
	}

	if coauthor != "" {
		t.Errorf("expected: %s, received: %s", "", coauthor)
		t.Fail()
	}
}

func TestShouldReturnErrorIfResolvingFails(t *testing.T) {
	mr := "mr"
	expectedErr := errors.New("failed to resolve alias team.alias.mr")

	gitconfigGet := func(scope gitconfigscope.Scope, alias string) (string, error) {
		if scope != gitconfigscope.Global {
			return "", errors.New("wrong scope")
		}
		return "", errors.New("git command failed")
	}

	coauthor, err := resolveAlias(gitconfigGet)(mr)

	if err.Error() != expectedErr.Error() {
		t.Errorf("expected: %s, received: %s", expectedErr, err)
		t.Fail()
	}

	if coauthor != "" {
		t.Errorf("expected: %s, received: %s", "", coauthor)
		t.Fail()
	}
}
