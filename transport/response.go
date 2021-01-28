package transport

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/status"

	"github.com/painkuter/errors"
	_ "github.com/painkuter/errors"
	"github.com/painkuter/errors/protobuf"
)

// Error ...
type Error struct {
	cause *errors.Error
}

// NewError ...
func NewError(err *errors.Error) *Error {
	return &Error{cause: err}
}

// Error ...
func (e *Error) Error() string {
	return e.cause.Error()
}

// Cause ...
func (e *Error) Cause() error {
	return e.cause
}

// GRPCStatus ...
func (e *Error) GRPCStatus() *status.Status {
	st := status.New(e.cause.GRPCCode(), e.cause.Error())
	stWithDetails, err := st.WithDetails(e.ProtoMessage())
	if err != nil {
		return st
	}
	return stWithDetails
}

// Details is implementation of platform/errors.errDetails interface
func (e *Error) Details() interface{} {
	return []proto.Message{e.ProtoMessage()}
}

// ProtoMessage get error in form of proto message
func (e *Error) ProtoMessage() *protobuf.Error {
	kv := e.cause.Details()
	data := make([]*protobuf.Error_Detail, 0, len(kv)/2)
	for i := 1; i < len(kv); i += 2 {
		entry := &protobuf.Error_Detail{
			Key:   fmt.Sprint(kv[i-1]),
			Value: fmt.Sprint(kv[i]),
		}

		data = append(data, entry)
	}

	return &protobuf.Error{
		Code:    e.cause.Code(),
		Message: e.Error(),
		Details: data,
	}
}

// ProtoResponse ...
func (e *Error) ProtoResponse() interface{} {
	return &protobuf.ErrorResponse{Error: e.ProtoMessage()}
}
