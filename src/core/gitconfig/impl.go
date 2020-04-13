package gitconfig

// RawDataSource read data directly from gitconfig
type RawDataSource struct {
}

// NewRawDataSource construct new RawDataSource
func NewRawDataSource() RawDataSource {
	return RawDataSource{}
}

// Get read a setting
func (ds RawDataSource) Get(key string) (string, error) {
	return Get(key)
}

// RawDataSink write data directly to gitconfig
type RawDataSink struct {
}

// NewRawDataSink construct new RawDataSink
func NewRawDataSink() RawDataSink {
	return RawDataSink{}
}

// Set modify a setting
func (ds RawDataSink) Set(key string, value string) error {
	return ReplaceAll(key, value)
}
