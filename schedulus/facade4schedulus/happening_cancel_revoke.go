package facade4schedulus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/schedulus/const4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/dto4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/models4schedulus"
	"github.com/sneat-co/sneat-core-modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
	"log"
)

// RevokeHappeningCancellation marks happening as canceled
func RevokeHappeningCancellation(ctx context.Context, user facade.User, request dto4schedulus.CancelHappeningRequest) (err error) {
	log.Printf("RevokeHappeningCancellation() %+v", request)
	if err = request.Validate(); err != nil {
		return err
	}

	happening := models4schedulus.NewHappeningContext(request.TeamID, request.HappeningID)
	err = dal4teamus.RunModuleTeamWorker(ctx, user, request.TeamRequest,
		const4schedulus.ModuleID,
		new(models4schedulus.SchedulusTeamDto),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.ModuleTeamWorkerParams[*models4schedulus.SchedulusTeamDto]) (err error) {
			if err = tx.Get(ctx, happening.Record); err != nil {
				return fmt.Errorf("failed to get happening: %w", err)
			}
			switch happening.Dto.Type {
			case "":
				return fmt.Errorf("happening record has no type: %w", validation.NewErrRecordIsMissingRequiredField("type"))
			case "single":
				return revokeSingleHappeningCancellation(ctx, tx, happening)
			case "recurring":
				return revokeRecurringHappeningCancellation(ctx, tx, params, happening, request.Date, request.SlotID)
			default:
				return validation.NewErrBadRecordFieldValue("type", "happening has unknown type: "+happening.Dto.Type)
			}
		})
	if err != nil {
		return fmt.Errorf("failed to cancel happening: %w", err)
	}
	return
}

func revokeSingleHappeningCancellation(ctx context.Context, tx dal.ReadwriteTransaction, happening models4schedulus.HappeningContext) error {
	return removeCancellationFromHappeningRecord(ctx, tx, happening)
}

func revokeRecurringHappeningCancellation(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *dal4teamus.ModuleTeamWorkerParams[*models4schedulus.SchedulusTeamDto],
	happening models4schedulus.HappeningContext,
	dateID string,
	slotID string,
) error {
	log.Printf("revokeRecurringHappeningCancellation(): teamID=%v, dateID=%v, happeningID=%v, slotID=%+v", params.Team.ID, dateID, happening.ID, slotID)
	if happening.Dto.Status == models4schedulus.HappeningStatusCanceled {
		if err := removeCancellationFromHappeningRecord(ctx, tx, happening); err != nil {
			return fmt.Errorf("failed to remove cancellation from happening record: %w", err)
		}
	}
	if dateID == "" {
		if err := removeCancellationFromHappeningBrief(params, happening); err != nil {
			return fmt.Errorf("failed to remove cancellation from happening brief in team record: %w", err)
		}
	} else if err := removeCancellationFromScheduleDay(ctx, tx, params.Team.ID, dateID, happening.ID, slotID); err != nil {
		return fmt.Errorf("failed to remove cancellation from schedule day record: %w", err)
	}
	return nil
}

func removeCancellationFromHappeningBrief(params *dal4teamus.ModuleTeamWorkerParams[*models4schedulus.SchedulusTeamDto], happening models4schedulus.HappeningContext) error {
	happeningBrief := params.TeamModuleEntry.Data.GetRecurringHappeningBrief(happening.ID)
	if happeningBrief == nil {
		return nil
	}
	if happeningBrief.Status == models4schedulus.HappeningStatusCanceled {
		happeningBrief.Status = models4schedulus.HappeningStatusActive
		happeningBrief.Canceled = nil
		if err := happeningBrief.Validate(); err != nil {
			return err
		}
		params.TeamUpdates = append(params.TeamUpdates, dal.Update{
			Field: "recurringHappenings",
			Value: params.TeamModuleEntry.Data.RecurringHappenings,
		})
	}
	return nil
}

func removeCancellationFromHappeningRecord(ctx context.Context, tx dal.ReadwriteTransaction, happening models4schedulus.HappeningContext) error {
	if happening.Dto.Status != models4schedulus.HappeningStatusCanceled {
		return fmt.Errorf("not allowed to revoke cancelation for happening in status=" + happening.Dto.Status)
	}
	happening.Dto.Status = models4schedulus.HappeningStatusCanceled
	happening.Dto.Canceled = nil
	if err := happening.Dto.Validate(); err != nil {
		return err
	}
	updates := []dal.Update{
		{Field: "status", Value: models4schedulus.HappeningStatusActive},
		{Field: "canceled", Value: dal.DeleteField},
	}
	if err := tx.Update(ctx, happening.Key, updates); err != nil {
		return fmt.Errorf("failed to update happening record: %w", err)
	}
	return nil

}

func removeCancellationFromScheduleDay(ctx context.Context, tx dal.ReadwriteTransaction, teamID, dateID, happeningID string, slotID string) error {
	log.Printf("removeCancellationFromScheduleDay(): teamID=%v, dateID=%v, happeningID=%v, slotID=%+v", teamID, dateID, happeningID, slotID)
	if len(slotID) == 0 {
		return validation.NewErrRequestIsMissingRequiredField("slotID")
	}
	scheduleDay := models4schedulus.NewScheduleDayContext(teamID, dateID)
	if err := tx.Get(ctx, scheduleDay.Record); err != nil {
		if dal.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("failed to get schedule day record by ContactID")
	}
	if i, adjustment := scheduleDay.Dto.GetAdjustment(happeningID, slotID); adjustment != nil && adjustment.Canceled != nil {
		a := scheduleDay.Dto.HappeningAdjustments
		scheduleDay.Dto.HappeningAdjustments = append(a[:i], a[i+1:]...)
		if len(scheduleDay.Dto.HappeningAdjustments) == 0 {
			if err := tx.Delete(ctx, scheduleDay.Key); err != nil {
				return fmt.Errorf("failed to delete schedule day record: %w", err)
			}
		}
	}
	return nil
}
