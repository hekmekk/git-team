package gitconfig

// RawReader read data from gitconfig
type RawReader interface {
	Get(key string) (string, error)
}
