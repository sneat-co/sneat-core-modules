package facade4invitus

import (
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
)

type CreateInviteResponse struct {
	Invite InviteEntry
	Space  dbo4spaceus.SpaceEntry
}
