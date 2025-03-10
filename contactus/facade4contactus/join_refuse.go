package facade4contactus

import (
	"errors"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// RefuseToJoinSpaceRequest request
type RefuseToJoinSpaceRequest struct {
	SpaceID string `json:"id"`
	Pin     int32  `json:"pin"`
}

// Validate validates request
func (v *RefuseToJoinSpaceRequest) Validate() error {
	if v.SpaceID == "" {
		return validation.NewErrRecordIsMissingRequiredField("space")
	}
	if v.SpaceID == "" {
		return validation.NewErrRecordIsMissingRequiredField("pin")
	}
	return nil
}

// RefuseToJoinSpace refuses to join space
func RefuseToJoinSpace(_ facade.ContextWithUser, request RefuseToJoinSpaceRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	return errors.New("not implemented")
}
