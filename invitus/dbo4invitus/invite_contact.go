package dbo4invitus

import (
	"fmt"
	"github.com/strongo/validation"
	"net/mail"
	"strings"
)

// InviteContact holds invitation contact data
type InviteContact struct {
	Channel   InviteChannel `json:"channel,omitempty" firestore:"channel,omitempty"`
	Address   string        `json:"address,omitempty" firestore:"address,omitempty"`
	Title     string        `json:"title,omitempty" firestore:"title,omitempty"`
	UserID    string        `json:"userID,omitempty" firestore:"userID,omitempty"`
	ContactID string        `json:"contactID,omitempty" firestore:"contactID,omitempty"`
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
	//if v.ContactID == "" {
	//	return validation.NewErrRecordIsMissingRequiredField("contactID")
	//}
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
