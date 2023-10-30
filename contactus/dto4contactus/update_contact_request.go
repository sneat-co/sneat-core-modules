package dto4contactus

import (
	"errors"
	"github.com/sneat-co/sneat-core-modules/linkage/models4linkage"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"strings"
)

type UpdateContactRequest struct {
	ContactRequest
	Address   *dbmodels.Address       `json:"address,omitempty"`
	AgeGroup  string                  `json:"ageGroup,omitempty"`
	Roles     *SetContactRolesRequest `json:"roles,omitempty"`
	VatNumber *string                 `json:"vatNumber,omitempty"`
	models4linkage.WithRelated
}

func (v UpdateContactRequest) Validate() error {
	if err := v.ContactRequest.Validate(); err != nil {
		return err
	}
	if v.Address == nil && v.AgeGroup == "" && v.Roles == nil && v.Related == nil && v.VatNumber == nil {
		return validation.NewBadRequestError(errors.New("at least one of contact fields must be provided for an update"))
	}
	if v.Address != nil {
		if err := v.Address.Validate(); err != nil {
			return validation.NewErrBadRequestFieldValue("address", err.Error())
		}
	}
	if err := dbmodels.ValidateAgeGroup(v.AgeGroup, false); err != nil {
		return validation.NewErrBadRequestFieldValue("ageGroup", err.Error())
	}
	if v.Roles != nil {
		if err := v.Roles.Validate(); err != nil {
			return validation.NewErrBadRequestFieldValue("roles", err.Error())
		}
	}
	if v.VatNumber != nil {
		vat := *v.VatNumber
		if strings.TrimSpace(vat) == vat {
			return validation.NewErrBadRequestFieldValue("vatNumber", "must not have leading or trailing spaces")
		}

	}
	if err := v.WithRelated.Validate(); err != nil {
		return validation.NewBadRequestError(err)
	}
	return nil
}
