package facade4linkage

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-core-modules/linkage/dto4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"time"
)

var relatedDboFactories = make(map[coretypes.ModuleID]map[string]RelatedDboFactory)

func RegisterDboFactory(moduleID coretypes.ModuleID, collection string, f RelatedDboFactory) {
	m, ok := relatedDboFactories[moduleID]
	if !ok {
		m = make(map[string]RelatedDboFactory)
		relatedDboFactories[moduleID] = m
	}
	if _, ok = m[collection]; ok {
		panic(fmt.Sprintf("duplicate registration of dbo factory for module=%s and collection=%s", string(moduleID), collection))
	}
	m[collection] = f
}

func getDboFactory(moduleID coretypes.ModuleID, collection string) RelatedDboFactory {
	m := relatedDboFactories[moduleID]
	if m == nil {
		return nil
	}
	return m[collection]
}

func UpdateRelatedAndIDsOfSpaceItem(ctx facade.ContextWithUser, request dto4linkage.UpdateRelatedRequest) (err error) {
	dboFactory := getDboFactory(request.ModuleID, request.Collection)
	if dboFactory == nil {
		return validation.NewBadRequestError(fmt.Errorf("unknown moduleID or collection: %s/%s", request.ModuleID, request.Collection))
	}

	moduleDbo := dboFactory.NewSpaceModuleDbo()
	itemDbo := dboFactory.NewItemDbo()

	return dal4spaceus.RunSpaceItemWorker(ctx,
		request.SpaceItemRequest, request.ModuleID, moduleDbo, request.Collection, itemDbo,
		func(
			ctx facade.ContextWithUser, tx dal.ReadwriteTransaction,
			params *dal4spaceus.SpaceItemWorkerParams[dal4spaceus.SpaceModuleDbo, SpaceItemDboWithRelatedAndIDs],
		) (err error) {
			if err = params.GetRecords(ctx, tx); err != nil {
				return err
			}
			spaceItemUpdates, recordsUpdates, err := updateRelatedTxWorker(ctx, tx, params.Started, request, params.SpaceModuleEntry.Data, params.SpaceItem.Data)
			if err != nil {
				return
			}
			params.SpaceItemUpdates = append(params.SpaceItemUpdates, spaceItemUpdates...)
			// This is updates of the related DBOs
			params.RecordUpdates = append(params.RecordUpdates, recordsUpdates...)
			return nil
		},
	)
}

func updateRelatedTxWorker(
	ctx facade.ContextWithUser,
	tx dal.ReadwriteTransaction,
	now time.Time,
	request dto4linkage.UpdateRelatedRequest,
	_ dal4spaceus.SpaceModuleDbo, // TODO: apply updates to this DBO
	spaceItemDbo SpaceItemDboWithRelatedAndIDs,
) (spaceItemUpdates []update.Update, recordsUpdates []record.Updates, err error) {
	itemRef := dbo4linkage.ItemRef{
		Module:     const4contactus.ModuleID,
		Collection: const4contactus.ContactsCollection,
		ItemID:     request.ID,
	}
	userID := ctx.User().GetUserID()
	recordsUpdates, err = UpdateRelatedFields(ctx, tx,
		now,
		userID,
		request.SpaceID,
		itemRef, request.UpdateRelatedFieldRequest,
		&dbo4linkage.WithRelatedAndIDsAndUserID{
			WithUserID: dbmodels.WithUserID{
				UserID: userID,
			},
			WithRelatedAndIDs: spaceItemDbo.RelatedAndIDs(),
		},
		func(updates []update.Update) {
			spaceItemUpdates = append(spaceItemUpdates, updates...)
		})
	if err != nil {
		return
	}
	// Below is done in UpdateRelatedFields
	//if err = updateRelatedDbos(ctx, tx, request); err != nil {
	//	return nil, err
	//}
	return
}
