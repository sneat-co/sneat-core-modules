package dal4schedulus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/schedulus/models4schedulus"
)

type SchedulusTeamContext = record.DataWithID[string, *models4schedulus.SchedulusTeamDto]

func NewSchedulusTeamKey(teamID string) *dal.Key {
	return dal.NewKeyWithID("schedulus_team", teamID)
}

func NewSchedulusTeamContext(teamID string) SchedulusTeamContext {
	key := NewSchedulusTeamKey(teamID)
	return record.NewDataWithID(teamID, key, new(models4schedulus.SchedulusTeamDto))
}

func GetSchedulusTeam(ctx context.Context, tx dal.ReadwriteTransaction, teamID string) (SchedulusTeamContext, error) {
	schedulusTeam := NewSchedulusTeamContext(teamID)
	return schedulusTeam, tx.Get(ctx, schedulusTeam.Record)
}
