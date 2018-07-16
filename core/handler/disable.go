package handler

import (
	"fmt"
	"github.com/hekmekk/git-team/core/config"
	"github.com/hekmekk/git-team/core/git"
	"github.com/hekmekk/git-team/core/state"
	"os"
	"sync"
)

func DisableCommand() {
	defer state.Print()

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		if err := state.Save(state.DISABLED); err != nil {
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
		cfg, _ := config.Load()
		os.Remove(fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.TemplateFileName))
	}()

	wg.Wait()
}
