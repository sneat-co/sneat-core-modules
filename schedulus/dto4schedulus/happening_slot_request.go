package dto4schedulus

import (
	"github.com/sneat-co/sneat-core-modules/schedulus/models4schedulus"
)

// HappeningSlotRequest updates slot
type HappeningSlotRequest struct {
	HappeningRequest
	Slot models4schedulus.HappeningSlot `json:"slot"`
}

// Validate returns error if not valid
func (v HappeningSlotRequest) Validate() error {
	if err := v.HappeningRequest.Validate(); err != nil {
		return err
	}
	if err := v.Slot.Validate(); err != nil {
		return err
	}
	return nil
}
