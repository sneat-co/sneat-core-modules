package dto4contactus

import (
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"strings"
)

type CommChannelRequest struct {
	ContactRequest
	ChannelType with.CommChannelType `json:"channelType"`
	ChannelID   string               `json:"channelID"`
}

func (v *CommChannelRequest) Validate() (err error) {
	if err = v.ContactRequest.Validate(); err != nil {
		return
	}
	switch v.ChannelType {
	case with.CommChannelTypeEmail, with.CommChannelTypePhone: // ok
	case "":
		return validation.NewErrRequestIsMissingRequiredField("channelType")
	default:
		return validation.NewErrBadRequestFieldValue("channelType", "unknown value: "+v.ChannelType)
	}
	if s := strings.TrimSpace(v.ChannelID); s == "" {
		return validation.NewErrRequestIsMissingRequiredField("channelID")
	} else if s != v.ChannelID {
		return validation.NewErrBadRecordFieldValue("channelID", "has leading or trailing spaces")
	}
	return nil
}

type AddCommChannelRequest struct {
	CommChannelRequest
	Type string `json:"type"`
	Note string `json:"note,omitempty"`
}

func (v *AddCommChannelRequest) Validate() (err error) {
	if err = v.CommChannelRequest.Validate(); err != nil {
		return
	}
	switch v.ChannelType {
	case with.CommChannelTypeEmail, with.CommChannelTypePhone: // ok
	case "":
		return validation.NewErrRequestIsMissingRequiredField("type")
	default:
		return validation.NewErrBadRequestFieldValue("type", "unknown value: "+v.Type)
	}
	if s := strings.TrimSpace(v.Note); s == "" {
		v.Note = ""
	} else if s != v.Note {
		return validation.NewErrBadRecordFieldValue("note", "has leading or trailing spaces")
	}
	return nil
}

type UpdateCommChannelRequest struct {
	CommChannelRequest
	NewChannelID *string `json:"newChannelID,omitempty"`
	Type         *string `json:"type"`
	Note         *string `json:"note,omitempty"`
}

func (v *UpdateCommChannelRequest) Validate() error {
	if v.Type != nil {
		switch *v.Type {
		case with.CommChannelTypeEmail, with.CommChannelTypePhone: // ok
		case "":
			return validation.NewErrRequestIsMissingRequiredField("type")
		default:
			return validation.NewErrBadRequestFieldValue("type", "unknown value: "+*v.Type)
		}
	}
	if v.Note != nil {
		if s := strings.TrimSpace(*v.Note); s != *v.Note {
			return validation.NewErrBadRecordFieldValue("note", "has leading or trailing spaces")
		}
	}
	return nil
}

type DeleteCommChannelRequest = CommChannelRequest
