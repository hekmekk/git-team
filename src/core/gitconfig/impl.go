package gitconfig

// RawDataSource read data from gitconfig
type RawDataSource struct {
}

// NewRawDataSource construct new RawDataSource
func NewRawDataSource() RawDataSource {
	return RawDataSource{}
}

// Get git config --global --get <key>
func (ds RawDataSource) Get(key string) (string, error) {
	return Get(key)
}
