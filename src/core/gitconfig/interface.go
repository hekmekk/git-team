package gitconfig

// SettingsReader read git configuration settings
type SettingsReader interface {
	Get(key string) (string, error)
}

// SettingsWriter modify git configuration settings
type SettingsWriter interface {
	Set(key string, value string) error
	Unset(key string) error
}
