package gitconfig

import (
	"fmt"
	"os/exec"
	"strings"
)

const hooksPath = "core.hooksPath"
const commitTemplate = "commit.template"
const teamAlias = "team.alias"

// SetHooksPath set your global "core.hooksPath"
func SetHooksPath(path string) error {
	_, err := execGitConfig("--add", hooksPath, path)
	return err
}

// UnsetHooksPath unset your global "core.hooksPath"
func UnsetHooksPath() error {
	_, err := execGitConfig("--unset", hooksPath)
	return err
}

// SetCommitTemplate set your global "commit.template" globally
func SetCommitTemplate(path string) error {
	_, err := execGitConfig("--add", commitTemplate, path)
	return err
}

// UnsetCommitTemplate unset your global "commit.template"
func UnsetCommitTemplate() error {
	_, err := execGitConfig("--unset", commitTemplate)
	return err
}

// AddAlias add a co-author for "team.alias.<alias>"
func AddAlias(alias, author string) error {
	_, err := execGitConfig("--add", getAliasFullPath(alias), author)
	return err
}

// RemoveAlias remove "team.alias.<alias>"
func RemoveAlias(alias string) error {
	_, err := execGitConfig("--unset-all", getAliasFullPath(alias))
	return err
}

// GetAssignments get all alias -> co-author mappings
func GetAssignments() map[string]string {
	return getAssignments(execGitConfig)
}

func getAssignments(exec func(...string) ([]string, error)) map[string]string {
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

// ResolveAliases convenience function to resolve multiple aliases and accumulate errors
func ResolveAliases(aliases []string) ([]string, []error) {
	return resolveAliases(ResolveAlias)(aliases)
}

func resolveAliases(resolveAlias func(string) (string, error)) func([]string) ([]string, []error) {
	return func(aliases []string) ([]string, []error) {
		var resolvedAliases []string
		var resolveErrors []error

		for _, alias := range aliases {
			var resolvedCoauthor, err = resolveAlias(alias)
			if err != nil {
				resolveErrors = append(resolveErrors, err)
			} else {
				resolvedAliases = append(resolvedAliases, resolvedCoauthor)
			}
		}

		return resolvedAliases, resolveErrors
	}
}

// ResolveAlias lookup "team.alias.<alias>" globally
func ResolveAlias(alias string) (string, error) {
	return resolveAlias(execGitConfig)(alias)
}

func resolveAlias(exec func(...string) ([]string, error)) func(string) (string, error) {
	return func(alias string) (string, error) {
		aliasFullPath := getAliasFullPath(alias)
		lines, err := exec("--get", aliasFullPath)
		if err != nil || len(lines) == 0 {
			return "", fmt.Errorf("Failed to resolve alias %s", aliasFullPath)
		}

		return lines[0], nil
	}
}

func getAliasFullPath(alias string) string {
	return fmt.Sprintf("%s.%s", teamAlias, alias)
}

// execute /usr/bin/env git config --null --global <args>
func execGitConfig(args ...string) ([]string, error) {
	exec := func(theArgs ...string) ([]byte, error) {
		return exec.Command("/usr/bin/env", append([]string{"git"}, theArgs...)...).CombinedOutput()
	}

	return execGitConfigFactory(exec)(args...)
}

func execGitConfigFactory(cmd func(...string) ([]byte, error)) func(...string) ([]string, error) {
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
