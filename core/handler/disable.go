package handler

import (
	"fmt"
	"github.com/hekmekk/git-team/core/git"
	"github.com/hekmekk/git-team/core/state"
	"os"
	"sync"
)

func DisableCommand(baseDir, templateFile, stateFile string) {
	defer state.Print(baseDir, stateFile)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		if err := state.Save(baseDir, stateFile, state.DISABLED); err != nil {
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
		os.Remove(fmt.Sprintf("%s/%s", baseDir, templateFile))
	}()

	wg.Wait()
}
