package events

import (
	config "github.com/hekmekk/git-team/src/shared/config/entity/config"
)

// RetrievalSucceeded successfully got the configuration
type RetrievalSucceeded struct {
	Config config.Config
}

// RetrievalFailed failed to read the configuration
type RetrievalFailed struct {
	Reason error
}

// SettingModificationFailed failed to modify the setting
type SettingModificationFailed struct {
	Reason error
}

// ReadingSingleSettingNotYetImplemented trying to access a not yet implemented feature
type ReadingSingleSettingNotYetImplemented struct {
}

// SettingModificationSucceeded suceesfully modified the setting
type SettingModificationSucceeded struct {
	Key   string
	Value string
}
