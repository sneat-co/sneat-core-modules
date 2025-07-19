package auth

import (
	"testing"
)

func TestModule(t *testing.T) {
	m := Extension()
	if m == nil {
		t.Fatal("ExtID() returned nil")
	}
}
