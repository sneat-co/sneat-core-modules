package dbo4invitus

import (
	"fmt"
	core "github.com/sneat-co/sneat-go-core"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/models/dbprofile"
	"github.com/strongo/validation"
	"net/mail"
	"strings"
	"time"
)

var _ core.Validatable = (*InviteDbo)(nil)

// InviteDbo record - used in PersonalInviteDbo and MassInviteDbo
type InviteDbo struct {
	InviteBase
	Status     InviteStatus         `json:"status" firestore:"status" `
	Pin        string               `json:"pin,omitempty" firestore:"pin,omitempty"`
	SpaceID    coretypes.SpaceID    `json:"spaceID" firestore:"spaceID"`
	TargetType string               `json:"targetType,omitempty" firestore:"targetType,omitempty"`
	TargetIDs  []string             `json:"targetIDs,omitempty" firestore:"targetIDs,omitempty"`
	MessageID  string               `json:"messageID" firestore:"messageID"` // e.g. email message ContactID from AWS SES
	CreatedAt  time.Time            `json:"createdAt" firestore:"createdAt"`
	Created    dbmodels.CreatedInfo `json:"created" firestore:"created"`
	Claimed    *time.Time           `json:"claimed,omitempty" firestore:"claimed,omitempty"`
	Revoked    *time.Time           `json:"revoked" firestore:"revoked,omitempty"`
	Sending    *time.Time           `json:"sending,omitempty" firestore:"sending,omitempty"`
	Sent       *time.Time           `json:"sent,omitempty" firestore:"sent,omitempty"`
	Expires    *time.Time           `json:"expires,omitempty" firestore:"expires,omitempty"`
	Space      *InviteSpace         `json:"space,omitempty" firestore:"space,omitempty"`
	Roles      []string             `json:"roles,omitempty" firestore:"roles,omitempty"`
	FromUserID string               `json:"fromUserID" firestore:"fromUserID"`
	ToUserID   string               `json:"toUserID,omitempty" firestore:"toUserID,omitempty"`
	Message    string               `json:"message,omitempty" firestore:"message,omitempty"`

	// TODO: Document purpose
	Attempts int `json:"attempts,omitempty" firestore:"attempts,omitempty"`

	// Personal invite fields
	Address          string            `json:"address,omitempty" firestore:"address,omitempty"` // Can be empty for a channel=link
	ToSpaceContactID string            `json:"toSpaceContactId" firestore:"toSpaceContactId"`   // in format "<SPACE_ID>:<MEMBER_ID>"
	ToAvatar         *dbprofile.Avatar `json:"toAvatar,omitempty" firestore:"toAvatar,omitempty"`

	// Mass invite fields
	Limit         int `json:"limit,omitempty" firestore:"limit,omitempty"`
	AcceptedCount int `json:"acceptedCount,omitempty" firestore:"acceptedCount,omitempty"`
	DeclinedCount int `json:"declinedCount,omitempty" firestore:"declinedCount,omitempty"`
}

func (v InviteDbo) IsClaimed() bool {
	return v.Claimed != nil || v.Status == InviteStatusAccepted || v.Status == InviteStatusDeclined
}

func (v InviteDbo) validateType(expected InviteType) error {
	if v.Type == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	if v.Type != expected {
		return validation.NewErrBadRecordFieldValue("type", fmt.Sprintf("expected to have value '%s', got: %s", expected, v.Type))
	}
	return nil
}

// Validate validates record
func (v InviteDbo) Validate() error {
	if err := v.InviteBase.Validate(); err != nil {
		return err
	}

	if v.Status == "" {
		return validation.NewErrRecordIsMissingRequiredField("status")
	} else if !IsKnownInviteStatus(v.Status) {
		return validation.NewErrBadRecordFieldValue("status", "unknown value: "+string(v.Status))
	}

	if v.TargetType == "" && len(v.TargetIDs) > 0 {
		return validation.NewErrRecordIsMissingRequiredField("targetType")
	}
	if v.FromUserID == "" {
		return validation.NewErrRecordIsMissingRequiredField("fromUserID")
	}
	if v.From.UserID != v.FromUserID {
		return validation.NewErrBadRecordFieldValue("fromUserID", "does not match from.UserID")
	}
	if len(v.TargetIDs) == 0 && v.TargetType != "" {
		return validation.NewErrRecordIsMissingRequiredField("targetIDs")
	}
	switch v.TargetType {
	case "", "tracker": // known values
	default:
		return validation.NewErrBadRecordFieldValue("targetType", "unknown value: "+v.TargetType)
	}
	for i, targetID := range v.TargetIDs {
		if targetID == "" {
			return validation.NewErrRecordIsMissingRequiredField(fmt.Sprintf("targetIDs[%d]", i))
		}
		if strings.TrimSpace(targetID) != targetID {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("targetIDs[%d]", i), "should not have leading or trailing spaces")
		}
	}

	if v.SpaceID == "" && v.Space != nil {
		return validation.NewErrRecordIsMissingRequiredField("spaceID")
	}

	if err := v.Created.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("created", err.Error())
	}

	if v.Space != nil {
		if err := v.Space.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("space", err.Error())
		}
	}

	if err := v.From.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("from", err.Error())
	}
	if v.To != nil {
		if err := v.To.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("to", err.Error())
		}
	}
	if v.Type == "mass" && len(v.Roles) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("roles")
	}
	if len(v.Roles) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("roles")
	}
	for i, role := range v.Roles {
		if strings.TrimSpace(role) == "" {
			return validation.NewErrRecordIsMissingRequiredField(fmt.Sprintf("roles[%d]", i))
		}
	}
	if err := v.validateType(v.Type); err != nil {
		return err
	}
	if v.Type == InviteTypePersonal && v.ToSpaceContactID == "" {
		return validation.NewErrRecordIsMissingRequiredField("toSpaceContactID")
	}
	if v.ToSpaceContactID != "" {
		if v.ToSpaceContactID[0] == ':' {
			return validation.NewErrBadRecordFieldValue("toSpaceContactID", "starts with ':'")
		}
		if v.ToSpaceContactID[len(v.ToSpaceContactID)-1] == ':' {
			return validation.NewErrBadRecordFieldValue("toSpaceContactID", "ends with ':'")
		}
	}

	switch v.Channel {
	case "email":
		if v.Address != "" {
			if address, err := mail.ParseAddress(v.Address); err != nil {
				return validation.NewErrBadRequestFieldValue("address", fmt.Errorf("field channel is 'email': %w", err).Error())
			} else if address.Name != "" {
				return validation.NewErrBadRecordFieldValue("address", "should not have name, only email address")
			} else if v.Address != strings.ToLower(v.Address) {
				return validation.NewErrBadRecordFieldValue("address", "should be in lower case")
			}
		}
	default:
		if ValidateChannel(v.Channel, true) != nil {
			return validation.NewErrBadRecordFieldValue("channel", "unknown channel value: "+string(v.Channel))
		}
	}
	if !v.ComposeOnly && v.Type == InviteTypePersonal && v.Address == "" {
		return validation.NewErrRecordIsMissingRequiredField("address")
	}
	if v.Pin == "" {
		return validation.NewErrRecordIsMissingRequiredField("pin")
	}
	return nil
}
