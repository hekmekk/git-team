package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/hekmekk/git-team/src/add/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/add/interfaceadapter/event"
	"github.com/hekmekk/git-team/src/core/config"
	"github.com/hekmekk/git-team/src/core/effects"
	"github.com/hekmekk/git-team/src/core/events"
	git "github.com/hekmekk/git-team/src/core/gitconfig"
	"github.com/hekmekk/git-team/src/core/policy"
	"github.com/hekmekk/git-team/src/core/state_repository"
	"github.com/hekmekk/git-team/src/disable/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/disable/interfaceadapter/event"
	enableExecutor "github.com/hekmekk/git-team/src/enable"
	"github.com/hekmekk/git-team/src/remove/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/remove/interfaceadapter/event"
	"github.com/hekmekk/git-team/src/status"
	"github.com/hekmekk/git-team/src/status/interfaceadapter/cmd"
	"github.com/hekmekk/git-team/src/status/interfaceadapter/event"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	version = "v1.3.2-alpha2"
	author  = "Rea Sand <hekmek@posteo.de>"
)

func main() {
	application := newApplication(author, version)

	switch kingpin.MustParse(application.app.Parse(os.Args[1:])) {
	case application.add.CommandName:
		applyPolicy(application.add.Policy, addeventadapter.MapEventToEffects)
	case application.remove.CommandName:
		applyPolicy(application.remove.Policy, removeeventadapter.MapEventToEffects)
	case application.enable.command.FullCommand():
		runEnable(application.enable)
	case application.disable.CommandName:
		applyPolicy(application.disable.Policy, disableeventadapter.MapEventToEffectsFactory(application.status.Policy.Deps.StateRepositoryQuery))
	case application.status.CommandName:
		applyPolicy(application.status.Policy, statuseventadapter.MapEventToEffects)
	case application.list.command.FullCommand():
		runList()
	}

	os.Exit(0)
}

func applyPolicy(policy policy.Policy, adapter func(events.Event) []effects.Effect) {
	effects := adapter(policy.Apply())
	for _, effect := range effects {
		effect.Run()
	}
}

type enable struct {
	command             *kingpin.CmdClause
	aliasesAndCoauthors *[]string // can contain both aliases and coauthors
}

func newEnable(app *kingpin.Application) enable {
	command := app.Command("enable", "Enables injection of the provided co-authors whenever `git-commit` is used").Default()
	return enable{
		command:             command,
		aliasesAndCoauthors: command.Arg("coauthors", "The co-authors for the next commit(s). A co-author must either be an alias or of the shape \"Name <email>\"").Strings(),
	}
}

type list struct {
	command *kingpin.CmdClause
}

func newList(app *kingpin.Application) list {
	command := app.Command("ls", "List currently available aliases")
	command.Alias("list")
	return list{
		command: command,
	}
}

type application struct {
	app     *kingpin.Application
	add     addcmdadapter.Definition
	remove  removecmdadapter.Definition
	enable  enable
	disable disablecmdadapter.Definition
	status  statuscmdadapter.Definition
	list    list
}

func newApplication(author string, version string) application {
	app := kingpin.New("git-team", "Command line interface for managing and enhancing git commit messages with co-authors.")

	app.HelpFlag.Short('h')
	app.Version(version)
	app.Author(author)

	return application{
		app:     app,
		add:     addcmdadapter.NewDefinition(app),
		remove:  removecmdadapter.NewDefinition(app),
		enable:  newEnable(app),
		disable: disablecmdadapter.NewDefinition(app),
		status:  statuscmdadapter.NewDefinition(app),
		list:    newList(app),
	}
}

func runEnable(enable enable) {
	enableDeps := enableExecutor.Dependencies{
		CreateDir:         os.MkdirAll,           // TODO: CreateTemplateDir
		WriteFile:         ioutil.WriteFile,      // TODO: WriteTemplateFile
		SetCommitTemplate: git.SetCommitTemplate, // TODO: GitSetCommitTemplate
		GitSetHooksPath:   git.SetHooksPath,
		GitResolveAliases: git.ResolveAliases,
		PersistEnabled:    staterepository.PersistEnabled,
		LoadConfig:        config.Load,
	}
	execEnable := enableExecutor.ExecutorFactory(enableDeps)
	cmd := enableExecutor.Command{
		Coauthors: append(*enable.aliasesAndCoauthors),
	}
	enableErrs := execEnable(cmd)
	exitIfErr(enableErrs...)

	currState, err := staterepository.Query()
	exitIfErr(err)

	for _, effect := range statuseventadapter.MapEventToEffects(status.StateRetrievalSucceeded{State: currState}) {
		effect.Run()
	}
}

func runList() {
	assignments := git.GetAssignments()

	blackBold := color.New(color.FgBlack).Add(color.Bold)
	blackBold.Println("Aliases:")
	blackBold.Println("--------")

	var aliases []string

	for alias := range assignments {
		aliases = append(aliases, alias)
	}

	sort.Strings(aliases)

	for _, alias := range aliases {
		coauthor := assignments[alias]
		color.Magenta(fmt.Sprintf("'%s' -> '%s'", alias, coauthor))
	}
	os.Exit(0)
}

func exitIfErr(validationErrs ...error) {
	if len(validationErrs) > 0 && validationErrs[0] != nil {
		os.Stderr.WriteString(fmt.Sprintf("error: %s\n", foldErrors(validationErrs)))
		os.Exit(-1)
	}
}

// TODO: is this required for anything else than "enable"?
func foldErrors(validationErrors []error) error {
	var buffer bytes.Buffer
	for _, err := range validationErrors {
		buffer.WriteString(err.Error())
		buffer.WriteString("; ")
	}
	return errors.New(strings.TrimRight(buffer.String(), "; "))
}
