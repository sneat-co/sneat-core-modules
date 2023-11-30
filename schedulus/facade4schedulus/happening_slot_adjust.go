package facade4schedulus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/schedulus/dto4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/models4schedulus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
)

func AdjustSlot(ctx context.Context, user facade.User, request dto4schedulus.HappeningSlotDateRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	var worker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *happeningWorkerParams) (err error) {
		switch params.Happening.Dto.Type {
		case "single":
			return errors.New("only recurring happenings can be adjusted, single happenings should be updated")
		case "recurring":
			if err = adjustRecurringSlot(ctx, tx, params.Happening, request); err != nil {
				return fmt.Errorf("failed to adjust recurring happening: %w", err)
			}
			return err
		}
		return
	}

	if err = modifyHappening(ctx, user, request.HappeningRequest, worker); err != nil {
		return err
	}
	return nil
}

func adjustRecurringSlot(ctx context.Context, tx dal.ReadwriteTransaction, happening models4schedulus.HappeningContext, request dto4schedulus.HappeningSlotDateRequest) (err error) {
	//for _, teamID := range happening.Dto.TeamIDs { // TODO: run in parallel in go routine if > 1
	if err := adjustSlotInScheduleDay(ctx, tx, request.TeamID, happening.ID, request); err != nil {
		return fmt.Errorf("failed to adjust slot in schedule day record for teamID=%v: %w", request.TeamID, err)
	}
	//}
	return nil
}

func adjustSlotInScheduleDay(ctx context.Context, tx dal.ReadwriteTransaction, teamID, happeningID string, request dto4schedulus.HappeningSlotDateRequest) error {
	scheduleDay := models4schedulus.NewScheduleDayContext(teamID, request.Date)
	if err := tx.Get(ctx, scheduleDay.Record); err != nil {
		if !dal.IsNotFound(err) {
			return fmt.Errorf("failed to get schedule day record: %w", err)
		}
	}
	_, adjustment := scheduleDay.Dto.GetAdjustment(happeningID, request.Slot.ID)
	if adjustment == nil {
		adjustment = &models4schedulus.HappeningAdjustment{
			HappeningID: happeningID,
		}
		scheduleDay.Dto.HappeningAdjustments = append(scheduleDay.Dto.HappeningAdjustments, adjustment)
	}
	adjustment.Slot = request.Slot
	var happeningIDsChanged bool
	if happeningIDsChanged = slice.Index(scheduleDay.Dto.HappeningIDs, happeningID) < 0; happeningIDsChanged {
		scheduleDay.Dto.HappeningIDs = append(scheduleDay.Dto.HappeningIDs, happeningID)
	}

	if err := scheduleDay.Dto.Validate(); err != nil {
		return fmt.Errorf("schedule day record is not valid: %w", err)
	}

	if scheduleDay.Record.Exists() {
		updates := []dal.Update{
			{Field: "happeningAdjustments", Value: scheduleDay.Dto.HappeningAdjustments},
		}
		if happeningIDsChanged {
			updates = append(updates, dal.Update{
				Field: "happeningIDs", Value: scheduleDay.Dto.HappeningIDs,
			})
		}
		if err := tx.Update(ctx, scheduleDay.Key, updates); err != nil {
			return fmt.Errorf("failed to update schedule day record with happening adjustment: %w", err)
		}
	} else {
		if err := tx.Insert(ctx, scheduleDay.Record); err != nil {
			return fmt.Errorf("failed to insert schedule day record with happening adjustment: %w", err)
		}
	}
	return nil
}
