package facade4schedulus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/schedulus/dal4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/dto4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/models4schedulus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

func AddParticipantToHappening(ctx context.Context, userID string, request dto4schedulus.HappeningContactRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	if userID == "" {
		return fmt.Errorf("%w: user ContactID is missing", facade.ErrUnauthorized)
	}

	var worker = func(ctx context.Context, tx dal.ReadwriteTransaction, params happeningWorkerParams) (err error) {

		if err = tx.GetMulti(ctx, []dal.Record{params.TeamModuleEntry.Record, params.SchedulusTeam.Record}); err != nil {
			return fmt.Errorf("failed to get records: %w", err)
		}

		if params.Happening.Dto.HasTeamContactID(dbmodels.NewTeamItemID(request.TeamID, request.ContactID)) {
			return
		}

		if !params.TeamModuleEntry.Data.HasContact(request.ContactID) {
			return validation.NewErrBadRequestFieldValue("teamContactID", "unknown member ContactID")
		}

		switch params.Happening.Dto.Type {
		case "single":
		case "recurring":
			if err = addContactToHappeningBriefInTeamDto(ctx, tx, params.SchedulusTeam, params.Happening.Dto, request.HappeningID, request.ContactID); err != nil {
				return fmt.Errorf("failed to add member to happening brief in team DTO: %w", err)
			}
		default:
			return fmt.Errorf("invalid happenning record: %w",
				validation.NewErrBadRecordFieldValue("type",
					fmt.Sprintf("unknown value: [%v]", params.Happening.Dto.Type)))
		}
		teamContactID := dbmodels.NewTeamItemID(request.TeamID, request.ContactID)
		if !params.Happening.Dto.HasTeamContactID(teamContactID) {
			params.Happening.Dto.AddTeamContactID(teamContactID)
			params.HappeningUpdates = params.Happening.Dto.WithMultiTeamContacts.Updates()
		}
		return
	}

	if err = modifyHappening(ctx, userID, request.HappeningRequest, worker); err != nil {
		return fmt.Errorf("failed to add member to happening: %w", err)
	}
	return nil
}

func addContactToHappeningBriefInTeamDto(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	schedulusTeam dal4schedulus.SchedulusTeamContext,
	happeningDto *models4schedulus.HappeningDto,
	happeningID string,
	contactID string,
) (err error) {
	happeningBrief := schedulusTeam.Data.GetRecurringHappeningBrief(happeningID)
	teamContactID := dbmodels.NewTeamItemID(schedulusTeam.ID, contactID)
	if happeningBrief != nil && happeningBrief.Participants[string(teamContactID)] != nil {
		return nil // Already added to happening brief in schedulusTeam record
	}
	happeningBrief = &happeningDto.HappeningBrief
	// We have to check again as DTO can have member ContactID while brief does not.
	if happeningBrief.Participants[string(teamContactID)] == nil {
		happeningBrief.Participants[string(teamContactID)] = &models4schedulus.HappeningParticipant{}
	}
	schedulusTeam.Data.RecurringHappenings[happeningID] = happeningBrief
	teamUpdates := []dal.Update{
		{
			Field: "recurringHappenings." + happeningID,
			Value: happeningBrief,
		},
	}
	if err = tx.Update(ctx, schedulusTeam.Key, teamUpdates); err != nil {
		return fmt.Errorf("failed to update schedulusTeam record with a member added to a recurring happening: %w", err)
	}
	return nil
}
