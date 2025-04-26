package facade4linkage

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-core-modules/linkage/dto4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
)

func UpdateItemRelationships(ctx facade.ContextWithUser, request dto4linkage.UpdateItemRequest) (item record.DataWithID[string, *dbo4linkage.WithRelatedAndIDsAndUserID], err error) {
	if err = dal4spaceus.RunSpaceWorkerWithUserContext(ctx, ctx.User(), request.SpaceID, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4spaceus.SpaceWorkerParams) (err error) {
		item, err = txUpdateItemRelationships(ctx, tx, params, request)
		return err
	}); err != nil {
		return item, err
	}
	if err = UpdateRelatedItemsWithLatestRelationships(ctx, ctx.User(), request, *item.Data.WithRelatedAndIDs); err != nil {
		return item, err
	}
	return item, err
}

func txUpdateItemRelationships(
	ctx context.Context, tx dal.ReadwriteTransaction,
	params *dal4spaceus.SpaceWorkerParams,
	request dto4linkage.UpdateItemRequest,
) (item record.DataWithID[string, *dbo4linkage.WithRelatedAndIDsAndUserID], err error) {
	key := dbo4spaceus.NewSpaceModuleItemKey(request.SpaceID, request.Module, request.Collection, request.ItemID)
	item = record.NewDataWithID[string, *dbo4linkage.WithRelatedAndIDsAndUserID](request.ItemID, key, new(dbo4linkage.WithRelatedAndIDsAndUserID))
	if err = tx.Get(ctx, item.Record); err != nil {
		return item, err
	}
	var itemUpdates []update.Update
	userID := params.UserID()
	params.RecordUpdates, err = UpdateRelatedFields(ctx, tx,
		params.Started,
		userID,
		request.SpaceID,
		request.SpaceModuleItemRef, request.UpdateRelatedFieldRequest, item.Data,
		func(updates []update.Update) {
			itemUpdates = append(itemUpdates, updates...)
		})
	if err != nil {
		return item, err
	}
	if len(itemUpdates) > 0 {
		if err = tx.Update(ctx, item.Key, itemUpdates); err != nil {
			return item, err
		}
	}
	return item, nil
}
