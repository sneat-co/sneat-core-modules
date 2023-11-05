package facade4schedulus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/schedulus/const4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/dto4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/models4schedulus"
	"github.com/sneat-co/sneat-core-modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/slice"
	"strings"
)

// CreateHappening creates a recurring happening
func CreateHappening(
	ctx context.Context, user facade.User, request dto4schedulus.CreateHappeningRequest,
) (
	response dto4schedulus.CreateHappeningResponse, err error,
) {
	request.Happening.Title = strings.TrimSpace(request.Happening.Title)
	if err = request.Validate(); err != nil {
		return
	}
	var counter string
	if request.Happening.Type == models4schedulus.HappeningTypeRecurring {
		counter = "recurringHappenings"
	}
	happeningDto := &models4schedulus.HappeningDto{
		HappeningBrief: *request.Happening,
		WithTeamDates: dbmodels.WithTeamDates{
			WithTeamIDs: dbmodels.WithTeamIDs{
				TeamIDs: []string{request.TeamID},
			},
		},
	}
	happeningDto.ContactIDs = append(happeningDto.ContactIDs, "*")
	if len(happeningDto.AssetIDs) == 0 {
		happeningDto.AssetIDs = []string{"*"}
	}

	if happeningDto.Type == models4schedulus.HappeningTypeSingle {
		for _, slot := range happeningDto.Slots {
			date := slot.Start.Date
			if slice.Index(happeningDto.Dates, date) < 0 {
				happeningDto.Dates = append(happeningDto.Dates, date)
			}
		}
	}
	err = dal4teamus.CreateTeamItem(ctx, user, counter, request.TeamRequest,
		const4schedulus.ModuleID,
		new(models4schedulus.SchedulusTeamDto),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.ModuleTeamWorkerParams[*models4schedulus.SchedulusTeamDto]) (err error) {
			contactusTeam := dal4contactus.NewContactusTeamModuleEntry(params.Team.ID)
			if err = params.GetRecords(ctx, tx, params.UserID, contactusTeam.Record); err != nil {
				return err
			}

			happeningDto.UserIDs = params.Team.Data.UserIDs
			happeningDto.Status = "active"
			if happeningDto.Type == "single" {
				date := happeningDto.Slots[0].Start.Date
				happeningDto.Dates = []string{date}
				happeningDto.DateMin = date
				happeningDto.DateMax = date
			}

			contactsByTeamID := make(map[string][]dal4contactus.ContactEntry)

			for participantID := range happeningDto.Participants {
				participantKey := dbmodels.TeamItemID(participantID)
				teamID := participantKey.TeamID()
				if teamID == params.Team.ID {
					contactBrief := contactusTeam.Data.Contacts[participantKey.ItemID()]
					if contactBrief == nil {
						teamContacts := contactsByTeamID[teamID]
						if teamContacts == nil {
							teamContacts = make([]dal4contactus.ContactEntry, 0, 1)
						}
						contactsByTeamID[teamID] = append(teamContacts, dal4contactus.NewContactEntry(teamID, participantKey.ItemID()))
					} else {
						happeningDto.AddContact(teamID, participantKey.ItemID(), contactBrief)
					}
				} else {
					return errors.New("not implemented yet: adding participants from other teams at happening creation")
				}
			}

			if len(contactsByTeamID) > 0 {
				contactRecords := make([]dal.Record, 0)
				for _, teamContacts := range contactsByTeamID {
					for _, contact := range teamContacts {
						contactRecords = append(contactRecords, contact.Record)
					}
				}
				if err = tx.GetMulti(ctx, contactRecords); err != nil {
					return err
				}
				for teamID, teamContacts := range contactsByTeamID {
					for _, contact := range teamContacts {
						happeningDto.AddContact(teamID, contact.ID, &contact.Data.ContactBrief)
					}
				}
			}

			var happeningID string
			var happeningKey *dal.Key
			if happeningID, happeningKey, err = dal4teamus.GenerateNewTeamModuleItemKey(
				ctx, tx, params.Team.ID, moduleID, happeningsCollection, 5, 10); err != nil {
				return err
			}
			response.ID = happeningID
			record := dal.NewRecordWithData(happeningKey, happeningDto)
			if err = happeningDto.Validate(); err != nil {
				return fmt.Errorf("happening record is not valid for insertion: %w", err)
			}
			//panic("teamDates: " + strings.Join(happeningDto.TeamDates, ","))
			if err = tx.Insert(ctx, record); err != nil {
				return fmt.Errorf("failed to insert new happening record: %w", err)
			}
			if happeningDto.Type == models4schedulus.HappeningTypeRecurring {
				brief := &happeningDto.HappeningBrief
				if params.TeamModuleEntry.Data.RecurringHappenings == nil {
					params.TeamModuleEntry.Data.RecurringHappenings = make(map[string]*models4schedulus.HappeningBrief)
				}
				params.TeamModuleEntry.Data.RecurringHappenings[happeningID] = brief
				params.TeamModuleUpdates = append(params.TeamUpdates, dal.Update{
					Field: "recurringHappenings." + happeningID,
					Value: brief,
				})
			}
			return nil
		},
	)
	response.Dto = *happeningDto
	return
}
