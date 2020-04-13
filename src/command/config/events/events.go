package events

import (
	config "github.com/hekmekk/git-team/src/command/config/entity/config"
)

// RetrievalSucceeded successfully got the configuration
type RetrievalSucceeded struct {
	Config config.Config
}

// RetrievalFailed failed to read the configuration
type RetrievalFailed struct {
	Reason error
}

// SettingModificationFailed failed to modify the a setting
type SettingModificationFailed struct {
	Reason error
}

// ReadingSingleSettingNotYetImplemented trying to access a not yet implemented feature
type ReadingSingleSettingNotYetImplemented struct {
}
