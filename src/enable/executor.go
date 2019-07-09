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

		go func() {
			defer wg.Done()
			writeTemplateFileErr := effect.WriteFile(templatePath, []byte(coauthorsString), 0644)
			if writeTemplateFileErr != nil {
				mutex.Lock()
				errs = append(errs, writeTemplateFileErr)
				mutex.Unlock()
			}
		}()

		go func() {
			defer wg.Done()
			gitConfigErr := effect.SetCommitTemplate(templatePath)
			if gitConfigErr != nil {
				mutex.Lock()
				errs = append(errs, gitConfigErr)
				mutex.Unlock()
			}
		}()

		go func() {
			defer wg.Done()
			writeStatusFileErr := effect.SaveStatus(status.ENABLED, cmd.Coauthors...)
			if writeStatusFileErr != nil {
				mutex.Lock()
				errs = append(errs, writeStatusFileErr)
				mutex.Unlock()
			}
		}()

		wg.Wait()
		return errs
	}
}
