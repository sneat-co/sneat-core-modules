package facade4linkage

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/linkage/models4linkage"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
	"time"
)

type RelatableAdapter[D models4linkage.Relatable] interface {
	VerifyItem(ctx context.Context, tx dal.ReadTransaction, recordRef models4linkage.TeamModuleDocRef) (err error)
	//GetRecord(ctx context.Context, tx dal.ReadTransaction, recordRef models4linkage.TeamModuleDocRef) (record.DataWithID[string, D], error)
}
type relatableAdapter[D models4linkage.Relatable] struct {
	verifyItem func(ctx context.Context, tx dal.ReadTransaction, recordRef models4linkage.TeamModuleDocRef) (err error)
}

func (v relatableAdapter[D]) VerifyItem(ctx context.Context, tx dal.ReadTransaction, recordRef models4linkage.TeamModuleDocRef) (err error) {
	return v.verifyItem(ctx, tx, recordRef)
}

func NewRelatableAdapter[D models4linkage.Relatable](
	verifyItem func(ctx context.Context, tx dal.ReadTransaction, recordRef models4linkage.TeamModuleDocRef) (err error),
) RelatableAdapter[D] {
	return relatableAdapter[D]{
		verifyItem: verifyItem,
	}
}

//func (relatableAdapter[D]) GetRecord(ctx context.Context, tx dal.ReadTransaction, recordRef models4linkage.TeamModuleDocRef) (record.DataWithID[string, D], error) {
//	return nil, nil
//}

// SetRelated updates related records to define relationships
func SetRelated[D models4linkage.Relatable](
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userID string,
	now time.Time,
	adapter RelatableAdapter[D],
	object record.DataWithID[string, D],
	objectRef models4linkage.TeamModuleDocRef,
	related models4linkage.RelatedByTeamID,

) (updates []dal.Update, err error) {

	if err = objectRef.Validate(); err != nil {
		return nil, fmt.Errorf("facade4linkage.SetRelated got invalid argument `objectRef models4linkage.TeamModuleDocRef`: %w", err)
	}

	var updatedFields []string

	var relUpdates []dal.Update

	//var userContactID string
	//if userContactID, err = facade4userus.GetUserTeamContactID(ctx, tx, params.UserID, params.ContactusTeamWorkerParams.TeamModuleEntry); err != nil {
	//	return fmt.Errorf("failed to get user's contact ID: %w", err)
	//}

	for teamID, relatedByModuleID := range related {
		if teamID != objectRef.TeamID {
			return nil, validation.NewBadRequestError(errors.New("adding related item from other team is not supported yet"))
		}
		for moduleID, relatedByCollection := range relatedByModuleID {
			for collection, relatedByItemID := range relatedByCollection {
				for itemID, relatedItem := range relatedByItemID {
					itemRef := models4linkage.TeamModuleDocRef{
						TeamID:     teamID,
						ModuleID:   moduleID,
						Collection: collection,
						ItemID:     itemID,
					}
					if err := adapter.VerifyItem(ctx, tx, itemRef); err != nil {
						return nil, fmt.Errorf("failed to verify related item: %w", err)
					}
					objectWithRelated := object.Data.GetRelated()
					if objectWithRelated.Related == nil {
						objectWithRelated.Related = make(models4linkage.RelatedByTeamID, 1)
					}

					if relUpdates, err = objectWithRelated.SetRelationshipsToItem(
						userID,
						objectRef,
						models4linkage.TeamModuleDocRef{
							ModuleID:   moduleID,
							TeamID:     teamID,
							Collection: collection,
							ItemID:     itemID,
						},
						relatedItem.RelatedAs,
						relatedItem.RelatesAs,
						now,
					); err != nil {
						return updates, err
					}
					updates = append(updates, relUpdates...)
					for _, update := range relUpdates {
						if !slice.Contains(updatedFields, update.Field) {
							updatedFields = append(updatedFields, update.Field)
						}
					}
				}
			}
		}
	}
	return updates, err
}
