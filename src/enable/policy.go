package enable

import (
	"fmt"
	"os"

	"github.com/hekmekk/git-team/src/core/config"
	"github.com/hekmekk/git-team/src/core/events"
	utils "github.com/hekmekk/git-team/src/enable/utils"
)

// Dependencies the dependencies of the enable Policy module
type Dependencies struct {
	SanityCheckCoauthors          func([]string) []error
	LoadConfig                    func() (config.Config, error)
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

	cfg, err := deps.LoadConfig()
	if err != nil {
		return Failed{Reason: []error{err}}
	}

	// TODO: extract these 3 commit template functions into a method
	if err := deps.CreateTemplateDir(cfg.BaseDir, os.ModePerm); err != nil {
		return Failed{Reason: []error{err}}
	}

	templatePath := fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.TemplateFileName)

	if err := deps.WriteTemplateFile(templatePath, []byte(utils.PrepareForCommitMessage(uniqueCoauthors)), 0644); err != nil {
		return Failed{Reason: []error{err}}
	}
	if err := deps.GitSetCommitTemplate(templatePath); err != nil {
		return Failed{Reason: []error{err}}
	}
	if err := deps.GitSetHooksPath(cfg.GitHooksPath); err != nil {
		return Failed{Reason: []error{err}}
	}
	if err := deps.StateRepositoryPersistEnabled(uniqueCoauthors); err != nil {
		return Failed{Reason: []error{err}}
	}
	return Succeeded{}
}
