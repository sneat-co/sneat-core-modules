package models4calendarium

import (
	"fmt"
	"github.com/strongo/validation"
)

type CalendariumTeamDto struct {
	RecurringHappenings map[string]*HappeningBrief `json:"recurringHappenings,omitempty" firestore:"recurringHappenings,omitempty"`
}

func (v *CalendariumTeamDto) Validate() error {
	for i, rh := range v.RecurringHappenings {
		if err := rh.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("recurringHappenings", fmt.Errorf("invalid value at index %v: %w", i, err).Error())
		}
		if rh.Type != HappeningTypeRecurring {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("recurringHappenings[%v].type", i),
				fmt.Errorf("expected to have value 'recurring' got '%v'", rh.Type).Error())
		}
	}
	return nil
}

func (v *CalendariumTeamDto) GetRecurringHappeningBrief(id string) *HappeningBrief {
	return v.RecurringHappenings[id]
}
