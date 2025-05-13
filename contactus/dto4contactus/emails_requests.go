package dto4contactus

import (
	"github.com/strongo/validation"
	"strings"
)

type EmailRequest struct {
	ContactRequest
	EmailAddress string `json:"emailAddress,omitempty"`
}

func (v EmailRequest) Validate() error {
	if err := v.ContactRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.EmailAddress) == "" {
		return validation.NewErrRequestIsMissingRequiredField("emailAddress")
	}
	return nil
}

type AddEmailRequest struct {
	EmailRequest
	Type  string `json:"type"`
	Title string `json:"title,omitempty"`
}

func (v AddEmailRequest) Validate() error {
	if v.Type == "" {
		return validation.NewErrRequestIsMissingRequiredField("type")
	}
	if v.Title != "" {
		if trimmed := strings.TrimSpace(v.Title); trimmed != v.Title {
			return validation.NewErrBadRequestFieldValue("title", "title have leading or trailing spaces")
		}
	}
	return nil
}

type DeleteEmailRequest = EmailRequest
