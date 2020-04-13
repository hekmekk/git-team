package policy

import (
	"errors"
	"fmt"

	configevents "github.com/hekmekk/git-team/src/command/config/events"
	configreader "github.com/hekmekk/git-team/src/command/config/reader"
	configwriter "github.com/hekmekk/git-team/src/command/config/writer"
	"github.com/hekmekk/git-team/src/core/events"
)

// Request defines which config setting to modify or if the config should just be displayed
type Request struct {
	Key   *string
	Value *string
}

// Dependencies the dependencies of the config Policy module
type Dependencies struct {
	ConfigReader configreader.ConfigReader
	ConfigWriter configwriter.ConfigWriter
}

// Policy the policy to apply
type Policy struct {
	Req  Request
	Deps Dependencies
}

// Apply Edit or show configuration
func (policy Policy) Apply() events.Event {
	deps := policy.Deps
	req := policy.Req
	keyPtr := req.Key
	valuePtr := req.Value

	if (keyPtr == nil || *keyPtr == "") && (valuePtr == nil || *valuePtr == "") {
		cfg, err := deps.ConfigReader.Read()
		if err != nil {
			return configevents.RetrievalFailed{Reason: err}
		}

		return configevents.RetrievalSucceeded{Config: cfg}
	}

	if keyPtr == nil || *keyPtr == "" || valuePtr == nil || *valuePtr == "" {
		return configevents.ReadingSingleSettingNotYetImplemented{}
	}

	key := *keyPtr
	value := *valuePtr

	if key != "activation-scope" {
		return configevents.SettingModificationFailed{Reason: fmt.Errorf("unknown setting '%s'", key)}
	}

	if value == "global" || value == "repo-local" {
		return configevents.SettingModificationFailed{Reason: errors.New("not yet")}
	}

	return configevents.SettingModificationFailed{Reason: fmt.Errorf("unknown activation-scope '%s'", value)}
}
