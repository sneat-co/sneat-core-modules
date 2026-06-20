package facade4userus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
)

// ContactusCountryUpdater updates the user's contact country within a space.
// It is implemented and registered by the contactus module so that userus does
// not depend on contactus DAL types directly.
type ContactusCountryUpdater interface {
	UpdateUserContactCountryInSpace(
		ctx facade.ContextWithUser,
		tx dal.ReadwriteTransaction,
		spaceID coretypes.SpaceID,
		userID string,
		countryID string,
	) error
}

var contactusCountryUpdater ContactusCountryUpdater

// RegisterContactusCountryUpdater registers the contactus implementation used by SetUserCountry.
// Called once at startup from contactus.Extension().
func RegisterContactusCountryUpdater(u ContactusCountryUpdater) {
	contactusCountryUpdater = u
}
