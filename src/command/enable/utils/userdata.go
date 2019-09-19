package enableutils

import (
	"strings"
)

// Partition expects input to be both coauthors and aliases to separate the two by a very simple heuristic D:
func Partition(userProvidedData []string) ([]string, []string) {
	var coauthorCandidates []string
	var aliasCandidates []string

	for _, datum := range userProvidedData {
		if strings.ContainsRune(datum, ' ') {
			coauthorCandidates = append(coauthorCandidates, datum)
		} else {
			aliasCandidates = append(aliasCandidates, datum)
		}
	}

	return coauthorCandidates, aliasCandidates
}
