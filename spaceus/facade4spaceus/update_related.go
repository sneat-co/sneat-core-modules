package facade4spaceus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-core-modules/linkage/facade4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

func UpdateRelated(ctx facade.ContextWithUser, request dto4spaceus.UpdateRelatedRequest) (err error) {
	userCtx := ctx.User()
	var moduleDbo dal4spaceus.SpaceModuleDbo
	var itemDbo interface {
		dal4spaceus.SpaceItemDbo
		GetUserID() string
	}
	var itemWithRelatedAndIDs *dbo4linkage.WithRelatedAndIDs

	switch request.ModuleID {
	case const4contactus.ModuleID:
		moduleDbo = new(dbo4contactus.ContactusSpaceDbo)
		contactDbo := new(dbo4contactus.ContactDbo)
		itemDbo = contactDbo
		itemWithRelatedAndIDs = &contactDbo.WithRelatedAndIDs
	}
	return dal4spaceus.RunSpaceItemWorker(ctx, userCtx,
		request.SpaceItemRequest,
		request.ModuleID, moduleDbo,
		request.Collection, itemDbo,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4spaceus.SpaceItemWorkerParams[dal4spaceus.SpaceModuleDbo, interface {
			dal4spaceus.SpaceItemDbo
			GetUserID() string
		}]) (err error) {
			itemRef := dbo4linkage.SpaceModuleItemRef{
				Module:     const4contactus.ModuleID,
				Collection: const4contactus.ContactsCollection,
				Space:      request.SpaceID,
				ItemID:     request.ID,
			}
			var recordsUpdates []record.Updates
			userID := params.UserID()
			recordsUpdates, err = facade4linkage.UpdateRelatedFields(ctx, tx,
				params.Started,
				userID,
				itemRef, request.UpdateRelatedFieldRequest,
				&dbo4linkage.WithRelatedAndIDsAndUserID{
					WithUserID: dbmodels.WithUserID{
						UserID: itemDbo.GetUserID(),
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
