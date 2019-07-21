package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	addExecutor "github.com/hekmekk/git-team/src/add"
	"github.com/hekmekk/git-team/src/config"
	execDisable "github.com/hekmekk/git-team/src/disable"
	enableExecutor "github.com/hekmekk/git-team/src/enable"
	git "github.com/hekmekk/git-team/src/gitconfig"
	removeExecutor "github.com/hekmekk/git-team/src/remove"
	statusApi "github.com/hekmekk/git-team/src/status"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	version = "v1.1.0"
	author  = "Rea Sand <hekmek@posteo.de>"
)

var (
	addDeps = addExecutor.Dependencies{
		AddGitAlias: git.AddAlias,
	}
	enableDeps = enableExecutor.Dependencies{
		CreateDir:         os.MkdirAll,      // TODO: CreateTemplateDir
		WriteFile:         ioutil.WriteFile, // TODO: WriteTemplateFile
		SetCommitTemplate: git.SetCommitTemplate,
		PersistEnabled:    statusApi.PersistEnabled,
	}
	rmDeps = removeExecutor.Dependencies{
		GitResolveAlias: git.ResolveAlias,
		GitRemoveAlias:  git.RemoveAlias,
	}
	execAdd    = addExecutor.ExecutorFactory(addDeps)
	execEnable = enableExecutor.ExecutorFactory(enableDeps)
	execRemove = removeExecutor.ExecutorFactory(rmDeps)
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
	list.Alias("ls")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case enable.FullCommand():
		// TODO: should be in some validation package
		validCoAuthors, validationErrs := validateUserInput(coauthors)
		printErrAndExit(validationErrs...)

		// TODO: resolve inconsitencies (where to load config?)
		cfg, configErr := config.Load()
		printErrAndExit(configErr)

		cmd := enableExecutor.Command{
			Coauthors:        validCoAuthors,
			BaseDir:          cfg.BaseDir,
			TemplateFileName: cfg.TemplateFileName,
		}
		enableErr := execEnable(cmd)
		printErrAndExit(enableErr)

		status, err := statusApi.Fetch()
		printErrAndExit(err)

		fmt.Println(status.ToString())
		os.Exit(0)
	case disable.FullCommand():
		err := execDisable.Exec()
		printErrAndExit(err)

		status, err := statusApi.Fetch()
		printErrAndExit(err)

		fmt.Println(status.ToString())
		os.Exit(0)
	case status.FullCommand():
		// TODO: should we return the "effect" PrintStatus?
		status, err := statusApi.Fetch()
		printErrAndExit(err)

		fmt.Println(status.ToString())
		os.Exit(0)
	case add.FullCommand():
		checkErr := sanityCheckCoauthor(*addCoauthor)
		printErrAndExit(checkErr)

		cmd := addExecutor.Command{
			Alias:    *addAlias,
			Coauthor: *addCoauthor,
		}
		addErr := execAdd(cmd)
		printErrAndExit(addErr)

		color.Green(fmt.Sprintf("Alias '%s' -> '%s' has been added.", *addAlias, *addCoauthor))
		os.Exit(0)
	case rm.FullCommand():
		cmd := removeExecutor.Command{
			Alias: *rmAlias,
		}

		rmErr := execRemove(cmd)
		printErrAndExit(rmErr)

		color.Red(fmt.Sprintf("Alias '%s' has been removed.", cmd.Alias))
		os.Exit(0)
	case list.FullCommand():
		assignments := git.GetAddedAliases()

		blackBold := color.New(color.FgBlack).Add(color.Bold)
		blackBold.Println("Aliases:")
		blackBold.Println("--------")

		for alias, coauthor := range assignments {
			color.Magenta(fmt.Sprintf("'%s' -> '%s'", alias, coauthor))
		}
		os.Exit(0)
	}
}

func printErrAndExit(validationErrs ...error) {
	if len(validationErrs) > 0 && validationErrs[0] != nil {
		os.Stderr.WriteString(fmt.Sprintf("error: %s\n", foldErrors(validationErrs)))
		os.Exit(-1)
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
