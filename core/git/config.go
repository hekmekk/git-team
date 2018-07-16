package git

import (
	"errors"
	"fmt"
	"gopkg.in/libgit2/git2go.v24"
	"os"
	"strings"
)

const commitTemplate = "commit.template"
const teamAlias = "team.alias"

func ResolveAlias(alias string) (string, error) {
	aliasFullPath := fmt.Sprintf("%s.%s", teamAlias, alias)
	coauthor, err := resolveAlias(getGlobalConfig)(aliasFullPath)
	if err != nil {
		return resolveAlias(getRepoLocalConfig)(aliasFullPath)
	}
	return coauthor, nil
}

func SetCommitTemplate(path string) error {
	globalConfig, err := getGlobalConfig()
	if err != nil {
		return err
	}

	return globalConfig.SetString(commitTemplate, path)
}

func UnsetCommitTemplate() error {
	globalConfig, err := getGlobalConfig()
	if err != nil {
		return err
	}
	_, err = globalConfig.LookupString(commitTemplate)
	if err != nil {
		return nil
	}

	return globalConfig.Delete(commitTemplate)
}

func resolveAlias(configProvider func() (*git.Config, error)) func(string) (string, error) {
	return func(aliasFullPath string) (string, error) {
		config, err := configProvider()
		if err != nil {
			return "", resolveErr(aliasFullPath)
		}
		coauthor, err := config.LookupString(aliasFullPath)
		if err != nil {
			return "", resolveErr(aliasFullPath)
		}
		return strings.TrimRight(coauthor, "\n"), nil
	}
}

func resolveErr(aliasFullPath string) error {
	return errors.New(fmt.Sprintf("Failed to resolve alias %s", aliasFullPath))
}

func getGlobalConfig() (*git.Config, error) {
	globalConfigPath, err := git.ConfigFindGlobal()
	if err != nil {
		return nil, err
	}

	return git.OpenOndisk(nil, globalConfigPath)
}

func getRepoLocalConfig() (*git.Config, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	repo, err := git.OpenRepository(workDir)
	if err != nil {
		return nil, errors.New("Not a git repository")
	}

	return repo.Config()
}
