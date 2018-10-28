package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hekmekk/git-team/core/git"
	"os"
	"strings"
)

func ToLine(coauthor string) string {
	return fmt.Sprintf("Co-authored-by: %s\n", coauthor)
}

func PrepareForCommitMessage(coauthors []string) string {
	if len(coauthors) == 0 {
		return ""
	}

	var buffer bytes.Buffer
	buffer.WriteString("\n\n")
	for _, coauthor := range coauthors {
		buffer.WriteString(ToLine(coauthor))
	}
	return strings.TrimRight(buffer.String(), "\n")
}

func ValidateUserInput(coauthors *[]string) ([]string, []error) {
	var userInputErrors []error

	coauthorsWithResolvedAliases, resolveErrors := ReplaceAliasWithActualCoAuthor(*coauthors)

	if resolveErrors != nil {
		userInputErrors = append(userInputErrors, resolveErrors...)
	}

	validCoauthors, validationErrors := git.CoAuthorValidation(coauthorsWithResolvedAliases)

	if validationErrors != nil {
		userInputErrors = append(userInputErrors, validationErrors...)
	}

	if len(userInputErrors) > 0 {
		return nil, userInputErrors
	}

	return validCoauthors, nil
}

func ReplaceAliasWithActualCoAuthor(coauthors []string) ([]string, []error) {
	var coauthorsWithResolvedAliases []string
	var resolveErrors []error

	for _, maybeAlias := range coauthors {
		if strings.ContainsRune(maybeAlias, ' ') {
			coauthorsWithResolvedAliases = append(coauthorsWithResolvedAliases, maybeAlias)
		} else {
			var resolvedCoauthor, err = git.ResolveAlias(maybeAlias)
			if err != nil {
				resolveErrors = append(resolveErrors, err)
			} else {
				coauthorsWithResolvedAliases = append(coauthorsWithResolvedAliases, resolvedCoauthor)
			}
		}
	}

	if len(resolveErrors) > 0 {
		return coauthorsWithResolvedAliases, resolveErrors
	}

	return coauthorsWithResolvedAliases, nil
}

func ToStderrAndExit(err ...error) {
	if len(err) > 0 && err[0] != nil {
		os.Stderr.WriteString(fmt.Sprintf("error: %s\n", FoldErrors(err)))
		os.Exit(-1)
	}
}

func FoldErrors(validationErrors []error) error {
	var buffer bytes.Buffer
	for _, err := range validationErrors {
		buffer.WriteString(err.Error())
		buffer.WriteString("; ")
	}
	return errors.New(strings.TrimRight(buffer.String(), "; "))
}
