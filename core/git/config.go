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
	coauthor, err := resolveAliasFromGlobalConfig(aliasFullPath)
	if err != nil {
		return resolveAliasFromRepoLocalConfig(aliasFullPath)
	}
	return coauthor, nil
}

func resolveAliasFromGlobalConfig(aliasFullPath string) (string, error) {
	globalConfig, err := getGlobalConfig()
	if err != nil {
		return "", err
	}
	coauthor, err := globalConfig.LookupString(aliasFullPath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to resolve alias %s", aliasFullPath))
	}
	return strings.TrimRight(coauthor, "\n"), nil
}

func resolveAliasFromRepoLocalConfig(aliasFullPath string) (string, error) {
	repoLocalConfig, err := getRepoLocalConfig()
	if err != nil {
		return "", err
	}
	coauthor, err := repoLocalConfig.LookupString(aliasFullPath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to resolve alias %s", aliasFullPath))
	}
	return strings.TrimRight(coauthor, "\n"), nil
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

func lookupEntry(config *git.Config, key string) (string, error) {
	return config.LookupString(key)
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
