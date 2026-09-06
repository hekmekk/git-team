package validation

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
	return fmt.Errorf("not a valid coauthor: %s", candidateCoauthor)
}
