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

// StateFactoryFunc is an adapter to use a function as StateFactory.
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

// Inspector is interface which retrieves the facts of the resource state.
//
// Get tries to retrieve information of the resource, and returns the facts or
// an error.
type Inspector interface {
	Get() (Facts, error)
}

// InspectorFunc is an adapter to use a function as an Inspector.
type InspectorFunc func() (Facts, error)

// Get implements Inspector interface to use as an Inspector.
func (f InspectorFunc) Get() (Facts, error) {
	return f()
}

// Facts is interface of the attributes of information of the resource.
type Facts interface {
	Get(name string) string
}
