package list

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hekmekk/git-team/src/core/assignment"
	"github.com/hekmekk/git-team/src/core/events"
	giterror "github.com/hekmekk/git-team/src/shared/gitconfig/error"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/interface"
	gitconfigscope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
)

// ListRequest how to show the available assignments
type ListRequest struct {
	OnlyAlias *bool
}

// Dependencies the dependencies of the list Policy module
type Dependencies struct {
	GitConfigReader gitconfig.Reader
}

// Policy the policy to apply
type Policy struct {
	Deps Dependencies
	Req  ListRequest
}

// Apply show the available assignments
func (policy Policy) Apply() events.Event {
	deps := policy.Deps

	aliasCoauthorMap, err := deps.GitConfigReader.GetRegexp(gitconfigscope.Global, "team.alias")
	if err != nil && !errors.Is(err, giterror.ErrSectionOrKeyIsInvalid) {
		return RetrievalFailed{Reason: fmt.Errorf("failed to retrieve assignments: %s", err)}
	}

	assignments := []assignment.Assignment{}

	for rawAlias, coauthor := range aliasCoauthorMap {
		alias := strings.TrimPrefix(rawAlias, "team.alias.")
		assignments = append(assignments, assignment.Assignment{Alias: alias, Coauthor: coauthor})
	}

	return RetrievalSucceeded{Assignments: assignments, OnlyAlias: *policy.Req.OnlyAlias}
}
