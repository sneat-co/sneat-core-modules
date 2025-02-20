package facade4invitus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
)

type CreateInviteResponse struct {
	Invite         InviteEntry
	Contact        dal4contactus.ContactEntry
	ContactusSpace dal4contactus.ContactusSpaceEntry
	Space          dbo4spaceus.SpaceEntry
}
