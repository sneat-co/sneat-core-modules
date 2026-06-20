package contactus

import (
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/strongoapp/with"
)

// userusContactusContributor implements facade4userus.ContactusCountryUpdater.
// It updates the user's contact country within a space, keeping userus decoupled
// from contactus DAL types. Logic moved here from facade4userus/set_user_country.go
// as part of the contactus cycle-break.
type userusContactusContributor struct{}

func (userusContactusContributor) UpdateUserContactCountryInSpace(
	ctx facade.ContextWithUser,
	tx dal.ReadwriteTransaction,
	spaceID coretypes.SpaceID,
	userID string,
	countryID string,
) error {
	return dal4contactus.RunContactusSpaceWorkerNoUpdate(ctx, tx, spaceID,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusSpaceWorkerParams) (err error) {
			if err = params.GetRecords(ctx, tx, params.Space.Record); err != nil {
				return
			}

			type recordToUpdate struct {
				key     *dal.Key
				updates []update.Update
			}
			var recordsToUpdate []recordToUpdate

			if isUnknownCountryID(params.Space.Data.CountryID) {
				params.Space.Data.CountryID = countryID
				params.SpaceUpdates = append(params.SpaceUpdates, update.ByFieldName("countryID", countryID))
				params.Space.Record.MarkAsChanged()
			}
			contactID, userContactBrief := params.SpaceModuleEntry.Data.GetContactBriefByUserID(userID)
			if userContactBrief != nil && isUnknownCountryID(userContactBrief.CountryID) {
				userContactBrief.CountryID = countryID
				params.SpaceModuleUpdates = append(params.SpaceModuleUpdates,
					update.ByFieldPath([]string{"contacts", contactID, "countryID"}, countryID))
				params.SpaceModuleEntry.Record.MarkAsChanged()
				userContact := dal4contactus.NewContactEntry(spaceID, contactID)
				if err = tx.Get(ctx, userContact.Record); err != nil {
					return
				}
				if isUnknownCountryID(userContact.Data.CountryID) {
					userContact.Data.CountryID = countryID
					recordsToUpdate = append(recordsToUpdate, recordToUpdate{key: userContact.Key, updates: []update.Update{update.ByFieldName("countryID", countryID)}})
				}
			}
			if params.Space.Record.HasChanged() && len(params.SpaceUpdates) > 0 {
				recordsToUpdate = append(recordsToUpdate, recordToUpdate{key: params.Space.Key, updates: params.SpaceUpdates})
			}
			if params.SpaceModuleEntry.Record.HasChanged() && len(params.SpaceModuleUpdates) > 0 {
				recordsToUpdate = append(recordsToUpdate, recordToUpdate{key: params.SpaceModuleEntry.Key, updates: params.SpaceModuleUpdates})
			}
			for _, rec := range recordsToUpdate {
				if err = tx.Update(ctx, rec.key, rec.updates); err != nil {
					return fmt.Errorf("failed to update record %s: %w", rec.key, err)
				}
			}
			return
		})
}

// isUnknownCountryID mirrors facade4userus.IsUnknownCountryID without importing it,
// to avoid a contactus -> userus dependency cycle.
func isUnknownCountryID(countryID string) bool {
	return countryID == "" || countryID == with.UnknownCountryID
}
