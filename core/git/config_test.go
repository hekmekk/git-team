package git

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestGetAliasSucceeds(t *testing.T) {

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

	aliasMap := getAliasMap(execSucceeds)

	if reflect.DeepEqual(aliasMap, expectedMap) != true {
		t.Errorf("expected: %s, received %s", expectedMap, aliasMap)
		t.Fail()
	}
}

func TestGetAliasFailsDueToEmptyReturnFromGitConfigCommand(t *testing.T) {

	expectedMap := make(map[string]string)

	execFails := func(args ...string) ([]string, error) { return nil, nil }

	aliasMap := getAliasMap(execFails)

	if reflect.DeepEqual(aliasMap, expectedMap) != true {
		t.Errorf("expected: %s, received %s", expectedMap, aliasMap)
		t.Fail()
	}
}

func TestShouldReturnEmptyMapIfGitConfigCommandFails(t *testing.T) {

	expectedMap := make(map[string]string)

	execFails := func(args ...string) ([]string, error) { return nil, errors.New("failed to exec git config command") }

	aliasMap := getAliasMap(execFails)

	if reflect.DeepEqual(aliasMap, expectedMap) != true {
		t.Errorf("expected: %s, received %s", expectedMap, aliasMap)
		t.Fail()
	}
}
