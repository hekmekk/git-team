package handler

import (
	"fmt"
	"github.com/hekmekk/git-team/core/config"
	"github.com/hekmekk/git-team/core/git"
	"github.com/hekmekk/git-team/core/status"
	"io/ioutil"
	"os"
	"sync"
)

func EnableCommand(coauthors *[]string) {
	cfg, _ := config.Load()

	defer status.Print()

	if len(*coauthors) == 0 {
		return
	}

	validCoAuthors, userInputErrors := ValidateUserInput(coauthors)
	ToStderrAndExit(userInputErrors...)

	coauthorsString := PrepareForCommitMessage(validCoAuthors)

	mkdirErr := os.MkdirAll(cfg.BaseDir, os.ModePerm)
	ToStderrAndExit(mkdirErr)

	templatePath := fmt.Sprintf("%s/%s", cfg.BaseDir, cfg.TemplateFileName)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		writeTemplateFileErr := ioutil.WriteFile(templatePath, []byte(coauthorsString), 0644)
		ToStderrAndExit(writeTemplateFileErr)
	}()

	go func() {
		defer wg.Done()
		gitConfigErr := git.SetCommitTemplate(templatePath)
		ToStderrAndExit(gitConfigErr)
	}()

	go func() {
		defer wg.Done()
		writestatusFileErr := status.Save(status.ENABLED, validCoAuthors...)
		ToStderrAndExit(writestatusFileErr)
	}()

	wg.Wait()
}
