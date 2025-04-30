package dbo4linkage

import (
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/strongo/validation"
)

type SpaceItemRef struct {
	SpaceID coretypes.SpaceID `json:"spaceID,omitempty" firestore:"spaceID,omitempty"`
	ItemRef
}

func (v *SpaceItemRef) Validate() error {
	if v.SpaceID == "" {
		return validation.NewErrRecordIsMissingRequiredField("spaceID")
	}
	if err := v.ItemRef.Validate(); err != nil {
		return err
	}
	return nil
}
