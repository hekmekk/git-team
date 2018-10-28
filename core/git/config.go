package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/tcnksm/go-gitconfig"
)

const commitTemplate = "commit.template"
const teamAlias = "team.alias"

func ResolveAlias(alias string) (string, error) {
	aliasFullPath := getAliasFullPath(alias)
	coauthor, err := gitconfig.Local(aliasFullPath)
	if err != nil {
		coauthor, err = gitconfig.Global(aliasFullPath)
		if err != nil {
			return "", errors.New(fmt.Sprintf("Failed to resolve alias %s", aliasFullPath))
		}
	}
	return coauthor, nil
}

func SetCommitTemplate(path string) error {
	_, err := execGitConfig(commitTemplate, path)
	return err
}

func UnsetCommitTemplate() error {
	_, err := execGitConfig("--unset", commitTemplate)
	return err
}

func RemoveCommitSection() error {
	_, err := execGitConfig("--remove-section", "commit")
	return err
}

func AddAlias(alias, author string) error {
	_, err := execGitConfig("--add", getAliasFullPath(alias), author)
	return err
}

func RemoveAlias(alias string) error {
	_, err := execGitConfig("--unset-all", getAliasFullPath(alias))
	return err
}

func ListAlias() ([]string, error) {
	return execGitConfig("--get-regexp", teamAlias)
}

func getAliasFullPath(alias string) string {
	return fmt.Sprintf("%s.%s", teamAlias, alias)
}

func execGitConfig(args ...string) ([]string, error) {
	gitArgs := append([]string{"config", "--null", "--global"}, args...)
	out, err := exec.Command("/usr/bin/env", append([]string{"git"}, gitArgs...)...).CombinedOutput()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to exec git config command with args: %s", args))
	}

	stringOut := string(out)

	if stringOut == "" {
		return []string{}, nil
	}

	lines := strings.Split(strings.TrimRight(stringOut, "\000"), "\000")

	return lines, nil
}
