package gitconfiginterface

import (
	"github.com/hekmekk/git-team/v2/src/shared/gitconfig/scope"
)

// Writer modify git configuration settings
type Writer interface {
	Add(scope scope.Scope, key string, value string) error
	ReplaceAll(scope scope.Scope, key string, value string) error
	UnsetAll(scope scope.Scope, key string) error
}
