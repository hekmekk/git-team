package gitconfig

import (
	"os/exec"
	"regexp"
	"strings"

	scope "github.com/hekmekk/git-team/src/shared/gitconfig/scope"
)

// Get git config --global --get <key>
func Get(key string) (string, error) {
	return get(execGitConfig)(scope.Global, key)
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

// GetAll git config --global --get-all <key>
func GetAll(key string) ([]string, error) {
	return execGitConfig(scope.Global, "--get-all", key)
}

// Add git config --global --add <key> <value>
func Add(key string, value string) error {
	_, err := execGitConfig(scope.Global, "--add", key, value)
	return err
}

// ReplaceAll git config --global --replace-all <key> <value>
func ReplaceAll(key string, value string) error {
	_, err := execGitConfig(scope.Global, "--replace-all", key, value)
	return err
}

// UnsetAll git config --global --unset-all <key>
func UnsetAll(key string) error {
	_, err := execGitConfig(scope.Global, "--unset-all", key)
	return err
}

// GetRegexp git config --global --gex-regexp <pattern>
func GetRegexp(pattern string) (map[string]string, error) {
	return getRegexp(execGitConfig)(scope.Global, pattern)
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

// execute /usr/bin/env git config --<scope> <options>
func execGitConfig(scope scope.Scope, options ...string) ([]string, error) {
	gitConfigCommand := func(additionalOptions ...string) ([]byte, error) {
		return exec.Command("/usr/bin/env", append([]string{"git", "config"}, additionalOptions...)...).CombinedOutput()
	}

	return execGitConfigFactory(gitConfigCommand)(scope, options...)
}

func execGitConfigFactory(cmd func(...string) ([]byte, error)) func(scope.Scope, ...string) ([]string, error) {
	return func(scope scope.Scope, args ...string) ([]string, error) {
		gitArgs := append([]string{scope.Flag()}, args...)

		out, err := cmd(gitArgs...)

		if err != nil {
			return nil, err
		}

		stringOut := string(out)

		if stringOut == "" {
			return []string{}, nil
		}

		lines := strings.Split(strings.TrimRight(stringOut, "\n"), "\n")

		return lines, nil
	}
}
