package git

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"syscall"

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
	return execGitConfig(commitTemplate, path)
}

func UnsetCommitTemplate() error {
	return execGitConfig("--unset", commitTemplate)
}

func RemoveCommitSection() error {
	return execGitConfig("--remove-section", "commit")
}

func AddAlias(alias, author string) error {
	return execGitConfig("--add", getAliasFullPath(alias), author)
}

func RemoveAlias(alias string) error {
	return execGitConfig("--unset-all", getAliasFullPath(alias))
}

func ListAlias() error {
	return execGitConfig("--get-regexp", teamAlias)
}

func getAliasFullPath(alias string) string {
	return fmt.Sprintf("%s.%s", teamAlias, alias)
}

func execGitConfig(args ...string) error {
	gitArgs := append([]string{"config", "--null", "--global"}, args...)
	var stdout bytes.Buffer
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdout = &stdout
	cmd.Stderr = ioutil.Discard

	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() == 1 {
				return errors.New(fmt.Sprintf("Failed to exec git config command with args: %s", args))
			}
		}
		return err
	}

	line, _ := stdout.ReadString([]byte("\n"))

	print(line)

	// print(stdout.String())

	// print(string(stdout.Bytes()))

	return nil
}
