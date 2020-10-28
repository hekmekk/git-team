package list

import (
	"strings"

	"github.com/hekmekk/git-team/src/core/assignment"
	"github.com/hekmekk/git-team/src/core/events"
	giterror "github.com/hekmekk/git-team/src/shared/gitconfig/error"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/interface"
	gitconfigscope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
)

// Dependencies the dependencies of the list Policy module
type Dependencies struct {
	GitConfigReader gitconfig.Reader
}

// Policy the policy to apply
type Policy struct {
	Deps Dependencies
}

// Apply show the available assignments
func (policy Policy) Apply() events.Event {
	deps := policy.Deps

	aliasCoauthorMap, err := deps.GitConfigReader.GetRegexp(gitconfigscope.Global, "team.alias")
	if err != nil && err.Error() != giterror.SectionOrKeyIsInvalid {
		return RetrievalFailed{Reason: err}
	}

	assignments := []assignment.Assignment{}

	for rawAlias, coauthor := range aliasCoauthorMap {
		alias := strings.TrimPrefix(rawAlias, "team.alias.")
		assignments = append(assignments, assignment.Assignment{Alias: alias, Coauthor: coauthor})
	}

	return RetrievalSucceeded{Assignments: assignments}
}
