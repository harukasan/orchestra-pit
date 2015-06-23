package commands

// State is interface of the state of the resource.
//
// Apply tries to change the state to the desired. Apply returns an error,
// if failed to modify the state.
//
// Test tests whether the state is the same to requested. If the state is
// different to requested state, Test returns an error.
type State interface {
	Apply() error
	Test() error
}

// StateFactory is interface of the factory struct of the state.
//
// New returns the initialized State object with given options. If failed to
// creates a new state or passed invalid options, it returns nil and an error.
type StateFactory interface {
	New(options Options) (State, error)
}

// StateFactoryFunc is function that creates and returns a new State with given
// options parameters. If given options is invalid, StateFactory returns nil and
// an error.
type StateFactoryFunc func(options Options) (State, error)

// New implements the StateFactory interface to use a StateFactoryFunc as
// StateFactory interface.
func (f StateFactoryFunc) New(options Options) (State, error) {
	return f(options)
}

// Options is interface of the parameters of the initialize function of the state.
//
// Get returns the value string of the named parameter.
type Options interface {
	Get(name string) string
}
