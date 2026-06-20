package facade4userus

import (
	"errors"
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
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

func txSetUserCountry(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, request SetUserCountryRequest, params *dal4userus.UserWorkerParams) (err error) {
	if params.User.Data.CountryID != request.CountryID {
		params.User.Data.CountryID = request.CountryID
		params.User.Record.MarkAsChanged()
		params.UserUpdates = append(params.UserUpdates,
			update.ByFieldName("countryID", request.CountryID))
	}

	if contactusCountryUpdater == nil {
		return errors.New("contactus country updater is not registered")
	}

	userID := ctx.User().GetUserID()
	for spaceID, spaceBrief := range params.User.Data.Spaces {
		if IsUnknownCountryID(spaceBrief.CountryID) && spaceBrief.Type == coretypes.SpaceTypeFamily || spaceBrief.Type == coretypes.SpaceTypePrivate {
			spaceBrief.CountryID = request.CountryID
			params.UserUpdates = append(params.UserUpdates, update.ByFieldName(fmt.Sprintf("spaces.%s.countryID", spaceID), request.CountryID))
		}
		if err = contactusCountryUpdater.UpdateUserContactCountryInSpace(ctx, tx, coretypes.SpaceID(spaceID), userID, request.CountryID); err != nil {
			return fmt.Errorf("failed to update space %s: %w", spaceID, err)
		}
	}
	return
}

// IsUnknownCountryID checks if countryID is empty or "--" - TODO: move next to dbmodels.UnknownCountryID
func IsUnknownCountryID(countryID string) bool {
	return countryID == "" || countryID == with.UnknownCountryID
}
