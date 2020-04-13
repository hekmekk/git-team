package gitconfig

// RawReader read data from gitconfig
type RawReader interface {
	Get(key string) (string, error)
}

// RawWriter set data to gitconfig
type RawWriter interface {
	Set(key string, value string) error
}
