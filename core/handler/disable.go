package handler

import (
	"fmt"
	"github.com/hekmekk/git-team/core/status"
	"github.com/hekmekk/git-team/src/config"
	git "github.com/hekmekk/git-team/src/gitconfig"
	"os"
	"sync"
)

func DisableCommand() {
	defer status.Print()

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		if err := status.Save(status.DISABLED); err != nil {
			os.Stderr.WriteString(fmt.Sprintf("error: %s\n", err))
			os.Exit(-1)
		}
	}()

	go func() {
		defer wg.Done()
		git.UnsetCommitTemplate()
	}()

	go func() {
		defer wg.Done()
		git.RemoveCommitSection()
	}()

	go func() {
		defer wg.Done()
		cfg, _ := config.Load()
		os.Remove(fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.TemplateFileName))
	}()

	wg.Wait()
}
