package facade4linkage

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-core-modules/linkage/dto4linkage"
	"github.com/strongo/validation"
	"time"
)

func UpdateRelatedFields(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	now time.Time,
	userID string,
	objectRef dbo4linkage.SpaceModuleItemRef,
	request dto4linkage.UpdateRelatedFieldRequest,
	item *dbo4linkage.WithRelatedAndIDsAndUserID,
	addUpdatesToParams func(updates []update.Update),
) (
	recordsUpdates []record.Updates, err error,
) {
	var setRelatedResult SetRelatedResult

	for i, command := range request.Related {
		itemRef := command.ItemRef
		if itemRef == objectRef {
			return nil, validation.NewErrBadRequestFieldValue(fmt.Sprintf("request.Related[%d].ItemRef", i), fmt.Sprintf("same as objectRef: %+v", objectRef))
		}
		if setRelatedResult, err = SetRelated(now, userID, item, objectRef, command); err != nil {
			return nil, err
		}

		addUpdatesToParams(setRelatedResult.ItemUpdates)
		//params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, setRelatedResult.SpaceModuleUpdates...)

		if recordsUpdates, err = updateRelatedItem(ctx, tx, now, objectRef, command); err != nil {
			return recordsUpdates, fmt.Errorf("failed to update related record for command [%d=%s]: %w", i, itemRef.ID(), err)
		}
	}

	return recordsUpdates, nil
}
