package facade4linkage

import (
	"context"
	"errors"
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-core-modules/linkage/dto4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
)

func UpdateRelatedItemsWithLatestRelationships(
	ctx facade.ContextWithUser,
	request dto4linkage.UpdateItemRequest,
	itemData dbo4linkage.WithRelatedAndIDs,
) (err error) {
	var updateErrors []error
	for i, related := range request.Related {
		err = updateItemWithLatestRelationshipsFromRelatedItem(ctx, request.SpaceID, related.ItemRef, request.ItemRef, itemData.Related)
		if err != nil {
			updateErrors = append(updateErrors, fmt.Errorf("failed to update related item (%d=%s): %w", i, related.ItemRef.ID(), err))
		}
	}
	if len(updateErrors) > 0 {
		return fmt.Errorf("failed to update %d related items: %w", len(updateErrors), errors.Join(updateErrors...))
	}
	return nil
}

func updateItemWithLatestRelationshipsFromRelatedItem(
	ctx facade.ContextWithUser,
	spaceID coretypes.SpaceID,
	itemRef dbo4linkage.ItemRef,
	relatedItemRef dbo4linkage.ItemRef,
	relatedByModuleOfRelatedItem dbo4linkage.RelatedModules,
) (err error) {
	itemRelationshipInRelatedItem := dbo4linkage.GetRelatedItemByRef(relatedByModuleOfRelatedItem, itemRef, false)
	if itemRelationshipInRelatedItem == nil || len(itemRelationshipInRelatedItem.RolesToItem) == 0 {
		return nil
	}

	// specscore: decisions/0002-reserved-extension-space-ids
	// Fail closed: a cross-space ("@otherSpace") ref is rejected outright; a
	// spaceless (/ext/) ref is allowed through here but its write is authorized
	// per-record against the loaded record's ACL below.
	// https://github.com/sneat-co/sneat-specs/blob/main/spec/decisions/0002-reserved-extension-space-ids.md
	spaceless, err := classifyRelatedItemRef(spaceID, itemRef)
	if err != nil {
		return err
	}
	userID := ctx.User().GetUserID()

	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return err
	}

	return db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {

		key := dbo4spaceus.NewSpaceModuleItemKeyFromItemRef(spaceID, itemRef)
		// The record id must match the document id used to build the key (the
		// bare itemID with any "@{spaceID}" suffix stripped).
		item := record.NewDataWithID(itemRef.DocID(), key, new(dbo4linkage.WithRelatedAndIDsAndUserID))
		if err = tx.Get(ctx, item.Record); err != nil {
			return err
		}
		// specscore: decisions/0002-reserved-extension-space-ids
		// A spaceless /ext/ record is authorized per-record against its own ACL.
		if spaceless {
			if err = authorizeSpacelessRelatedWrite(userID, item.Data.ACL); err != nil {
				return err
			}
		}
		relatedItem := dbo4linkage.GetRelatedItemByRef(item.Data.Related, relatedItemRef, true)

		// We do not override existing roles in related item, so we do not lose a role in case of race condition
		for roleID, role := range itemRelationshipInRelatedItem.RolesToItem {
			relatedItem.RolesOfItem[roleID] = role
		}
		//if item.Data.UserID != "" {
		//	users := make(map[string]dbo4userus.UserEntry)
		//	if err = updateUserWithRelatedTx(ctx, tx, item.Data.UserID, spaceID, users, itemRef, *relatedItem); err != nil {
		//		return err
		//	}
		//
		//}
		return nil
	})
}
