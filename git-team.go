package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/hekmekk/git-team/core/git"
	"github.com/hekmekk/git-team/core/handler"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	version = "v1.0.1"
	author  = "Rea Sand <hekmek@posteo.de>"
)

var (
	handleAdd = handler.RunAddCommand(git.AddAlias)
)

func main() {
	app := kingpin.New("git-team", "Command line interface for creating git commit templates provisioned with one or more co-authors. Please note that \"git commit -m\" is not affected by commit templates.")

	app.HelpFlag.Short('h')
	app.Version(version)
	app.Author(author)

	enable := app.Command("enable", "Provisions a git-commit template with the provided co-authors. A co-author must either be an alias or of the shape \"Name <email>\"").Default()
	coauthors := enable.Arg("coauthors", "Git co-authors").Strings()

	disable := app.Command("disable", "Use default template")
	status := app.Command("status", "Print the current status")

	add := app.Command("add", "Add an alias")
	addAlias := add.Arg("alias", "The alias to be added").Required().String()
	addCoauthor := add.Arg("coauthor", "The co-author").Required().String()

	rm := app.Command("rm", "Remove an alias")
	rmAlias := rm.Arg("alias", "The alias to be removed").Required().String()

	list := app.Command("list", "List currently available aliases")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case enable.FullCommand():
		validCoAuthors, err := validateUserInput(coauthors)
		if len(err) > 0 && err[0] != nil {
			os.Stderr.WriteString(fmt.Sprintf("error: %s\n", foldErrors(err)))
			os.Exit(-1)
		}
		handler.EnableCommand(validCoAuthors)
	case disable.FullCommand():
		handler.DisableCommand()
	case status.FullCommand():
		handler.Status()
	case add.FullCommand():
		checkErr := sanityCheckCoauthor(*addCoauthor)
		if checkErr != nil {
			fmt.Println(checkErr)
			os.Exit(-1)
		}
		aliasAdded, err := handleAdd(*addAlias, *addCoauthor)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
		fmt.Println(color.GreenString(fmt.Sprintf("Alias '%s' -> '%s' has been added.", aliasAdded.Alias, aliasAdded.CoAuthor)))
		os.Exit(0)
	case rm.FullCommand():
		handler.RemoveCommand(rmAlias)
	case list.FullCommand():
		handler.ListCommand()
	}
}

func foldErrors(validationErrors []error) error {
	var buffer bytes.Buffer
	for _, err := range validationErrors {
		buffer.WriteString(err.Error())
		buffer.WriteString("; ")
	}
	return errors.New(strings.TrimRight(buffer.String(), "; "))
}

func validateUserInput(coauthors *[]string) ([]string, []error) {
	var userInputErrors []error

	normalizedCoAuthors, resolveErrors := normalize(*coauthors)

	if resolveErrors != nil {
		userInputErrors = append(userInputErrors, resolveErrors...)
	}

	validCoauthors, validationErrors := coAuthorValidation(normalizedCoAuthors)

	if validationErrors != nil {
		userInputErrors = append(userInputErrors, validationErrors...)
	}

	if len(userInputErrors) > 0 {
		return nil, userInputErrors
	}

	return validCoauthors, nil
}

func normalize(coauthors []string) ([]string, []error) {
	var normalizedCoAuthors []string
	var resolveErrors []error

	for _, maybeAlias := range coauthors {
		if strings.ContainsRune(maybeAlias, ' ') {
			normalizedCoAuthors = append(normalizedCoAuthors, maybeAlias)
		} else {
			var resolvedCoauthor, err = git.ResolveAlias(maybeAlias)
			if err != nil {
				resolveErrors = append(resolveErrors, err)
			} else {
				normalizedCoAuthors = append(normalizedCoAuthors, resolvedCoauthor)
			}
		}
	}

	if len(resolveErrors) > 0 {
		return normalizedCoAuthors, resolveErrors
	}

	return normalizedCoAuthors, nil
}

func coAuthorValidation(coauthors []string) ([]string, []error) {
	var validCoauthors []string
	var validationErrors []error

	for _, coauthor := range coauthors {
		if err := sanityCheckCoauthor(coauthor); err != nil {
			validationErrors = append(validationErrors, err)
		} else {
			validCoauthors = append(validCoauthors, coauthor)
		}
	}

	if len(validationErrors) > 0 {
		return coauthors, validationErrors
	}

	return coauthors, nil
}

func sanityCheckCoauthor(candidateCoauthor string) error {
	var hasArrowBrackets = strings.Contains(candidateCoauthor, " <") && strings.HasSuffix(candidateCoauthor, ">")
	var containsAtSign = strings.ContainsRune(candidateCoauthor, '@')

	if hasArrowBrackets && containsAtSign {
		return nil
	}
	return fmt.Errorf(fmt.Sprintf("Not a valid coauthor: %s", candidateCoauthor))
}

/*
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
*/

/*
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
*/
