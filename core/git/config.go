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
	aliasFullPath := fmt.Sprintf("%s.%s", teamAlias, alias)
	coauthor, err := gitconfig.Local(aliasFullPath)
	if err != nil {
		coauthor, err = gitconfig.Global(aliasFullPath)
		if err != nil {
			return "", resolveErr(aliasFullPath)
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
	return nil
}

func resolveErr(aliasFullPath string) error {
	return errors.New(fmt.Sprintf("Failed to resolve alias %s", aliasFullPath))
}
