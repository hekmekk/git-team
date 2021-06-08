package validation

import (
	"fmt"
	"strings"
)

// SanityCheckCoauthors convenience function to check multiple co-authors and accumulate errors
func SanityCheckCoauthors(coauthors []string) []error {
	var validationErrors []error

	for _, coauthor := range coauthors {
		if err := SanityCheckCoauthor(coauthor); err != nil {
			validationErrors = append(validationErrors, err)
		}
	}

	return validationErrors
}

// SanityCheckCoauthor check if provided co-author candidate seem to be a valid co-author
func SanityCheckCoauthor(candidateCoauthor string) error {
	var hasArrowBrackets = strings.Contains(candidateCoauthor, " <") && strings.HasSuffix(candidateCoauthor, ">")
	var containsAtSign = strings.ContainsRune(candidateCoauthor, '@')

	if hasArrowBrackets && containsAtSign {
		return nil
	}
	return fmt.Errorf(fmt.Sprintf("not a valid coauthor: %s", candidateCoauthor))
}
