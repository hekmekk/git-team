package datasource

import (
	"fmt"
	"os"

	"github.com/hekmekk/git-team/src/command/enable/commitsettings/entity"
)

type dependencies struct {
	getEnv func(string) string
}

// StaticValueDataSource reads configuration from constant values
type StaticValueDataSource struct {
	deps dependencies
}

// NewStaticValueDataSource constructs new StaticValueDataSource
func NewStaticValueDataSource() StaticValueDataSource {
	return newStaticValueDataSource(dependencies{getEnv: os.Getenv})
}

// for tests
func newStaticValueDataSource(deps dependencies) StaticValueDataSource {
	return StaticValueDataSource{deps: deps}
}

func (ds StaticValueDataSource) Read() entity.CommitSettings {
	homeDir := ds.deps.getEnv("HOME")

	cfg := entity.CommitSettings{
		TemplatesBaseDir: fmt.Sprintf("%s/.git-team/commit-templates", homeDir),
		HooksDir:         fmt.Sprintf("%s/.git-team/hooks", homeDir),
	}
	return cfg
}
