package dbo4invitus

import (
	"fmt"
	"github.com/sneat-co/sneat-core-modules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-go-core"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/models/dbprofile"
	"github.com/strongo/validation"
	"net/mail"
	"strings"
	"time"
)

// InviteContact holds invitation contact data
type InviteContact struct {
	Channel   InviteChannel `json:"channel,omitempty" firestore:"channel,omitempty"`
	Address   string        `json:"address,omitempty" firestore:"address,omitempty"`
	Title     string        `json:"title,omitempty" firestore:"title,omitempty"`
	UserID    string        `json:"userID,omitempty" firestore:"userID,omitempty"`
	ContactID string        `json:"contactID,omitempty" firestore:"contactID,omitempty"`
}

func ValidateChannel(v InviteChannel, required bool) error {
	switch v {
	case "":
		if required {
			return validation.NewErrRecordIsMissingRequiredField("channel")
		}
	case "email", "sms", "link", "telegram":
		// known channels
	default:
		return validation.NewErrBadRecordFieldValue("channel", fmt.Sprintf("unsupported value: [%s]", v))
	}
	return nil
}

// Validate returns error if not valid
func (v *InviteContact) Validate() error {
	if err := ValidateChannel(v.Channel, false); err != nil {
		return err
	}
	if v.Channel == "email" && v.Address != "" {
		if _, err := mail.ParseAddress(v.Address); err != nil {
			return validation.NewErrBadRequestFieldValue("address", fmt.Errorf("failed to parse email: %w", err).Error())
		}
	}
	return nil
}

// InviteFrom describes who created the invite
type InviteFrom struct {
	InviteContact
}

// Validate returns error if not valid
func (v InviteFrom) Validate() error {
	if v.UserID == "" {
		return validation.NewErrRecordIsMissingRequiredField("userID")
	}
	if v.ContactID == "" {
		return validation.NewErrRecordIsMissingRequiredField("contactID")
	}
	if err := v.InviteContact.Validate(); err != nil {
		return err
	}
	return nil
}

// InviteTo record
type InviteTo struct {
	InviteContact
}

// Validate returns error if not valid
func (v *InviteTo) Validate() error {
	if err := v.InviteContact.Validate(); err != nil {
		return err
	}
	if v.Channel == "" {
		return validation.NewErrRecordIsMissingRequiredField("channel")
	}
	if v.Channel == "email" {
		if strings.TrimSpace(v.Address) == "" {
			return validation.NewErrRecordIsMissingRequiredField("address")
		}
		if _, err := mail.ParseAddress(v.Address); err != nil {
			return validation.NewErrBadRecordFieldValue("address", "not a valid email")
		}
	}
	const maxTitleLen = 100
	if len(v.Title) > maxTitleLen {
		return validation.NewErrBadRecordFieldValue("title",
			fmt.Sprintf("contact title should not exceed max length of %d, got: %d",
				maxTitleLen, len(v.Title)))
	}
	//if strings.TrimSpace(v.Title) == "" {
	//	return validation.NewErrRecordIsMissingRequiredField("Title")
	//}
	return nil
}

// InviteSpace a summary on space for which an invitation has been created
type InviteSpace struct {
	ID    string                 `json:"id,omitempty" firestore:"id,omitempty"`
	Type  core4spaceus.SpaceType `json:"type" firestore:"type"`
	Title string                 `json:"title,omitempty" firestore:"title,omitempty"`
}

// Validate returns error if not valid
func (v InviteSpace) Validate() error {
	//if v.InviteID == "" {
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

type InviteChannel string

type InviteType string

const (
	InviteTypePersonal InviteType = "personal" // To a specific person
	InviteTypePrivate  InviteType = "private"  // To a single person
	InviteTypeMass     InviteType = "mass"     // To a group of people
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
	if err := ValidateChannel(v.Channel, true); err != nil {
		return err
	}
	return nil
}

// InviteBrief summary about invite
type InviteBrief struct {
	ID   string      `json:"id" firestore:"id"`
	Pin  string      `json:"pin,omitempty" firestore:"pin,omitempty"`
	From *InviteFrom `json:"from,omitempty" firestore:"from,omitempty"`
	To   *InviteTo   `json:"to,omitempty" firestore:"to,omitempty"`
}

// Validate returns error if not valid
func (v InviteBrief) Validate() error {
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	if err := v.From.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("from", err.Error())
	}
	if v.To != nil {
		if err := v.To.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("to", err.Error())
		}
	}
	return nil
}

