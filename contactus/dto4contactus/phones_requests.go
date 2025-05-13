package dto4contactus

import (
	"github.com/strongo/validation"
)

type PhoneRequest struct {
	ContactRequest
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

func (v PhoneRequest) Validate() error {
	if err := v.ContactRequest.Validate(); err != nil {
		return err
	}
	if v.PhoneNumber == "" {
		return validation.NewErrRequestIsMissingRequiredField("phoneNumber")
	}
	return nil
}

type AddPhoneRequest struct {
	PhoneRequest
	Type string `json:"type"`
}

func (v AddPhoneRequest) Validate() error {
	if err := v.PhoneRequest.Validate(); err != nil {
		return err
	}
	if v.Type == "" {
		return validation.NewErrRequestIsMissingRequiredField("type")
	}
	return nil
}

type DeletePhoneRequest = PhoneRequest
