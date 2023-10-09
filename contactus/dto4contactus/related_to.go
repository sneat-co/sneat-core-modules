package dto4contactus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
)

type RelatedToRequest struct {
	ModuleID   string `json:"moduleID"`
	Collection string `json:"collection"`
	ItemID     string `json:"itemID,omitempty"` // if empty use current user ID
	RelatedAs  string `json:"relatedAs"`
}

func (v RelatedToRequest) Validate() error {
	if v.ModuleID == "" {
		return validation.NewErrRecordIsMissingRequiredField("moduleID")
	}
	if v.Collection == "" {
		return validation.NewErrRecordIsMissingRequiredField("collection")
	}
	if v.ItemID == "" {
		if !(v.ModuleID == const4contactus.ModuleID && v.Collection == const4contactus.ContactsCollection) {
			return validation.NewErrRecordIsMissingRequiredField("itemID")
		}
		// OK to be empty, will use contact ID of current user
	} else if err := validate.RecordID(v.ItemID); err != nil {
		return validation.NewErrBadRecordFieldValue("itemID", err.Error())
	}
	if v.RelatedAs == "" {
		return validation.NewErrRequestIsMissingRequiredField("relatedAs")
	}
	return nil
}
