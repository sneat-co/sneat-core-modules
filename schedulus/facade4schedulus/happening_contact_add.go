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

func AddParticipantToHappening(ctx context.Context, user facade.User, request dto4schedulus.HappeningContactRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	var worker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *happeningWorkerParams) error {
		contact, err := getHappeningContactRecords(ctx, tx, &request, params)
		if err != nil {
			return err
		}

		switch params.Happening.Dto.Type {
		case "single":
			break // No special processing needed
		case "recurring":
			if err = addContactToHappeningBriefInTeamDto(ctx, tx, params.TeamModuleEntry, params.Happening, request.Contact.ID); err != nil {
				return fmt.Errorf("failed to add member to happening brief in team DTO: %w", err)
			}
		default:
			return fmt.Errorf("invalid happenning record: %w",
				validation.NewErrBadRecordFieldValue("type",
					fmt.Sprintf("unknown value: [%v]", params.Happening.Dto.Type)))
		}
		params.HappeningUpdates = append(params.HappeningUpdates, params.Happening.Dto.AddContact(request.Contact.TeamID, contact.ID, &contact.Data.ContactBrief)...)
		params.HappeningUpdates = append(params.HappeningUpdates, params.Happening.Dto.AddParticipant(request.Contact.TeamID, contact.ID, nil)...)
		return err
	}

	if err = modifyHappening(ctx, user, request.HappeningRequest, worker); err != nil {
		return fmt.Errorf("failed to add member to happening: %w", err)
	}
	return nil
}

func addContactToHappeningBriefInTeamDto(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	schedulusTeam dal4schedulus.SchedulusTeamContext,
	happening models4schedulus.HappeningContext,
	contactID string,
) (err error) {
	teamID := schedulusTeam.Key.Parent().ID.(string)
	happeningBriefPointer := schedulusTeam.Data.GetRecurringHappeningBrief(happening.ID)
	teamContactID := dbmodels.NewTeamItemID(teamID, contactID)
	var happeningBrief models4schedulus.HappeningBrief
	if happeningBriefPointer == nil {
		happeningBrief = happening.Dto.HappeningBrief // Make copy so we do not affect the DTO object
		happeningBriefPointer = &happeningBrief
	} else if happeningBriefPointer.Participants[string(teamContactID)] != nil {
		return nil // Already added to happening brief in schedulusTeam record
	}
	if happeningBriefPointer.Participants == nil {
		happeningBriefPointer.Participants = make(map[string]*models4schedulus.HappeningParticipant)
	}
	if happeningBriefPointer.Participants[string(teamContactID)] == nil {
		happeningBriefPointer.Participants[string(teamContactID)] = &models4schedulus.HappeningParticipant{}
	}
	if schedulusTeam.Data.RecurringHappenings == nil {
		schedulusTeam.Data.RecurringHappenings = make(map[string]*models4schedulus.HappeningBrief, 1)
	}
	schedulusTeam.Data.RecurringHappenings[happening.ID] = happeningBriefPointer
	teamUpdates := []dal.Update{
		{
			Field: "recurringHappenings." + happening.ID,
			Value: happeningBriefPointer,
		},
	}
	if err = tx.Update(ctx, schedulusTeam.Key, teamUpdates); err != nil {
		return fmt.Errorf("failed to update schedulusTeam record with a member added to a recurring happening: %w", err)
	}
	return nil
}
