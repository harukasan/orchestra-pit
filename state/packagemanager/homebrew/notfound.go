package homebrew

// NotFound is an error which is used when the package is not found.
type NotFound struct {
	s string
}

func (e NotFound) Error() string {
	return e.s
}

func notFound(text string) error {
	return &NotFound{text}
}

// isNotFound returns true if whether the error is known to report that a
// package is not found.
func isNotFound(err error) bool {
	_, ok := err.(*NotFound)
	return ok
}
