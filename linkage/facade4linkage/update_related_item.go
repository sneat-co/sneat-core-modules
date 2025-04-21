package facade4linkage

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"time"
)

func updateRelatedItem(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	now time.Time,
	objectRef dbo4linkage.SpaceModuleItemRef,
	command dbo4linkage.RelationshipItemRolesCommand,
) (
	recordsUpdates []record.Updates,
	err error,
) {
	relatedItemRef := command.ItemRef
	if objectRef == relatedItemRef {
		return nil, fmt.Errorf("objectRef and command.ItemRef are the same: %+v", objectRef)
	}
	relatedItemID := relatedItemRef.ItemID

	key := dbo4spaceus.NewSpaceModuleItemKeyFromItemRef(relatedItemRef)
	dbo := new(dbo4linkage.WithRelatedAndIDsAndUserID)
	dbo.WithRelatedAndIDs = new(dbo4linkage.WithRelatedAndIDs)
	related := record.NewDataWithID[string, *dbo4linkage.WithRelatedAndIDsAndUserID](relatedItemID, key, dbo)
	if err = tx.Get(ctx, related.Record); err != nil {
		return recordsUpdates, fmt.Errorf("failed to get related record: %w", err)
	}
	if err = related.Data.Validate(); err != nil {
		return recordsUpdates, fmt.Errorf("record is not valid after loading from DB: %w", err)
	}

	var result SetRelatedResult

	relatedItemCommand := dbo4linkage.RelationshipItemRolesCommand{
		ItemRef: objectRef,
	}
	if command.Add != nil {
		relatedItemCommand.Add = &dbo4linkage.RolesCommand{
			RolesOfItem: command.Add.RolesToItem,
			RolesToItem: command.Add.RolesOfItem,
		}
	}
	if command.Remove != nil {
		relatedItemCommand.Remove = &dbo4linkage.RolesCommand{
			RolesOfItem: command.Remove.RolesToItem,
			RolesToItem: command.Remove.RolesOfItem,
		}
	}

	userID := "u123"
	if result, err = SetRelated(now, userID, dbo, relatedItemRef, relatedItemCommand); err != nil {
		return nil, fmt.Errorf("failed to update related item: %w", err)
	}
	recordsUpdates = append(recordsUpdates, record.Updates{
		Record:  related.Record,
		Updates: result.ItemUpdates,
	})
	if related.Data.UserID != "" {
		var userUpdates record.Updates
		// TODO: Document use case when this is needed and if it is really used
		if userUpdates, err = updateUserRelated(ctx, tx, related); err != nil {
			return recordsUpdates, fmt.Errorf("failed to update related in user record: %w", err)
		} else if len(userUpdates.Updates) > 0 {
			recordsUpdates = append(recordsUpdates, userUpdates)
		}
	}
	return recordsUpdates, nil
}
