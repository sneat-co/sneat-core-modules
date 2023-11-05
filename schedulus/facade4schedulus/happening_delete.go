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
)

// DeleteHappening deletes happening
func DeleteHappening(ctx context.Context, user facade.User, request dto4schedulus.HappeningRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	happening := models4schedulus.NewHappeningContext(request.TeamID, request.HappeningID)
	err = dal4teamus.RunModuleTeamWorker(ctx, user, request.TeamRequest,
		const4schedulus.ModuleID,
		new(models4schedulus.SchedulusTeamDto),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.ModuleTeamWorkerParams[*models4schedulus.SchedulusTeamDto]) (err error) {
			hasHappeningRecord := true
			if err = tx.Get(ctx, happening.Record); err != nil {
				if dal.IsNotFound(err) {
					hasHappeningRecord = false
					happening.Dto.Type = request.HappeningType
				} else {
					return fmt.Errorf("failed to get happening: %w", err)
				}
			}
			switch happening.Dto.Type {
			case "":
				return fmt.Errorf("unknown happening type: %w", validation.NewErrRecordIsMissingRequiredField("type"))
			case "single":
			case "recurring":
				happeningBrief := params.TeamModuleEntry.Data.GetRecurringHappeningBrief(request.HappeningID)

				if happeningBrief != nil {
					delete(params.TeamModuleEntry.Data.RecurringHappenings, request.HappeningID)
					params.TeamUpdates = append(params.TeamUpdates, dal.Update{
						Field: "recurringHappenings." + request.HappeningID,
						Value: dal.DeleteField,
					})
					params.TeamUpdates = append(params.TeamUpdates, dal.Update{
						Field: "numberOf.recurringHappenings",
						Value: len(params.TeamModuleEntry.Data.RecurringHappenings),
					})
				}
			default:
				return validation.NewErrBadRecordFieldValue("type", "happening has unknown type: "+happening.Dto.Type)
			}
			if hasHappeningRecord {
				if err = tx.Delete(ctx, happening.Key); err != nil {
					return fmt.Errorf("faield to delete happening record: %w", err)
				}
			}
			return nil
		})
	if err != nil {
		return fmt.Errorf("failed to delete happening: %w", err)
	}
	return
}
