package dal4contactus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core"
)

// NewContactKey creates a new contact's key in format "spaceID:memberID"
func NewContactKey(spaceID, contactID string) *dal.Key {
	if !core.IsAlphanumericOrUnderscore(contactID) {
		panic(fmt.Errorf("contactID should be alphanumeric, got: [%s]", contactID))
	}
	spaceModuleKey := dbo4spaceus.NewSpaceModuleKey(spaceID, const4contactus.ModuleID)
	return dal.NewKeyWithParentAndID(spaceModuleKey, dbo4contactus.SpaceContactsCollection, contactID)
}
