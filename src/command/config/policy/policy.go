package policy

import (
	"fmt"

	configevents "github.com/hekmekk/git-team/src/command/config/events"
	"github.com/hekmekk/git-team/src/core/events"
	activationscope "github.com/hekmekk/git-team/src/shared/activation/scope"
	config "github.com/hekmekk/git-team/src/shared/config/interface"
)

// Request defines which config setting to modify or if the config should just be displayed
type Request struct {
	Key   *string
	Value *string
}

// Dependencies the dependencies of the config Policy module
type Dependencies struct {
	ConfigWriter config.Writer
	ConfigReader config.Reader
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
		return configevents.SettingModificationFailed{Reason: fmt.Errorf("Unknown setting '%s'", key)}
	}

	desiredScope := activationscope.FromString(value)
	if desiredScope == activationscope.Unknown {
		return configevents.SettingModificationFailed{Reason: fmt.Errorf("Unknown activation-scope '%s'", value)}
	}

	if err := deps.ConfigWriter.SetActivationScope(desiredScope); err != nil {
		return configevents.SettingModificationFailed{Reason: fmt.Errorf("Failed to modify setting 'activation-scope': %s", err)}
	}

	return configevents.SettingModificationSucceeded{Key: key, Value: value}
}
