package gitconfig

import (
	"os/exec"
	"regexp"
	"strings"

	gitconfigerror "github.com/hekmekk/git-team/v2/src/shared/gitconfig/error"
	scope "github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
)

// Get git config --<scope> --get <key>
func Get(scope scope.Scope, key string) (string, error) {
	return get(execGitConfig)(scope, key)
}

func get(exec func(scope.Scope, ...string) ([]string, error)) func(scope.Scope, string) (string, error) {
	return func(scope scope.Scope, key string) (string, error) {
		lines, err := exec(scope, "--get", key)
		if err != nil {
			return "", err
		}

		if len(lines) == 0 {
			return "", nil
		}

		return lines[0], nil
	}
}

// GetAll git config --<scope> --get-all <key>
func GetAll(scope scope.Scope, key string) ([]string, error) {
	return execGitConfig(scope, "--get-all", key)
}

// Add git config --<scope> --add <key> <value>
func Add(scope scope.Scope, key string, value string) error {
	_, err := execGitConfig(scope, "--add", key, value)
	return err
}

// ReplaceAll git config --<scope> --replace-all <key> <value>
func ReplaceAll(scope scope.Scope, key string, value string) error {
	_, err := execGitConfig(scope, "--replace-all", key, value)
	return err
}

// UnsetAll git config --<scope> --unset-all <key>
func UnsetAll(scope scope.Scope, key string) error {
	_, err := execGitConfig(scope, "--unset-all", key)
	return err
}

// GetRegexp git config --<scope> --gex-regexp <pattern>
func GetRegexp(scope scope.Scope, pattern string) (map[string]string, error) {
	return getRegexp(execGitConfig)(scope, pattern)
}

func getRegexp(exec func(scope.Scope, ...string) ([]string, error)) func(scope.Scope, string) (map[string]string, error) {
	return func(scope scope.Scope, pattern string) (map[string]string, error) {
		mapping := make(map[string]string, 0)

		lines, err := exec(scope, "--get-regexp", pattern)
		if err != nil {
			return mapping, err
		}

		for _, line := range lines {
			keyAndValue := regexp.MustCompile("\\s").Split(line, 2)
			mapping[keyAndValue[0]] = keyAndValue[1]
		}

		return mapping, nil
	}
}

// List git config --<scope> --list
func List(scope scope.Scope) (map[string]string, error) {
	return list(execGitConfig)(scope)
}

func list(exec func(scope.Scope, ...string) ([]string, error)) func(scope.Scope) (map[string]string, error) {
	return func(scope scope.Scope) (map[string]string, error) {
		mapping := make(map[string]string, 0)

		lines, err := exec(scope, "--list")
		if err != nil {
			return mapping, err
		}

		for _, line := range lines {
			keyAndValue := regexp.MustCompile("=").Split(line, 2)
			mapping[keyAndValue[0]] = keyAndValue[1]
		}

		return mapping, nil
	}
}

// execute /usr/bin/env git config --<scope> <options>
func execGitConfig(scope scope.Scope, options ...string) ([]string, error) {
	gitConfigCommand := func(additionalOptions ...string) (string, error) {
		cmd := exec.Command("/usr/bin/env", append([]string{"git", "config"}, additionalOptions...)...)
		out, err := cmd.CombinedOutput()
		stringOut := string(out)
		return stringOut, gitconfigerror.New(err, cmd.String(), stringOut)
	}

	return execGitConfigFactory(gitConfigCommand)(scope, options...)
}

func execGitConfigFactory(cmd func(...string) (string, error)) func(scope.Scope, ...string) ([]string, error) {
	return func(scope scope.Scope, args ...string) ([]string, error) {
		gitArgs := append([]string{scope.Flag()}, args...)

		out, err := cmd(gitArgs...)

		if err != nil {
			return nil, err
		}

		if out == "" {
			return []string{}, nil
		}

		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")

		return lines, nil
	}
}
