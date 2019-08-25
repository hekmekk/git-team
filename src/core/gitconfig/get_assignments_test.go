package gitconfig

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestShouldReturnAliasAssignments(t *testing.T) {

	mr := "mr"
	mrs := "mrs"
	mrNoujz := "Mr. Noujz <noujz@mr.se>"
	mrsNoujz := "Mrs. Noujz <noujz@mrs.se>"

	expectedMap := make(map[string]string)
	expectedMap[mr] = mrNoujz
	expectedMap[mrs] = mrsNoujz

	lines := make([]string, 0)
	lines = append(lines, fmt.Sprintf("team.alias.%s\n%s\n", mr, mrNoujz))
	lines = append(lines, fmt.Sprintf("team.alias.%s\n%s\n", mrs, mrsNoujz))

	execSucceeds := func(args ...string) ([]string, error) { return lines, nil }

	aliasMap, _ := getAssignments(execSucceeds)

	if !reflect.DeepEqual(aliasMap, expectedMap) {
		t.Errorf("expected: %s, received %s", expectedMap, aliasMap)
		t.Fail()
	}
}

func TestShouldReturnEmptyMapIfEmptyReturnFromGitConfigCommand(t *testing.T) {

	expectedMap := make(map[string]string)

	execFails := func(args ...string) ([]string, error) { return nil, nil }

	aliasMap, _ := getAssignments(execFails)

	if !reflect.DeepEqual(aliasMap, expectedMap) {
		t.Errorf("expected: %s, received %s", expectedMap, aliasMap)
		t.Fail()
	}
}

func TestShouldFailIfGitConfigCommandFails(t *testing.T) {
	expectedErr := errors.New("failed to exec git config command")
	expectedMap := make(map[string]string)

	execFails := func(args ...string) ([]string, error) { return nil, expectedErr }

	aliasMap, err := getAssignments(execFails)

	if !reflect.DeepEqual(aliasMap, expectedMap) {
		t.Errorf("expected: %s, received %s", expectedMap, aliasMap)
		t.Fail()
	}

	if expectedErr != err {
		t.Errorf("expected: %s, received %s", expectedErr, err)
		t.Fail()
	}
}
