package dbo4linkage

import (
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
)

type ShortSpaceModuleItemRef struct {
	ID      string            `json:"id" firestore:"id"`
	SpaceID coretypes.SpaceID `json:"spaceID,omitempty" firestore:"spaceID,omitempty"`
}

func (v *ShortSpaceModuleItemRef) Validate() error {
	// SpaceID can be empty for global collections like Happening
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("itemID")
	} else if err := validate.RecordID(v.ID); err != nil {
		return validation.NewErrBadRecordFieldValue("itemID", err.Error())
	}
	return nil
}
