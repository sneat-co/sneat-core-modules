package dto4spaceus

import (
	"testing"
)

func TestValidSpaceRequest(t *testing.T) {
	if err := ValidSpaceRequest().Validate(); err != nil {
		t.Errorf("ValidSpaceRequest().Validate() returned error: %s", err)
	}
}
