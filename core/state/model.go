package state

type Status string

const (
	ENABLED  Status = "ENABLED"
	DISABLED Status = "DISABLED"
)

type State struct {
	Status    Status
	CoAuthors []string
}
