package enable

import (
	"fmt"
	"os"
	"sync"

	"github.com/hekmekk/git-team/core/status"
)

type EnableCommand struct {
	Coauthors        []string
	BaseDir          string
	TemplateFileName string
}

type EnableEffect struct {
	CreateDir         func(path string, perm os.FileMode) error
	WriteFile         func(path string, data []byte, mode os.FileMode) error
	SetCommitTemplate func(path string) error
	SaveStatus        func(state status.State, coauthors ...string) error
}

func ExecutorFactory(effect EnableEffect) func(cmd EnableCommand) []error {
	return func(cmd EnableCommand) []error {
		errs := make([]error, 0)
		if len(cmd.Coauthors) == 0 {
			return errs
		}

		coauthorsString := PrepareForCommitMessage(cmd.Coauthors)

		mkdirErr := effect.CreateDir(cmd.BaseDir, os.ModePerm)
		if mkdirErr != nil {
			return append(errs, mkdirErr)
		}

		templatePath := fmt.Sprintf("%s/%s", cmd.BaseDir, cmd.TemplateFileName)

		mutex := &sync.Mutex{}
		var wg sync.WaitGroup
		wg.Add(3)

		runTask := runTaskFactory(&wg, mutex, errs)

		go runTask(func() error { return effect.WriteFile(templatePath, []byte(coauthorsString), 0644) })
		go runTask(func() error { return effect.SetCommitTemplate(templatePath) })
		go runTask(func() error { return effect.SaveStatus(status.ENABLED, cmd.Coauthors...) })

		wg.Wait()
		return errs
	}
}

func runTaskFactory(wg *sync.WaitGroup, mutex *sync.Mutex, errs []error) func(effect func() error) {
	return func(effect func() error) {
		defer wg.Done()
		if err := effect(); err != nil {
			mutex.Lock()
			errs = append(errs, err)
			mutex.Unlock()
		}
	}
}
