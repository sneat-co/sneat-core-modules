package dto4contactus

import (
	"errors"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/validation"
	"strings"
	"time"
)

type UpdateContactRequest struct {
	ContactRequest
	Address     *dbmodels.Address  `json:"address,omitempty"`
	AgeGroup    string             `json:"ageGroup,omitempty"`
	Roles       *SetRolesRequest   `json:"roles,omitempty"`
	VatNumber   *string            `json:"vatNumber,omitempty"`
	Gender      dbmodels.Gender    `json:"gender,omitempty"`
	DateOfBirth *string            `json:"dateOfBirth,omitempty"`
	Names       *person.NameFields `json:"names,omitempty"`
}

func (v UpdateContactRequest) Validate() error {
	if err := v.ContactRequest.Validate(); err != nil {
		return err
	}
	if v.Address == nil && v.AgeGroup == "" && v.Roles == nil && v.VatNumber == nil && v.Gender == "" && v.Names == nil && v.DateOfBirth == nil {
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
	if v.Names != nil {
		if err := v.Names.Validate(); err != nil {
			return validation.NewErrBadRequestFieldValue("names", err.Error())
		}
	}
	if v.Gender != "" && !dbmodels.IsKnownGender(v.Gender) {
		return validation.NewErrBadRequestFieldValue("gender", "unknown value: "+v.Gender)
	}
	if v.DateOfBirth != nil {
		if dob := *v.DateOfBirth; dob != "" {
			if trimmed := strings.TrimSpace(dob); trimmed != dob {
				return validation.NewErrBadRequestFieldValue("dateOfBirth", "must not have leading or trailing spaces")
			}
			if _, err := time.Parse(time.DateOnly, dob); err != nil {
				return validation.NewErrBadRequestFieldValue("dateOfBirth", err.Error())
			}
		}
	}
	return nil
}
