package errors

import (
	"google.golang.org/grpc/codes"
)

type ErrorType struct {
	code codes.Code
}

func NewType(errorCode codes.Code) *ErrorType {
	return &ErrorType{code: errorCode}
}

var (
	Internal        = NewType(codes.Internal)
	NotFound        = NewType(codes.NotFound)
	Conflict        = NewType(codes.FailedPrecondition)
	Unauthenticated = NewType(codes.Unauthenticated)
	BadRequest      = NewType(codes.InvalidArgument)
	AccessDenied    = NewType(codes.PermissionDenied)
	Timeout         = NewType(codes.DeadlineExceeded)
	TooManyRequests = NewType(codes.ResourceExhausted)
	NullError       = NewType(codes.Internal)

	clientErrorTypes = [6]*ErrorType{
		NotFound, Conflict, Unauthenticated, BadRequest, AccessDenied, TooManyRequests,
	}
)

// IsTyped ...
func (et *ErrorType) IsTyped(err error) bool {
	if err == nil {
		return false
	}
	if appErr, ok := err.(*Error); ok {
		return appErr.IsTyped(et)
	}
	return false
}

// ErrWrap ...
func (et *ErrorType) ErrWrap(msg string, err error) *Error {
	// prevent double wrap
	if typed, ok := err.(*Error); ok {

		typed = typed.
			WithDetailsKV("wrapped_msg", typed.message).
			WithDetailsKV("wrapped_code", typed.code)

		typed.errorType = et
		typed.message = msg
		return typed
	}

	wrappedErr := et.new(msg, !IsClientErr(err) && et != Internal)
	wrappedErr.cause = err
	return wrappedErr
}

func (et *ErrorType) new(msg string, stacktrace bool) *Error {
	err := &Error{
		errorType:  et,
		message:    msg,
		properties: make(PropertyMap),
	}

	if stacktrace {
		err.stackTrace = getStackTrace()
	}

	return err
}

func Wrap(msg string, err error) error {
	if err == nil || msg == "" {
		return err
	}

	if typed, ok := err.(*Error); ok {
		typed = typed.
			WithDetailsKV("wrapped_msg", typed.message)
		typed.message = msg
		return typed
	}
	return Internal.ErrWrap(msg, err)
}
