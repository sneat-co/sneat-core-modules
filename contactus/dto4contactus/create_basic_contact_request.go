package dto4contactus

import (
	"strings"

	"github.com/strongo/validation"
)

// CreateBasicContactRequest - creates a basic contact
type CreateBasicContactRequest struct {
	Title string `json:"title"`
}

// Validate returns error if not valid
func (v CreateBasicContactRequest) Validate() error {
	if strings.TrimSpace(v.Title) == "" {
		return validation.NewErrRequestIsMissingRequiredField("title")
	}
	return nil
}
