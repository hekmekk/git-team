package git

import (
	"errors"
	"fmt"
	"gopkg.in/libgit2/git2go.v24"
	"strings"
)

const commitTemplate = "commit.template"

func ResolveAlias(alias string) (string, error) {
	globalConfig, err := getGlobalConfig()
	if err != nil {
		return "", err
	}
	coauthor, err := globalConfig.LookupString(fmt.Sprintf("team.alias.%s", alias))
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to resolve alias %s", alias))
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
