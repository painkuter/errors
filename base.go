package errors

import (
	"fmt"
	"io"

	"google.golang.org/grpc/codes"
)

const nilErrorMessage = "<nil_err>"

// Error ...
type Error struct {
	message    string
	errorType  *ErrorType
	code       string
	cause      error
	stackTrace []string
	tags       []*ErrorTag
	properties PropertyMap
	details    pairsList
	log        pairsList
}

// Error ...
func (e *Error) Error() string {
	if !e.IsNil() {
		return e.message
	}
	return nilErrorMessage
}

// IsNil ...
func (e *Error) IsNil() bool {
	return e == nil || e.errorType == NullError
}

func IsNil(e error) bool {
	if e == nil {
		return true
	}
	err, ok := Unwrap(e)
	return ok && err.IsNil()
}

// IsTyped ...
func (e *Error) IsTyped(et *ErrorType) bool {
	if !e.IsNil() {
		return e.errorType == et
	}
	return false
}

// IsClientErr checks is e client error
func (e *Error) IsClientErr() bool {
	if e.IsNil() {
		return false
	}

	for _, errorType := range clientErrorTypes {
		if e.errorType == errorType {
			return true
		}
	}

	return false
}

// IsClientErr is wrapper for *Error.IsClientErr() with unwrapping
func IsClientErr(e error) bool {
	if e == nil {
		return false
	}
	err, ok := Unwrap(e)
	return ok && err.IsClientErr()
}

// IsTagged ...
func (e *Error) IsTagged(et *ErrorTag) bool {
	if !e.IsNil() && len(e.tags) > 0 {
		for _, tag := range e.tags {
			if tag == et {
				return true
			}
		}
	}
	return false
}

// WithMessage ...
func (e *Error) WithMessage(message string) *Error {
	if !e.IsNil() {
		e.message = message
	}
	return e
}

// WithNotFoundMessage ...
func (e *Error) WithNotFoundMessage(message string) *Error {
	if !e.IsNil() && e.errorType == NotFound {
		e.message = message
	}
	return e
}

// WithMessageForType ...
func (e *Error) WithMessageForType(message string, errorType *ErrorType) *Error {
	if !e.IsNil() && e.errorType == errorType {
		e.message = message
	}
	return e
}

// WithCause ...
func (e *Error) WithCause(cause error) *Error {
	if !e.IsNil() {
		e.cause = cause
	}
	return e
}

// WithDetailsKV ...
func (e *Error) WithDetailsKV(pairs ...interface{}) *Error {
	if !e.IsNil() {
		details := newPairsList(pairs...)
		e.details = append(e.details, details...)
		e.log = append(e.log, details...)
	}
	return e
}

// WithLogKV ...
func (e *Error) WithLogKV(pairs ...interface{}) *Error {
	if !e.IsNil() {
		e.log = append(e.log, newPairsList(pairs...)...)
	}
	return e
}

// WithCode ...
func (e *Error) WithCode(s string) *Error {
	if !e.IsNil() {
		e.code = s
	}
	return e
}

// WithProperty ...
func (e *Error) WithProperty(p *Property, value interface{}) *Error {
	if !e.IsNil() {
		e.properties[p] = &PropertyValue{value}
	}
	return e
}

// WithProperties ...
func (e *Error) WithProperties(properties map[*Property]interface{}) *Error {
	if e.IsNil() {
		return e
	}
	for p, value := range properties {
		e.properties[p] = &PropertyValue{value}
	}
	return e
}

// Cause ...
func (e *Error) Cause() error {
	if !e.IsNil() {
		return e.cause
	}
	return nil
}

// GRPCCode ...
func (e *Error) GRPCCode() codes.Code {
	if !e.IsNil() {
		return e.errorType.code
	}
	return codes.OK
}

// ResponseCode ...
func (e *Error) ResponseCode() string {
	code, ok := e.Property(ResponseCode).(string)
	if !ok {
		return ""
	}
	return code
}

// LogDetails ...
func (e *Error) LogDetails() []interface{} {
	if e.IsNil() {
		return nil
	}
	return e.log.Flattened()
}

// Details ...
func (e *Error) Details() []interface{} {
	if e.IsNil() {
		return nil
	}
	return e.details.Flattened()
}

// Property ...
func (e *Error) Property(p *Property) interface{} {
	if e.IsNil() {
		return nil
	}
	if val, ok := e.properties[p]; ok {
		return val.Value()
	}
	return nil
}

// PropertyValue ...
func (e *Error) PropertyValue(p *Property) *PropertyValue {
	if e.IsNil() {
		return nil
	}
	if val, ok := e.properties[p]; ok {
		return val
	}
	return nil
}

// LogDetails ...
func (e *Error) Code() string {
	if e.IsNil() {
		return "OK"
	}
	return e.code
}

// Format ...
func (e *Error) Format(s fmt.State, verb rune) {
	if e.IsNil() {
		_, _ = io.WriteString(s, nilErrorMessage)
		return
	}

	_, _ = io.WriteString(s, e.message+"; details: "+e.log.String())
	if e.cause != nil {
		fmt.Fprintf(s, "; cause: %v", e.cause)
	}

	if verb == 'v' && s.Flag('+') {
		if len(e.stackTrace) > 0 {
			_, _ = io.WriteString(s, "\n ----------------------------------")
		}
		for _, call := range e.stackTrace {
			_, _ = io.WriteString(s, "\n")
			_, _ = io.WriteString(s, call)
		}

		if e.cause != nil {
			_, _ = io.WriteString(s, fmt.Sprintf("\ncause:\n%+v", e.cause))
		}
	}
}

// GetStackTrace ...
func (e *Error) GetStackTrace() string {
	if len(e.stackTrace) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", e.stackTrace)
}
