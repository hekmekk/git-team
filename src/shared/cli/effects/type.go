package effects

// Effect a side effect
type Effect interface {
	Run() error
}
