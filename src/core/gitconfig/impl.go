package gitconfig

// DataSource read data directly from gitconfig
type DataSource struct {
}

// NewDataSource construct new DataSource
func NewDataSource() DataSource {
	return DataSource{}
}

// Get read a setting
func (ds DataSource) Get(key string) (string, error) {
	return Get(key)
}

// DataSink write data directly to gitconfig
type DataSink struct {
}

// NewDataSink construct new DataSink
func NewDataSink() DataSink {
	return DataSink{}
}

// Set modify a setting
func (ds DataSink) Set(key string, value string) error {
	return ReplaceAll(key, value)
}

// Unset remove a setting
func (ds DataSink) Unset(key string) error {
	return UnsetAll(key)
}
