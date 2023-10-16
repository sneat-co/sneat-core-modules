package facade4schedulus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/schedulus/const4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/dto4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/models4schedulus"
	"github.com/sneat-co/sneat-core-modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
	"strings"
	"time"
)

// CancelHappening marks happening as canceled
func CancelHappening(ctx context.Context, user facade.User, request dto4schedulus.CancelHappeningRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	happening := models4schedulus.NewHappeningContext(request.HappeningID)
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
				return cancelSingleHappening(ctx, tx, params.UserID, happening)
			case "recurring":
				return cancelRecurringHappening(ctx, tx, params, params.UserID, happening, request)
			default:
				return validation.NewErrBadRecordFieldValue("type", "happening has unknown type: "+happening.Dto.Type)
			}
		})
	if err != nil {
		return fmt.Errorf("failed to cancel happening: %w", err)
	}
	return
}

func cancelSingleHappening(ctx context.Context, tx dal.ReadwriteTransaction, userID string, happening models4schedulus.HappeningContext) error {
	switch happening.Dto.Status {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("status")
	case models4schedulus.HappeningStatusActive:
		happening.Dto.Status = models4schedulus.HappeningStatusCanceled
		happening.Dto.Canceled = &models4schedulus.Canceled{
			At: time.Now(),
			By: dbmodels.ByUser{UID: userID},
		}
		if err := happening.Dto.Validate(); err != nil {
			return fmt.Errorf("happening record is not valid: %w", err)
		}
		happeningUpdates := []dal.Update{
			{Field: "status", Value: happening.Dto.Status},
			{Field: "canceled", Value: happening.Dto.Canceled},
		}
		if err := tx.Update(ctx, happening.Key, happeningUpdates); err != nil {
			return err
		}
	case models4schedulus.HappeningStatusDeleted:
		// Nothing to do
	default:
		return fmt.Errorf("only active happening can be canceled but happening is in status=[%v]", happening.Dto.Status)
	}
	happening.Dto.Status = "canceled"
	return nil
}

