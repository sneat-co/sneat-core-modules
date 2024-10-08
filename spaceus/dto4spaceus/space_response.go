package dto4spaceus

import (
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
)

type spaceRecord struct {
	ID  string               `json:"id"`
	Dbo dbo4spaceus.SpaceDbo `json:"dbo"`
}

// SpaceResponse response
type SpaceResponse struct {
	Space spaceRecord `json:"space"`
}
