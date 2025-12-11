package facade4spaceus

import (
	"testing"

	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
)

func TestCreateSpaceRequest_Validate(t *testing.T) {
	request := dto4spaceus.CreateSpaceRequest{Title: ""}
	if request.Validate() == nil {
		t.Error("request.Validate() == nil")
	}
}
