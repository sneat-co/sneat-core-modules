package facade4schedulus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/schedulus/dal4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/dto4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/models4schedulus"
	"github.com/sneat-co/sneat-core-modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

func RemoveParticipantFromHappening(ctx context.Context, userID string, request dto4schedulus.HappeningContactRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	var worker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *happeningWorkerParams) (err error) {
		teamContactID := dbmodels.NewTeamItemID(request.TeamID, request.Contact.ID)
		contactIDs := make([]string, 0, len(params.Happening.Dto.ContactIDs))
		for _, id := range params.Happening.Dto.ContactIDs {
			if id != string(teamContactID) {
				contactIDs = append(contactIDs, id)
			}
		}
		if len(contactIDs) < len(params.Happening.Dto.ContactIDs) {
			params.Happening.Dto.ContactIDs = contactIDs
			params.HappeningUpdates = []dal.Update{
				{
					Field: "contactIDs",
					Value: params.Happening.Dto.ContactIDs,
				},
			}
		}
		switch params.Happening.Dto.Type {
		case "single": // nothing to do
		case "recurring":
			team := dal4teamus.NewTeamContext(request.TeamID)
			if err = tx.Get(ctx, team.Record); err != nil {
				return fmt.Errorf("failed to get team record: %w", err)
			}
			if err = removeContactFromHappeningBriefInTeamDto(ctx, tx, params.TeamModuleEntry, params.Happening.Dto, request.HappeningID, teamContactID); err != nil {
				return fmt.Errorf("failed to remove member from happening brief in team DTO: %w", err)
			}
		default:
			return fmt.Errorf("invalid happenning record: %w",
				validation.NewErrBadRecordFieldValue("type",
					fmt.Sprintf("unknown value: [%v]", params.Happening.Dto.Type)))
		}
		return
	}

	if err = modifyHappening(ctx, userID, request.HappeningRequest, worker); err != nil {
		return err
	}
	return nil
}

func removeContactFromHappeningBriefInTeamDto(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	schedulusTeam dal4schedulus.SchedulusTeamContext,
	happeningDto *models4schedulus.HappeningDto,
	happeningID string,
	teamContactID dbmodels.TeamItemID,
) (err error) {
	happeningBrief := schedulusTeam.Data.GetRecurringHappeningBrief(happeningID)
	if happeningBrief == nil {
		schedulusTeam.Data.RecurringHappenings[happeningID] = &happeningDto.HappeningBrief
	} else if happeningBrief.Participants[string(teamContactID)] != nil {
		happeningBrief.Participants[string(teamContactID)] = new(models4schedulus.HappeningParticipant)
	} else {
		return nil
	}
	teamUpdates := []dal.Update{
		{
			Field: "recurringHappenings." + happeningID,
			Value: happeningBrief,
		},
	}
	if err = tx.Update(ctx, schedulusTeam.Key, teamUpdates); err != nil {
		return fmt.Errorf("failed to update schedulusTeam record: %w", err)
	}

	return nil
}