// NewInviteBriefFromDbo creates brief from DTO
func NewInviteBriefFromDbo(id string, dto InviteDbo) InviteBrief {
	from := dto.From
	to := *dto.To
	return InviteBrief{ID: id, From: &from, To: &to, Pin: dto.Pin}
}

type InviteStatus string

const (
	InviteStatusPending  InviteStatus = "pending"
	InviteStatusActive   InviteStatus = "active"
	InviteStatusAccepted InviteStatus = "accepted"
	InviteStatusDeclined InviteStatus = "declined"
	InviteStatusExpired  InviteStatus = "expired"
)

// InviteDbo record - used in PersonalInviteDbo and MassInviteDbo
type InviteDbo struct {
	InviteBase
	Status    InviteStatus         `json:"status" firestore:"status" `
	Pin       string               `json:"pin,omitempty" firestore:"pin,omitempty"`
	SpaceID   string               `json:"spaceID" firestore:"spaceID"`
	MessageID string               `json:"messageId" firestore:"messageId"` // e.g. email message ContactID from AWS SES
	CreatedAt time.Time            `json:"createdAt" firestore:"createdAt"`
	Created   dbmodels.CreatedInfo `json:"created" firestore:"created"`
	Claimed   *time.Time           `json:"claimed,omitempty" firestore:"claimed,omitempty"`
	Revoked   *time.Time           `json:"revoked" firestore:"revoked,omitempty"`
	Sending   *time.Time           `json:"sending,omitempty" firestore:"sending,omitempty"`
	Sent      *time.Time           `json:"sent,omitempty" firestore:"sent,omitempty"`
	Expires   *time.Time           `json:"expires,omitempty" firestore:"expires,omitempty"`
	Space     InviteSpace          `json:"space" firestore:"space"`
	Roles     []string             `json:"roles,omitempty" firestore:"roles,omitempty"`
	//FromUserID string     `json:"fromUserID" firestore:"fromUserID"`
	//ToUserID   string     `json:"toUserID,omitempty" firestore:"toUserID,omitempty"`
	Message string `json:"message,omitempty" firestore:"message,omitempty"`

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

// Validate validates record
func (v InviteDbo) Validate() error {
	if err := v.InviteBase.Validate(); err != nil {
		return err
	}
	switch v.Status {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("status")
	case InviteStatusPending, InviteStatusActive, InviteStatusAccepted, InviteStatusDeclined, InviteStatusExpired: // known statuses
	default:
		return validation.NewErrBadRecordFieldValue("status", "unknown value: "+string(v.Status))
	}
	if v.SpaceID == "" {
		return validation.NewErrRecordIsMissingRequiredField("spaceID")
	}
	if err := v.Created.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("created", err.Error())
	}
	if err := v.Space.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("space", err.Error())
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

func (v InviteDbo) validateType(expected InviteType) error {
	if v.Type == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	if v.Type != expected {
		return validation.NewErrBadRecordFieldValue("type", fmt.Sprintf("expected to have value '%s', got: %s", expected, v.Type))
	}
	return nil
}

var _ core.Validatable = (*InviteDbo)(nil)

// InviteClaim record
type InviteClaim struct {
	Time   time.Time `json:"time" firestore:"time"`
	UserID string    `json:"userId" firestore:"userId"`
}

// InviteCode record
type InviteCode struct {
	Claims []InviteClaim `json:"claims,omitempty" firestore:"claims,omitempty"`
}