func cancelRecurringHappening(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *dal4teamus.ModuleTeamWorkerParams[*models4schedulus.SchedulusTeamDto],
	uid string,
	happening models4schedulus.HappeningContext,
	request dto4schedulus.CancelHappeningRequest,
) error {
	happeningBrief := params.TeamModuleEntry.Data.GetRecurringHappeningBrief(happening.ID)
	if happeningBrief == nil {
		return errors.New("happening brief is not found in team record")
	}

	if request.Date == "" {
		if err := markRecurringHappeningRecordAsCanceled(ctx, tx, uid, happening, request); err != nil {
			return err
		}
		happeningBrief.Status = models4schedulus.HappeningStatusCanceled
		happeningBrief.Canceled = createCanceled(uid, request.Reason)
		if err := happeningBrief.Validate(); err != nil {
			return fmt.Errorf("happening brief in team record is not valid: %w", err)
		}
		params.TeamUpdates = append(params.TeamUpdates, dal.Update{
			Field: "recurringHappenings",
			Value: params.TeamModuleEntry.Data.RecurringHappenings,
		})
	} else {
		scheduleDay := models4schedulus.NewScheduleDayContext(params.Team.ID, request.Date)

		var isNewRecord bool
		if err := tx.Get(ctx, scheduleDay.Record); err != nil {
			if dal.IsNotFound(err) {
				isNewRecord = true
			} else {
				return fmt.Errorf("failed to get schedule day record by ContactID: %w", err)
			}
		}

		var dayUpdates []dal.Update
		_, adjustment := scheduleDay.Dto.GetAdjustment(happening.ID, request.SlotID)
		if adjustment == nil {
			_, slot := happening.Dto.GetSlot(request.SlotID)
			if slot == nil {
				return fmt.Errorf("%w: slot not found by ContactID=%v", facade.ErrBadRequest, request.SlotID)
			}
			adjustment = &models4schedulus.HappeningAdjustment{
				HappeningID: happening.ID,
				Slot:        *slot,
				Canceled: &models4schedulus.Canceled{
					At:     time.Now(),
					By:     dbmodels.ByUser{UID: uid},
					Reason: request.Reason,
				},
			}
			scheduleDay.Dto.HappeningAdjustments = append(scheduleDay.Dto.HappeningAdjustments, adjustment)
		}
		if i := slice.Index(scheduleDay.Dto.HappeningIDs, happening.ID); i < 0 {
			scheduleDay.Dto.HappeningIDs = append(scheduleDay.Dto.HappeningIDs, happening.ID)
			dayUpdates = append(dayUpdates, dal.Update{
				Field: "happeningIDs", Value: scheduleDay.Dto.HappeningIDs,
			})
		}
		var modified bool
		if adjustment.Slot.ID == request.SlotID {
			if strings.TrimSpace(request.Reason) != "" && (adjustment.Canceled == nil || request.Reason != adjustment.Canceled.Reason) {
				if adjustment.Canceled == nil {
					adjustment.Canceled = &models4schedulus.Canceled{
						At:     time.Now(),
						By:     dbmodels.ByUser{UID: uid},
						Reason: request.Reason,
					}
				}
				adjustment.Canceled.Reason = request.Reason
				modified = true
			}
		} else {
			_, slot := happening.Dto.GetSlot(request.SlotID)
			if slot == nil {
				return fmt.Errorf("%w: unknown slot ContactID=%v", facade.ErrBadRequest, request.SlotID)
			}
			adjustment.Slot = *slot
			adjustment.Canceled = &models4schedulus.Canceled{
				At:     time.Now(),
				By:     dbmodels.ByUser{UID: uid},
				Reason: request.Reason,
			}
			modified = true
		}

		if err := scheduleDay.Dto.Validate(); err != nil {
			return fmt.Errorf("schedule day record is not valid: %w", err)
		}

		if isNewRecord {
			if err := tx.Insert(ctx, scheduleDay.Record); err != nil {
				return fmt.Errorf("failed to create schedule day record: %w", err)
			}
		} else if modified {
			dayUpdates = append(dayUpdates, dal.Update{
				Field: "cancellations", Value: scheduleDay.Dto.HappeningAdjustments,
			})
			if err := tx.Update(ctx, scheduleDay.Key, dayUpdates); err != nil {
				return fmt.Errorf("failed to update schedule day record: %w", err)
			}
		}
	}

	return nil
}

func createCanceled(uid, reason string) *models4schedulus.Canceled {
	return &models4schedulus.Canceled{
		At:     time.Now(),
		By:     dbmodels.ByUser{UID: uid},
		Reason: strings.TrimSpace(reason),
	}
}
func markRecurringHappeningRecordAsCanceled(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	uid string,
	happening models4schedulus.HappeningContext,
	request dto4schedulus.CancelHappeningRequest,
) error {
	var happeningUpdates []dal.Update
	happening.Dto.Status = models4schedulus.HappeningStatusCanceled
	if happening.Dto.Canceled == nil {
		happening.Dto.Canceled = createCanceled(uid, request.Reason)
	} else if reason := strings.TrimSpace(request.Reason); reason != "" {
		happening.Dto.Canceled.Reason = reason
	}
	happeningUpdates = append(happeningUpdates,
		dal.Update{
			Field: "status",
			Value: happening.Dto.Status,
		},
		dal.Update{
			Field: "canceled",
			Value: happening.Dto.Canceled,
		},
	)
	if err := happening.Dto.Validate(); err != nil {
		return fmt.Errorf("happening record is not valid: %w", err)
	}
	if err := tx.Update(ctx, happening.Key, happeningUpdates); err != nil {
		return fmt.Errorf("faield to update happening record: %w", err)
	}
	return nil
}
