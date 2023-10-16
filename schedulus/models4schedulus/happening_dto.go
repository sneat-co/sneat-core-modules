package models4schedulus

import (
	"fmt"
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"strings"
)

// HappeningDto DTO
type HappeningDto struct {
	HappeningBrief
	dbmodels.WithTags
	dbmodels.WithUserIDs
	dbmodels.WithTeamDates
	briefs4contactus.WithMultiTeamContacts[*briefs4contactus.ContactBrief]
}

// Validate returns error if not valid
func (v *HappeningDto) Validate() error {
	if err := v.HappeningBrief.Validate(); err != nil {
		return err
	}
	if err := v.WithUserIDs.Validate(); err != nil {
		return err
	}
	if err := v.WithTeamDates.Validate(); err != nil {
		return err
	}
	if err := v.WithTags.Validate(); err != nil {
		return err
	}
	if len(v.TeamIDs) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("teamIDs")
	}
	for i, level := range v.Levels {
		if l := strings.TrimSpace(level); l == "" {
			return validation.NewErrRecordIsMissingRequiredField(
				fmt.Sprintf("levels[%v]", i),
			)
		} else if l != level {
			return validation.NewErrBadRecordFieldValue(
				fmt.Sprintf("levels[%v]", i),
				fmt.Sprintf("whitespaces at beginning or end: [%v]", level),
			)
		}
	}
	if err := v.WithMultiTeamContactIDs.Validate(); err != nil {
		return err
	}
	switch v.Type {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("type")
	case HappeningTypeSingle:
		if count := len(v.Slots); count > 1 {
			return validation.NewErrBadRecordFieldValue("slots", fmt.Sprintf("single time happening should have only single 'once' slot, got: %v", count))
		}
		if len(v.Dates) == 0 {
			return validation.NewErrRecordIsMissingRequiredField("dates")
		}
		if len(v.TeamDates) == 0 {
			return validation.NewErrRecordIsMissingRequiredField("teamDates")
		}
	case HappeningTypeRecurring:
		if len(v.Dates) > 0 {
			return validation.NewErrBadRequestFieldValue("dates", "should be empty for 'recurring' happening")
		}
	default:
		return validation.NewErrBadRecordFieldValue("type", "unknown value: "+v.Type)
	}
	//if v.Role == HappeningTypeRecurring && v.Status == HappeningStatusCanceled {
	//	for _, slot := range v.Slots {
	//		if slot.Status != SlotStatusCanceled {
	//
	//		}
	//	}
	//}
	return nil
}
