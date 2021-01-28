package errors_test

import (
	"errors"
	"testing"

	pkgerr "github.com/pkg/errors"
)

func TestErrorWrap(t *testing.T) {
	pkgerr.Wrap(errors.New("test"), "wrap message")
}
