package gitconfig

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
