package git

import (
	"fmt"
	"os/exec"
	"strings"
)

const commitTemplate = "commit.template"
const teamAlias = "team.alias"

func ResolveAlias(alias string) (string, error) {
	aliasFullPath := getAliasFullPath(alias)
	lines, err := ExecGitConfig("--get", getAliasFullPath(alias))
	if err != nil || len(lines) == 0 {
		return "", fmt.Errorf("Failed to resolve alias %s", aliasFullPath)
	}

	return lines[0], nil
}

func SetCommitTemplate(path string) error {
	_, err := ExecGitConfig(commitTemplate, path)
	return err
}

func UnsetCommitTemplate() error {
	_, err := ExecGitConfig("--unset", commitTemplate)
	return err
}

func RemoveCommitSection() error {
	_, err := ExecGitConfig("--remove-section", "commit")
	return err
}

func AddAlias(alias, author string) error {
	_, err := ExecGitConfig("--add", getAliasFullPath(alias), author)
	return err
}

func RemoveAlias(alias string) error {
	_, err := ExecGitConfig("--unset-all", getAliasFullPath(alias))
	return err
}

func GetAliasMap() map[string]string {
	return getAliasMap(ExecGitConfig)
}

func getAliasMap(exec func(...string) ([]string, error)) map[string]string {
	mapping := make(map[string]string)

	lines, err := exec("--get-regexp", teamAlias)
	if err != nil {
		lines = make([]string, 0)
	}

	for _, v := range lines {
		aliasAndCoauthor := strings.Split(strings.TrimRight(v, "\n"), "\n")
		mapping[strings.TrimPrefix(aliasAndCoauthor[0], fmt.Sprintf("%s.", teamAlias))] = aliasAndCoauthor[1]
	}

	return mapping
}

func getAliasFullPath(alias string) string {
	return fmt.Sprintf("%s.%s", teamAlias, alias)
}

func ExecGitConfig(args ...string) ([]string, error) {
	exec := func(theArgs ...string) ([]byte, error) {
		return exec.Command("/usr/bin/env", append([]string{"git"}, theArgs...)...).CombinedOutput()
	}

	return execGitConfig(exec)(args...)
}

func execGitConfig(cmd func(...string) ([]byte, error)) func(...string) ([]string, error) {
	return func(args ...string) ([]string, error) {
		gitArgs := append([]string{"config", "--null", "--global"}, args...)

		out, err := cmd(gitArgs...)

		if err != nil {
			return nil, fmt.Errorf("Failed to exec git config command with args: %s", args)
		}

		stringOut := string(out)

		if stringOut == "" {
			return []string{}, nil
		}

		lines := strings.Split(strings.TrimRight(stringOut, "\000"), "\000")

		return lines, nil
	}
}
