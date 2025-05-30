package facade4userus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
)

type SetUserCountryRequest struct {
	CountryID string `json:"countryID"`
}

func (v SetUserCountryRequest) Validate() error {
	if v.CountryID == "" {
		return validation.NewErrRequestIsMissingRequiredField("countryID")
	}
	if len(v.CountryID) != 2 {
		return validation.NewErrBadRequestFieldValue("countryID", "must be 2 characters long")
	}
	return nil
}

func SetUserCountry(ctx facade.ContextWithUser, request SetUserCountryRequest) (err error) {
	return dal4userus.RunUserWorker(ctx, true,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4userus.UserWorkerParams) (err error) {
			if err = txSetUserCountry(ctx, tx, request, params); err != nil {
				return fmt.Errorf("failed in txSetUserCountry(): %w", err)
			}
			return
		})
}

type RecordToUpdate struct {
	Key     *dal.Key
	Updates []update.Update
}

func txSetUserCountry(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, request SetUserCountryRequest, params *dal4userus.UserWorkerParams) (err error) {
	if params.User.Data.CountryID != request.CountryID {
		params.User.Data.CountryID = request.CountryID
		params.User.Record.MarkAsChanged()
		params.UserUpdates = append(params.UserUpdates,
			update.ByFieldName("countryID", request.CountryID))
	}

	recordsToUpdate := make([]RecordToUpdate, 0, len(params.User.Data.Spaces))

	for spaceID, spaceBrief := range params.User.Data.Spaces {
		if IsUnknownCountryID(spaceBrief.CountryID) && spaceBrief.Type == coretypes.SpaceTypeFamily || spaceBrief.Type == coretypes.SpaceTypePrivate {
			spaceBrief.CountryID = request.CountryID
			params.UserUpdates = append(params.UserUpdates, update.ByFieldName(fmt.Sprintf("spaces.%s.countryID", spaceID), request.CountryID))
		}
		if err = dal4contactus.RunContactusSpaceWorkerNoUpdate(ctx, tx, coretypes.SpaceID(spaceID),
			func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusSpaceWorkerParams) (err error) {
				if err = params.GetRecords(ctx, tx, params.Space.Record); err != nil {
					return
				}
				if IsUnknownCountryID(params.Space.Data.CountryID) {
					params.Space.Data.CountryID = request.CountryID
					params.SpaceUpdates = append(params.SpaceUpdates, update.ByFieldName("countryID", request.CountryID))
					params.Space.Record.MarkAsChanged()
				}
				userID := ctx.User().GetUserID()
				userContactID, userContactBrief := params.SpaceModuleEntry.Data.GetContactBriefByUserID(userID)
				if userContactBrief != nil && IsUnknownCountryID(userContactBrief.CountryID) {
					userContactBrief.CountryID = request.CountryID
					params.SpaceModuleUpdates = append(params.SpaceModuleUpdates,
						update.ByFieldPath([]string{"contacts", userContactID, "countryID"},
							request.CountryID))
					params.SpaceModuleEntry.Record.MarkAsChanged()
					userContact := dal4contactus.NewContactEntry(coretypes.SpaceID(spaceID), userContactID)
					if err = tx.Get(ctx, userContact.Record); err != nil {
						return
					}
					if IsUnknownCountryID(userContact.Data.CountryID) {
						userContact.Data.CountryID = request.CountryID
						recordsToUpdate = append(recordsToUpdate, RecordToUpdate{Key: userContact.Key, Updates: []update.Update{update.ByFieldName("countryID", request.CountryID)}})
					}
				}
				if params.Space.Record.HasChanged() && len(params.SpaceUpdates) > 0 {
					recordsToUpdate = append(recordsToUpdate, RecordToUpdate{Key: params.Space.Key, Updates: params.SpaceUpdates})
				}
				if params.SpaceModuleEntry.Record.HasChanged() && len(params.SpaceModuleUpdates) > 0 {
					recordsToUpdate = append(recordsToUpdate, RecordToUpdate{Key: params.SpaceModuleEntry.Key, Updates: params.SpaceModuleUpdates})
				}
				return
			}); err != nil {
			return fmt.Errorf("failed to update space %s: %w", spaceID, err)
		}
	}
	if len(recordsToUpdate) > 0 {
		for _, rec := range recordsToUpdate {
			if err = tx.Update(ctx, rec.Key, rec.Updates); err != nil {
				return fmt.Errorf("failed to update record %s: %w", rec.Key, err)
			}
		}
	}
	return
}

// IsUnknownCountryID checks if countryID is empty or "--" - TODO: move next to dbmodels.UnknownCountryID
func IsUnknownCountryID(countryID string) bool {
	return countryID == "" || countryID == with.UnknownCountryID
}
