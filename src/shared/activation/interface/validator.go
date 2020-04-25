package activationvalidatorinterface

// Validator check activation validation properties
type Validator interface {
	IsInsideAGitRepository() bool
}
