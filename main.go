package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	addExecutor "github.com/hekmekk/git-team/src/add"
	"github.com/hekmekk/git-team/src/config"
	execDisable "github.com/hekmekk/git-team/src/disable"
	enableExecutor "github.com/hekmekk/git-team/src/enable"
	git "github.com/hekmekk/git-team/src/gitconfig"
	removeExecutor "github.com/hekmekk/git-team/src/remove"
	statusApi "github.com/hekmekk/git-team/src/status"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	version = "v1.1.1"
	author  = "Rea Sand <hekmek@posteo.de>"
)

var (
	addDeps = addExecutor.Dependencies{
		AddGitAlias: git.AddAlias,
	}
	enableDeps = enableExecutor.Dependencies{
		CreateDir:         os.MkdirAll,           // TODO: CreateTemplateDir
		WriteFile:         ioutil.WriteFile,      // TODO: WriteTemplateFile
		SetCommitTemplate: git.SetCommitTemplate, // TODO: GitSetCommitTemplate
		GitResolveAliases: git.ResolveAliases,
		PersistEnabled:    statusApi.PersistEnabled,
		LoadConfig:        config.Load,
	}
	rmDeps = removeExecutor.Dependencies{
		GitResolveAlias: git.ResolveAlias,
		GitRemoveAlias:  git.RemoveAlias,
	}
	execAdd    = addExecutor.ExecutorFactory(addDeps)
	execEnable = enableExecutor.ExecutorFactory(enableDeps)
	execRemove = removeExecutor.ExecutorFactory(rmDeps)
)

func main() {
	app := kingpin.New("git-team", "Command line interface for creating git commit templates provisioned with one or more co-authors. Please note that \"git commit -m\" is not affected by commit templates.")

	app.HelpFlag.Short('h')
	app.Version(version)
	app.Author(author)

	enable := app.Command("enable", "Provisions a git-commit template with the provided co-authors. A co-author must either be an alias or of the shape \"Name <email>\"").Default()
	enableCoauthors := enable.Arg("coauthors", "Git co-authors").Strings()

	disable := app.Command("disable", "Use default template")
	status := app.Command("status", "Print the current status")

	add := app.Command("add", "Add an alias")
	addAlias := add.Arg("alias", "The alias to be added").Required().String()
	addCoauthor := add.Arg("coauthor", "The co-author").Required().String()

	rm := app.Command("rm", "Remove an alias")
	rmAlias := rm.Arg("alias", "The alias to be removed").Required().String()

	list := app.Command("list", "List currently available aliases")
	list.Alias("ls")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case enable.FullCommand():
		cmd := enableExecutor.Command{
			Coauthors: append(*enableCoauthors),
		}
		enableErrs := execEnable(cmd)
		exitIfErr(enableErrs...)

		status, err := statusApi.Fetch()
		exitIfErr(err)

		fmt.Println(status.ToString())
		os.Exit(0)
	case disable.FullCommand():
		err := execDisable.Exec()
		exitIfErr(err)

		status, err := statusApi.Fetch()
		exitIfErr(err)

		fmt.Println(status.ToString())
		os.Exit(0)
	case status.FullCommand():
		status, err := statusApi.Fetch()
		exitIfErr(err)

		fmt.Println(status.ToString())
		os.Exit(0)
	case add.FullCommand():
		cmd := addExecutor.Command{
			Alias:    *addAlias,
			Coauthor: *addCoauthor,
		}
		addErr := execAdd(cmd)
		exitIfErr(addErr)

		color.Green(fmt.Sprintf("Alias '%s' -> '%s' has been added.", *addAlias, *addCoauthor))
		os.Exit(0)
	case rm.FullCommand():
		cmd := removeExecutor.Command{
			Alias: *rmAlias,
		}

		rmErr := execRemove(cmd)
		exitIfErr(rmErr)

		color.Red(fmt.Sprintf("Alias '%s' has been removed.", cmd.Alias))
		os.Exit(0)
	case list.FullCommand():
		assignments := git.GetAddedAliases()

		blackBold := color.New(color.FgBlack).Add(color.Bold)
		blackBold.Println("Aliases:")
		blackBold.Println("--------")

		for alias, coauthor := range assignments {
			color.Magenta(fmt.Sprintf("'%s' -> '%s'", alias, coauthor))
		}
		os.Exit(0)
	}
}

func exitIfErr(validationErrs ...error) {
	if len(validationErrs) > 0 && validationErrs[0] != nil {
		os.Stderr.WriteString(fmt.Sprintf("error: %s\n", foldErrors(validationErrs)))
		os.Exit(-1)
	}
}

func foldErrors(validationErrors []error) error {
	var buffer bytes.Buffer
	for _, err := range validationErrors {
		buffer.WriteString(err.Error())
		buffer.WriteString("; ")
	}
	return errors.New(strings.TrimRight(buffer.String(), "; "))
}
