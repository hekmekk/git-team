package gitconfiginterface

// Reader read git configuration settings
type Reader interface {
	Get(key string) (string, error)
}
