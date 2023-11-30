package facade4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-core-modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-core-modules/calendarium/models4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

func AddParticipantToHappening(ctx context.Context, user facade.User, request dto4calendarium.HappeningContactRequest) (err error) {
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
	calendariumTeam dal4calendarium.CalendariumTeamContext,
	happening models4calendarium.HappeningContext,
	contactID string,
) (err error) {
	teamID := calendariumTeam.Key.Parent().ID.(string)
	happeningBriefPointer := calendariumTeam.Data.GetRecurringHappeningBrief(happening.ID)
	teamContactID := dbmodels.NewTeamItemID(teamID, contactID)
	var happeningBrief models4calendarium.HappeningBrief
	if happeningBriefPointer == nil {
		happeningBrief = happening.Dto.HappeningBrief // Make copy so we do not affect the DTO object
		happeningBriefPointer = &happeningBrief
	} else if happeningBriefPointer.Participants[string(teamContactID)] != nil {
		return nil // Already added to happening brief in calendariumTeam record
	}
	if happeningBriefPointer.Participants == nil {
		happeningBriefPointer.Participants = make(map[string]*models4calendarium.HappeningParticipant)
	}
	if happeningBriefPointer.Participants[string(teamContactID)] == nil {
		happeningBriefPointer.Participants[string(teamContactID)] = &models4calendarium.HappeningParticipant{}
	}
	if calendariumTeam.Data.RecurringHappenings == nil {
		calendariumTeam.Data.RecurringHappenings = make(map[string]*models4calendarium.HappeningBrief, 1)
	}
	calendariumTeam.Data.RecurringHappenings[happening.ID] = happeningBriefPointer
	teamUpdates := []dal.Update{
		{
			Field: "recurringHappenings." + happening.ID,
			Value: happeningBriefPointer,
		},
	}
	if err = tx.Update(ctx, calendariumTeam.Key, teamUpdates); err != nil {
		return fmt.Errorf("failed to update calendariumTeam record with a member added to a recurring happening: %w", err)
	}
	return nil
}