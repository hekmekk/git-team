package gitconfig

import (
	"os/exec"
	"regexp"
	"strings"
)

const commitTemplate = "commit.template"
const teamAlias = "team.alias"

// Get git config --global --get <key>
func Get(key string) (string, error) {
	return get(execGitConfig)(key)
}

func get(exec func(...string) ([]string, error)) func(string) (string, error) {
	return func(key string) (string, error) {
		lines, err := exec("--get", key)
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
	return execGitConfig("--get-all", key)
}

// Add git config --global --add <key> <value>
func Add(key string, value string) error {
	_, err := execGitConfig("--add", key, value)
	return err
}

// ReplaceAll git config --global --replace-all <key> <value>
func ReplaceAll(key string, value string) error {
	_, err := execGitConfig("--replace-all", key, value)
	return err
}

// UnsetAll git config --global --unset-all <key>
func UnsetAll(key string) error {
	_, err := execGitConfig("--unset-all", key)
	return err
}

// GetRegexp git config --global --gex-regexp <pattern>
func GetRegexp(pattern string) (map[string]string, error) {
	return getRegexp(execGitConfig)(pattern)
}

func getRegexp(exec func(...string) ([]string, error)) func(string) (map[string]string, error) {
	return func(pattern string) (map[string]string, error) {
		mapping := make(map[string]string, 0)

		lines, err := exec("--get-regexp", pattern)
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

// execute /usr/bin/env git config --global <args>
func execGitConfig(args ...string) ([]string, error) {
	exec := func(theArgs ...string) ([]byte, error) {
		return exec.Command("/usr/bin/env", append([]string{"git"}, theArgs...)...).CombinedOutput()
	}

	return execGitConfigFactory(exec)(args...)
}

func execGitConfigFactory(cmd func(...string) ([]byte, error)) func(...string) ([]string, error) {
	return func(args ...string) ([]string, error) {
		gitArgs := append([]string{"config", "--global"}, args...)

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
