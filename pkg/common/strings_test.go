package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPointerToString(t *testing.T) {
	assert.Equal(t, PointerToString(nil), "")
}
