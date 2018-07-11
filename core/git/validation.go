package git

import "fmt"
import "strings"
import "errors"

func SanityCheckCoauthor(candidateCoauthor string) error {
	var hasArrowBrackets = strings.Contains(candidateCoauthor, " <") && strings.HasSuffix(candidateCoauthor, ">")
	var containsAtSign = strings.ContainsRune(candidateCoauthor, '@')

	if hasArrowBrackets && containsAtSign {
		return nil
	}
	return errors.New(fmt.Sprintf("Not a valid coauthor: %s", candidateCoauthor))
}

func CoAuthorValidation(coauthors []string) ([]string, []error) {
	var validCoauthors []string
	var validationErrors []error

	for _, coauthor := range coauthors {
		if err := SanityCheckCoauthor(coauthor); err != nil {
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
