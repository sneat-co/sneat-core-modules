package dto4spaceus

import (
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

type spaceRecord struct {
	ID  coretypes.SpaceID    `json:"id"`
	Dbo dbo4spaceus.SpaceDbo `json:"dbo"`
}

// SpaceResponse response
type SpaceResponse struct {
	Space spaceRecord `json:"space"`
}
