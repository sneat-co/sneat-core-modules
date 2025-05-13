package dbo4all

import (
	"fmt"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"strings"
)

type WithPhones struct {
	Phones map[string]PhoneProps `json:"phones,omitempty" firestore:"phones,omitempty"`
}

func (v WithPhones) Validate() error {
	for k, p := range v.Phones {
		if strings.TrimSpace(k) == "" {
			return fmt.Errorf("phone key is empty")
		}
		if trimmedKey := strings.TrimSpace(k); trimmedKey == "" {
			return validation.NewErrBadRecordFieldValue("phones."+k, "phone key is empty")
		} else if k != trimmedKey {
			return validation.NewErrBadRecordFieldValue("phones."+k, "phone key has leading or trailing spaces")
		}
		if err := p.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("phones."+k, err.Error())
		}
	}
	return nil
}

type PhoneProps struct {
	with.CreatedFields
	with.TagsField
	Type  string `json:"type" firestore:"type"`
	Title string `json:"title,omitempty" firestore:"title,omitempty"`
}

func (v PhoneProps) Validate() error {
	if err := v.CreatedFields.Validate(); err != nil {
		return err
	}
	if err := v.TagsField.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.Type) == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	if v.Title != "" {
		if trimmed := strings.TrimSpace(v.Title); trimmed == "" {
			if v.Title != "" {
				return fmt.Errorf("title has spaces but is empty")
			}
		} else if v.Title != trimmed {
			return fmt.Errorf("title has leading or trailing spaces")
		}
	}
	return nil
}
