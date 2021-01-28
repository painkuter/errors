package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPairs(t *testing.T) {
	k := int64(1)
	assert.Equal(t, fmt.Sprintf("%d", k), "1")
}
