package dbo4invitus

import (
	"fmt"
	"github.com/strongo/validation"
)

type InviteChannel string

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
