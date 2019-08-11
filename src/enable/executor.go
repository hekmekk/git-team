package enable

import (
	"fmt"
	"os"

	"github.com/hekmekk/git-team/src/config"
	utils "github.com/hekmekk/git-team/src/enable/utils"
	"github.com/hekmekk/git-team/src/validation"
)

// Command add a <Coauthor> under "team.alias.<Alias>"
type Command struct {
	Coauthors []string
}

// Dependencies the real-world dependencies of the ExecutorFactory
type Dependencies struct {
	LoadConfig        func() (config.Config, error)
	CreateDir         func(path string, perm os.FileMode) error
	WriteFile         func(path string, data []byte, mode os.FileMode) error
	SetCommitTemplate func(path string) error
	GitSetHooksPath   func(path string) error
	GitResolveAliases func(aliases []string) ([]string, []error)
	PersistEnabled    func(coauthors []string) error
}

// ExecutorFactory provisions a Command Processor
func ExecutorFactory(deps Dependencies) func(cmd Command) []error {
	return func(cmd Command) []error {
		if len(cmd.Coauthors) == 0 {
			return []error{}
		}

		coauthorCandidates, aliases := utils.Partition(cmd.Coauthors)

		sanityCheckErrs := validation.SanityCheckCoauthors(coauthorCandidates)
		if len(sanityCheckErrs) > 0 {
			return sanityCheckErrs
		}

		resolvedAliases, resolveErrs := deps.GitResolveAliases(aliases)
		if len(resolveErrs) > 0 {
			return resolveErrs
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
			return []error{err}
		}

		if err := deps.CreateDir(cfg.BaseDir, os.ModePerm); err != nil {
			return []error{err}
		}

		templatePath := fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.TemplateFileName)

		if err := deps.WriteFile(templatePath, []byte(utils.PrepareForCommitMessage(uniqueCoauthors)), 0644); err != nil {
			return []error{err}
		}
		if err := deps.SetCommitTemplate(templatePath); err != nil {
			return []error{err}
		}
		if err := deps.GitSetHooksPath(cfg.GitHooksPath); err != nil {
			return []error{err}
		}
		if err := deps.PersistEnabled(uniqueCoauthors); err != nil {
			return []error{err}
		}
		return []error{}
	}
}
