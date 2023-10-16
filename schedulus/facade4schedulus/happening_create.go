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
	"github.com/strongo/random"
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

			for contactID := range happeningDto.Participants {
				contactBrief := contactusTeam.Data.Contacts[contactID]
				happeningDto.AddContact(params.Team.ID, contactID, contactBrief)
			}

			var happeningID string
			var happeningKey *dal.Key
			if happeningID, happeningKey, err = newHappeningKey(ctx, models4schedulus.HappeningsCollection, tx); err != nil {
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

// TODO: Implement & reuse a generic function?
func newHappeningKey(ctx context.Context, collection string, tx dal.ReadwriteTransaction) (id string, key *dal.Key, err error) {
	const maxAttemptsCount = 10
	for i := 0; i < maxAttemptsCount; i++ {
		id = random.ID(7)
		key = dal.NewKeyWithID(collection, id)
		record := dal.NewRecordWithData(key, make(map[string]interface{}))
		if err := tx.Get(ctx, record); err != nil { // TODO: use tx.Exists()
			if dal.IsNotFound(err) {
				return id, key, nil
			}
			return "", nil, err
		}
	}
	return "", nil, errors.New("too many attempts  to generate a random happening ContactID")
}
