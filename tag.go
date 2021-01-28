package errors

// ErrorTag ...
type ErrorTag struct {
	_ interface{} // fix issue with empty struct always has same addr
}

// Predefined tags
var (
	ApplicationError   = NewTag()
	LogicError         = NewTag()
	UnmarshalError     = NewTag()
	ValidationError    = NewTag()
	ErrBreak           = NewTag()
	AlreadyExistsError = NewTag()
	PlannedError       = NewTag() // for example: accessory check (404)
	TransportError     = NewTag()
)

// NewTag ...
func NewTag() *ErrorTag {
	return &ErrorTag{}
}

// WithTag ...
func (e *Error) WithTag(et *ErrorTag) *Error {
	if !e.IsNil() {
		e.tags = append(e.tags, et)
	}
	return e
}

// IsTagged ...
func (et *ErrorTag) IsTagged(err error) bool {
	if err == nil {
		return false
	}
	if appErr, ok := err.(*Error); ok {
		return appErr.IsTagged(et)
	}
	return false
}
