package dbo4invitus

import (
	"fmt"
	"github.com/strongo/validation"
)

// InviteBase base data about invite to be used in InviteBrief & InviteDbo
type InviteBase struct {
	Type        InviteType    `json:"type" firestore:"type"` // either "personal" or "mass"
	Channel     InviteChannel `json:"channel" firestore:"channel"`
	ComposeOnly bool          `json:"composeOnly" firestore:"composeOnly"`
	From        InviteFrom    `json:"from" firestore:"from"`
	To          *InviteTo     `json:"to" firestore:"to"`
}

// Validate returns error if not valid
func (v InviteBase) Validate() error {
	switch v.Type {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("type")
	case InviteTypePrivate:
		if v.To == nil {
			return fmt.Errorf("%w: expected to be either 'personal' or 'mass'", validation.NewErrRecordIsMissingRequiredField("to"))
		}
	case InviteTypePersonal:
		if v.To == nil {
			return fmt.Errorf("%w: expected to be either 'personal' or 'mass'", validation.NewErrRecordIsMissingRequiredField("to"))
		}
		if err := v.To.Validate(); err != nil {
			return err
		}
		// known values
	case InviteTypeMass:
		if v.To != nil {
			// TODO: we might want to change this to store a distribution channel?
			return validation.NewErrBadRecordFieldValue("to", "mass invite can not have 'to' value for now")
		}
		// known
	default:
		return validation.NewErrBadRecordFieldValue("type", "unknown invite type: "+string(v.Type))
	}
	if err := v.From.Validate(); err != nil {
		return err
	}
	if err := ValidateChannel(v.Channel, true); err != nil {
		return err
	}
	return nil
}
