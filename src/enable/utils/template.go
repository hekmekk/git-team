package enableutils

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// PrepareForCommitMessage create string from coauthors to write to a commit template
func PrepareForCommitMessage(coauthors []string) string {
	if len(coauthors) == 0 {
		return ""
	}

	sort.Strings(coauthors)

	var buffer bytes.Buffer
	buffer.WriteString("\n\n")
	for _, coauthor := range coauthors {
		buffer.WriteString(toLine(coauthor))
	}
	return strings.TrimRight(buffer.String(), "\n")
}

func toLine(coauthor string) string {
	return fmt.Sprintf("Co-authored-by: %s\n", coauthor)
}
