package errors

type stringer interface {
	String() string
}

// Causer is an interface
type Causer interface {
	Cause() error
}
