package gitconfiginterface

// Writer modify git configuration settings
type Writer interface {
	ReplaceAll(key string, value string) error
	UnsetAll(key string) error
}
