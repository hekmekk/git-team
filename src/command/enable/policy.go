package enable

import (
	"os"
	"path"

	utils "github.com/hekmekk/git-team/src/command/enable/utils"
	"github.com/hekmekk/git-team/src/core/config"
	"github.com/hekmekk/git-team/src/core/events"
)

// Dependencies the dependencies of the enable Policy module
type Dependencies struct {
	SanityCheckCoauthors          func([]string) []error
	LoadConfig                    func() config.Config
	CreateTemplateDir             func(path string, perm os.FileMode) error
	WriteTemplateFile             func(path string, data []byte, mode os.FileMode) error
	GitSetCommitTemplate          func(path string) error
	GitSetHooksPath               func(path string) error
	GitResolveAliases             func(aliases []string) ([]string, []error)
	StateRepositoryPersistEnabled func(coauthors []string) error
}

// Request the coauthors with which to enable git-team
type Request struct {
	AliasesAndCoauthors *[]string
}

// Policy add a <Coauthor> under "team.alias.<Alias>"
type Policy struct {
	Deps Dependencies
	Req  Request
}

// Apply enable git-team with the provided co-authors
func (policy Policy) Apply() events.Event {
	deps := policy.Deps
	req := policy.Req

	aliasesAndCoauthors := append(*req.AliasesAndCoauthors) // should be == *req.AliasesAndCoauthors

	if len(aliasesAndCoauthors) == 0 {
		return Aborted{}
	}

	coauthorCandidates, aliases := utils.Partition(aliasesAndCoauthors)

	sanityCheckErrs := deps.SanityCheckCoauthors(coauthorCandidates)
	if len(sanityCheckErrs) > 0 {
		return Failed{Reason: sanityCheckErrs}
	}

	resolvedAliases, resolveErrs := deps.GitResolveAliases(aliases)
	if len(resolveErrs) > 0 {
		return Failed{Reason: resolveErrs}
	}

	coauthors := append(coauthorCandidates, resolvedAliases...)

	var uniqueCoauthors []string
	temp := make(map[string]bool)
	for _, coauthor := range coauthors {
		temp[coauthor] = true
	}

	for coauthor := range temp {
		uniqueCoauthors = append(uniqueCoauthors, coauthor)
	}

	cfg := deps.LoadConfig()
	commitTemplatePath := cfg.GitTeamCommitTemplatePath
	commitTemplateDir := path.Dir(commitTemplatePath)

	// TODO: extract these 3 commit template functions into a method
	if err := deps.CreateTemplateDir(commitTemplateDir, os.ModePerm); err != nil {
		return Failed{Reason: []error{err}}
	}

	if err := deps.WriteTemplateFile(commitTemplatePath, []byte(utils.PrepareForCommitMessage(uniqueCoauthors)), 0644); err != nil {
		return Failed{Reason: []error{err}}
	}
	if err := deps.GitSetCommitTemplate(commitTemplatePath); err != nil {
		return Failed{Reason: []error{err}}
	}
	if err := deps.GitSetHooksPath(cfg.GitTeamHooksPath); err != nil {
		return Failed{Reason: []error{err}}
	}
	if err := deps.StateRepositoryPersistEnabled(uniqueCoauthors); err != nil {
		return Failed{Reason: []error{err}}
	}
	return Succeeded{}
}
