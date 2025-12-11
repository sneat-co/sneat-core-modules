package common4all

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeID(t *testing.T) {
	_, err := DecodeID("")
	assert.NotNil(t, err) // Should return error if empty string
}
