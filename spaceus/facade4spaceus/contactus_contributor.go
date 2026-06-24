package facade4spaceus

import (
	"time"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/contactus-ext/backend/contactusmodels/briefs4contactus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

// ContactusSpaceContributor builds the contactus records that must be persisted
// when a new space is created. It is implemented and registered by the contactus
// module so that spaceus does not depend on contactus DAL types directly.
type ContactusSpaceContributor interface {
	// BuildSpaceCreationRecords returns the contactus records (contactus space + creator member)
	// to insert as part of creating a new space.
	BuildSpaceCreationRecords(
		spaceID coretypes.SpaceID,
		userContactID string,
		creatorBrief briefs4contactus.ContactBrief,
		createdAt time.Time,
		byUserID string,
	) (records []dal.Record, err error)
}

var contactusSpaceContributor ContactusSpaceContributor

// RegisterContactusSpaceContributor registers the contactus implementation used by CreateSpace.
// Called once at startup from contactus.Extension().
func RegisterContactusSpaceContributor(c ContactusSpaceContributor) {
	contactusSpaceContributor = c
}
