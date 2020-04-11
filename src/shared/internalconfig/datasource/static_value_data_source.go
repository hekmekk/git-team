package datasource

import (
	"fmt"
	"os"

	"github.com/hekmekk/git-team/src/shared/internalconfig/entity"
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

func (ds StaticValueDataSource) Read() entity.InternalConfig {
	cfg := entity.InternalConfig{
		GitTeamCommitTemplatePath: fmt.Sprintf("%s/.config/git-team/COMMIT_TEMPLATE", ds.deps.getEnv("HOME")),
		GitTeamHooksPath:          "/usr/local/etc/git-team/hooks",
	}
	return cfg
}
