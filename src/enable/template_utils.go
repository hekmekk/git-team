package enable

import (
	"bytes"
	"fmt"
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
