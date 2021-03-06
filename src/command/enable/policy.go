package enable

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"

	commitsettings "github.com/hekmekk/git-team/src/command/enable/commitsettings/interface"
	utils "github.com/hekmekk/git-team/src/command/enable/utils"
	"github.com/hekmekk/git-team/src/core/events"
	activation "github.com/hekmekk/git-team/src/shared/activation/interface"
	activationscope "github.com/hekmekk/git-team/src/shared/activation/scope"
	config "github.com/hekmekk/git-team/src/shared/config/interface"
	giterror "github.com/hekmekk/git-team/src/shared/gitconfig/error"
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
	GitConfigReader      gitconfig.Reader
	StateWriter          state.Writer
	GetEnv               func(string) string
	GetWd                func() (string, error)
	ActivationValidator  activation.Validator
}

// Request the coauthors with which to enable git-team
type Request struct {
	AliasesAndCoauthors *[]string
	UseAll              *bool
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

	var coAuthors []string
	if *req.UseAll {
		availableCoauthors, err := lookupAllCoauthors(deps)

		if err != nil {
			return Failed{Reason: []error{err}}
		}

		if len(availableCoauthors) == 0 {
			return Aborted{}
		}

		coAuthors = availableCoauthors

	} else {
		aliasesAndCoauthors := append(*req.AliasesAndCoauthors) // should be == *req.AliasesAndCoauthors

		if len(aliasesAndCoauthors) == 0 {
			return Aborted{}
		}

		coauthors, errs := applyAdditionalGuards(deps, aliasesAndCoauthors)
		if len(errs) > 0 {
			return Failed{Reason: errs}
		}

		coAuthors = removeDuplicates(coauthors)
	}

	settings := deps.CommitSettingsReader.Read()

	cfg, err := deps.ConfigReader.Read()
	if err != nil {
		return Failed{Reason: []error{err}}
	}

	activationScope := cfg.ActivationScope

	if activationScope == activationscope.RepoLocal && !deps.ActivationValidator.IsInsideAGitRepository() {
		return Failed{Reason: []error{fmt.Errorf("Failed to enable with activation-scope=%s: not inside a git repository", activationScope)}}
	}

	var gitConfigScope gitconfigscope.Scope
	if activationScope == activationscope.Global {
		gitConfigScope = gitconfigscope.Global
	} else {
		gitConfigScope = gitconfigscope.Local
	}

	if err := setupTemplate(gitConfigScope, deps, settings.TemplatesBaseDir, coAuthors); err != nil {
		return Failed{Reason: []error{err}}
	}

	if err := deps.GitConfigWriter.ReplaceAll(gitConfigScope, "core.hooksPath", settings.HooksDir); err != nil {
		return Failed{Reason: []error{err}}
	}

	if err := deps.StateWriter.PersistEnabled(cfg.ActivationScope, coAuthors); err != nil {
		return Failed{Reason: []error{err}}
	}

	return Succeeded{}
}

func lookupAllCoauthors(deps Dependencies) ([]string, error) {
	aliasCoauthorMap, err := deps.GitConfigReader.GetRegexp(gitconfigscope.Global, "team.alias")
	if err != nil && err.Error() != giterror.SectionOrKeyIsInvalid {
		return []string{}, err
	}

	coAuthors := []string{}

	for _, coauthor := range aliasCoauthorMap {
		coAuthors = append(coAuthors, coauthor)
	}

	return coAuthors, nil
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

func setupTemplate(gitConfigScope gitconfigscope.Scope, deps Dependencies, commitTemplateBaseDir string, uniqueCoauthors []string) error {
	var templateDir string

	if gitConfigScope == gitconfigscope.Local {
		user := deps.GetEnv("USER")
		workingDir, err := deps.GetWd()
		if err != nil {
			return err
		}
		templateDir = fmt.Sprintf("%s/repo-local/%s", commitTemplateBaseDir, determineRepoChecksum(user, workingDir))
	} else {
		templateDir = fmt.Sprintf("%s/global", commitTemplateBaseDir)
	}

	if err := deps.CreateTemplateDir(templateDir, os.ModePerm); err != nil {
		return err
	}

	commitTemplatePath := fmt.Sprintf("%s/COMMIT_TEMPLATE", templateDir)

	if err := deps.WriteTemplateFile(commitTemplatePath, []byte(utils.PrepareForCommitMessage(uniqueCoauthors)), 0644); err != nil {
		return err
	}

	if err := deps.GitConfigWriter.ReplaceAll(gitConfigScope, "commit.template", commitTemplatePath); err != nil {
		return err
	}

	return nil
}

func determineRepoChecksum(user string, repoPath string) string {
	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("%s:%s", user, repoPath)))
	return hex.EncodeToString(hasher.Sum(nil))
}
