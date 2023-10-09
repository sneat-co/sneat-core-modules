package dal4teamus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core"
)

// TeamsCollection table name
const TeamsCollection = "teams"
const TeamBriefsCollection = "briefs"

// NewTeamKey create new doc ref
func NewTeamKey(id string) *dal.Key {
	const max = 30
	if id == "" {
		panic("empty team ItemID")
	}
	if l := len(id); l > max {
		panic(fmt.Sprintf("team ItemID is %v characters long exceded what is %v more then max %v", l, max-l, max))
	}
	if !core.IsAlphanumericOrUnderscore(id) {
		panic(fmt.Sprintf("team ItemID has non alphanumeric characters or letters in upper case: [%v]", id))
	}
	return dal.NewKeyWithID(TeamsCollection, id)
}
