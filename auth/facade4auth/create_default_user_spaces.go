package facade4auth

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/facade4spaceus"
)

func createDefaultUserSpacesTx(ctx context.Context, tx dal.ReadwriteTransaction, params *CreateUserWorkerParams) (err error) {
	for _, spaceType := range []core4spaceus.SpaceType{core4spaceus.SpaceTypeFamily, core4spaceus.SpaceTypePrivate} {
		if spaceID, _ := params.User.Data.GetFirstSpaceBriefBySpaceType(spaceType); spaceID == "" {
			createSpaceParams := facade4spaceus.CreateSpaceParams{
				User:              params.User,
				WithRecordChanges: &params.WithRecordChanges,
			}
			spaceRequest := dto4spaceus.CreateSpaceRequest{Type: spaceType}
			if err = facade4spaceus.CreateSpaceTxWorker(ctx, tx, params.Started, spaceRequest, &createSpaceParams); err != nil {
				return
			}
		}
	}
	return
}
