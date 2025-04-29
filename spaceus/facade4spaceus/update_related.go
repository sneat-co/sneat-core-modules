package facade4spaceus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-core-modules/linkage/dto4linkage"
	"github.com/sneat-co/sneat-core-modules/linkage/facade4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

type RelatedDboFactory func() (dal4spaceus.SpaceModuleDbo, dal4spaceus.SpaceItemDbo, *dbo4linkage.WithRelatedAndIDs)

var relatedDboFactory = make(map[coretypes.ModuleID]map[string]RelatedDboFactory)

func RegisterDboFactory(moduleID coretypes.ModuleID, collection string, f RelatedDboFactory) {
	m, ok := relatedDboFactory[moduleID]
	if !ok {
		m = make(map[string]RelatedDboFactory)
		relatedDboFactory[moduleID] = m
	}
	if _, ok = m[collection]; ok {
		panic(fmt.Sprintf("duplicate registration of dbo factory for module=%s and collection=%s", string(moduleID), collection))
	}
	m[collection] = f
}

func getDboFactory(moduleID coretypes.ModuleID, collection string) RelatedDboFactory {
	m := relatedDboFactory[moduleID]
	if m == nil {
		return nil
	}
	return m[collection]
}

func UpdateRelated(ctx facade.ContextWithUser, request dto4linkage.UpdateRelatedRequest) (err error) {
	dboFactory := getDboFactory(request.ModuleID, request.Collection)
	if dboFactory == nil {
		return validation.NewBadRequestError(fmt.Errorf("unknown moduleID or collection: %s/%s", request.ModuleID, request.Collection))
	}

	moduleDbo, itemDbo, itemWithRelatedAndIDs := dboFactory()

	return dal4spaceus.RunSpaceItemWorker(ctx,
		request.SpaceItemRequest,
		request.ModuleID, moduleDbo,
		request.Collection, itemDbo,
		func(
			ctx facade.ContextWithUser, tx dal.ReadwriteTransaction,
			params *dal4spaceus.SpaceItemWorkerParams[dal4spaceus.SpaceModuleDbo, dal4spaceus.SpaceItemDbo],
		) (err error) {
			itemRef := dbo4linkage.ItemRef{
				Module:     const4contactus.ModuleID,
				Collection: const4contactus.ContactsCollection,
				ItemID:     request.ID,
			}
			var recordsUpdates []record.Updates
			userID := ctx.User().GetUserID()
			recordsUpdates, err = facade4linkage.UpdateRelatedFields(ctx, tx,
				params.Started,
				userID,
				request.SpaceID,
				itemRef, request.UpdateRelatedFieldRequest,
				&dbo4linkage.WithRelatedAndIDsAndUserID{
					WithUserID: dbmodels.WithUserID{
						UserID: userID,
					},
					WithRelatedAndIDs: itemWithRelatedAndIDs,
				},
				func(updates []update.Update) {
					params.SpaceItemUpdates = append(params.SpaceItemUpdates, updates...)
				})
			if err != nil {
				return err
			}
			params.RecordUpdates = append(params.RecordUpdates, recordsUpdates...)
			return nil
		},
	)
}
