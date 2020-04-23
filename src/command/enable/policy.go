package enable

import (
	"fmt"
	"os"

	commitsettings "github.com/hekmekk/git-team/src/command/enable/commitsettings/interface"
	utils "github.com/hekmekk/git-team/src/command/enable/utils"
	"github.com/hekmekk/git-team/src/core/events"
	activationscope "github.com/hekmekk/git-team/src/shared/config/entity/activationscope"
	config "github.com/hekmekk/git-team/src/shared/config/interface"
	gitconfig "github.com/hekmekk/git-team/src/shared/gitconfig/interface"
	gitconfigscope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
	state "github.com/hekmekk/git-team/src/shared/state/interface"
)

// Dependencies the dependencies of the enable Policy module
type Dependencies struct {
	SanityCheckCoauthors func([]string) []error
	CommitSettingsReader commitsettings.Reader
	CreateTemplateDir    func(path string, perm os.FileMode) error
	WriteTemplateFile    func(path string, data []byte, mode os.FileMode) error
	GitResolveAliases    func(aliases []string) ([]string, []error)
	ConfigReader         config.Reader
	GitConfigWriter      gitconfig.Writer
	StateWriter          state.Writer
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

	coauthors, errs := applyAdditionalGuards(deps, aliasesAndCoauthors)
	if len(errs) > 0 {
		return Failed{Reason: errs}
	}

	uniqueCoauthors := removeDuplicates(coauthors)

	settings := deps.CommitSettingsReader.Read()

	_, err := deps.ConfigReader.Read()
	if err != nil {
		return Failed{Reason: []error{err}}
	}

	if err := setupTemplate(deps, settings.TemplatesBaseDir, uniqueCoauthors); err != nil {
		return Failed{Reason: []error{err}}
	}

	if err := deps.GitConfigWriter.ReplaceAll(gitconfigscope.Global, "core.hooksPath", settings.HooksDir); err != nil {
		return Failed{Reason: []error{err}}
	}

	if err := deps.StateWriter.PersistEnabled(activationscope.Global, uniqueCoauthors); err != nil {
		return Failed{Reason: []error{err}}
	}

	return Succeeded{}
}

func applyAdditionalGuards(deps Dependencies, aliasesAndCoauthors []string) ([]string, []error) {
	coauthorCandidates, aliases := utils.Partition(aliasesAndCoauthors)

	sanityCheckErrs := deps.SanityCheckCoauthors(coauthorCandidates)
	if len(sanityCheckErrs) > 0 {
		return []string{}, sanityCheckErrs
	}

	resolvedAliases, resolveErrs := deps.GitResolveAliases(aliases)
	if len(resolveErrs) > 0 {
		return []string{}, resolveErrs
	}

	return append(coauthorCandidates, resolvedAliases...), []error{}
}

func removeDuplicates(coauthors []string) []string {
	var uniqueCoauthors []string
	temp := make(map[string]bool)
	for _, coauthor := range coauthors {
		temp[coauthor] = true
	}

	for coauthor := range temp {
		uniqueCoauthors = append(uniqueCoauthors, coauthor)
	}

	return uniqueCoauthors
}

func setupTemplate(deps Dependencies, commitTemplateBaseDir string, uniqueCoauthors []string) error {
	templateDir := fmt.Sprintf("%s/global", commitTemplateBaseDir)

	if err := deps.CreateTemplateDir(templateDir, os.ModePerm); err != nil {
		return err
	}

	commitTemplatePath := fmt.Sprintf("%s/COMMIT_TEMPLATE", templateDir)

	if err := deps.WriteTemplateFile(commitTemplatePath, []byte(utils.PrepareForCommitMessage(uniqueCoauthors)), 0644); err != nil {
		return err
	}

	if err := deps.GitConfigWriter.ReplaceAll(gitconfigscope.Global, "commit.template", commitTemplatePath); err != nil {
		return err
	}

	return nil
}
