package dbo4invitus

import (
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/strongo/validation"
)

// InviteSpace a summary on space for which an invitation has been created
type InviteSpace struct {
	//ID    coretypes.SpaceID   `json:"id,omitempty" firestore:"id,omitempty"`
	Type  coretypes.SpaceType `json:"type" firestore:"type"`
	Title string              `json:"title,omitempty" firestore:"title,omitempty"`
}

// Validate returns error if not valid
func (v InviteSpace) Validate() error {
	//if v.InviteID == "" {
	//	return validation.NewErrRecordIsMissingRequiredField("id")
	//}
	//if v.ID == "" {
	//	return validation.NewErrRecordIsMissingRequiredField("id")
	//}
	if v.Type == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	switch v.Type {
	case "family":
		// Can be empty
	default:
		if v.Title == "" {
			return validation.NewErrRecordIsMissingRequiredField("title")
		}
	}
	return nil
}
