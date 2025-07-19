package facade4linkage

//import (
//	"fmt"
//	"github.com/dal-go/dalgo/dal"
//	"github.com/dal-go/dalgo/record"
//	"github.com/dal-go/dalgo/update"
//	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
//	"github.com/sneat-co/sneat-core-modules/linkage/dto4linkage"
//	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
//	"github.com/sneat-co/sneat-go-core/coretypes"
//	"github.com/sneat-co/sneat-go-core/facade"
//	"github.com/strongo/validation"
//)
//
//func updateRelatedDbos(
//	ctx facade.ContextWithUser,
//	tx dal.ReadwriteTransaction,
//	request dto4linkage.UpdateRelatedRequest,
//) error {
//	type relatedItem struct {
//		item    record.DataWithID[string, SpaceItemDboWithRelatedAndIDs]
//		updates []update.Update
//	}
//	itemsWithUpdates := make([]relatedItem, 0, len(request.Related))
//
//	itemRef := dbo4linkage.NewItemRef(request.ExtensionID, request.Collection, request.ID)
//
//	for _, command := range request.Related {
//		targetRef := command.ItemRef
//		targetCommand := dbo4linkage.RelationshipItemRolesCommand{
//			ItemRef: itemRef,
//		}
//		if command.Add != nil {
//			targetCommand.Add = &dbo4linkage.RolesCommand{}
//			for _, role := range command.Add.RolesOfItem {
//				targetCommand.Add.RolesToItem = append(targetCommand.Add.RolesToItem, role)
//			}
//			for _, role := range command.Add.RolesToItem {
//				targetCommand.Add.RolesOfItem = append(targetCommand.Add.RolesOfItem, role)
//			}
//		}
//		item, updates, err := updateRelatedDbo(ctx, tx, request.SpaceID, targetRef, targetCommand)
//		if err != nil {
//			return err
//		}
//		if len(updates) > 0 {
//			itemsWithUpdates = append(itemsWithUpdates, relatedItem{
//				item:    item,
//				updates: updates,
//			})
//		}
//	}
//
//	for _, related := range itemsWithUpdates {
//		if err := tx.Update(ctx, related.item.Record.Key(), related.updates); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func updateRelatedDbo(
//	ctx facade.ContextWithUser,
//	tx dal.ReadTransaction,
//	spaceID coretypes.SpaceID,
//	targetRef dbo4linkage.ItemRef,
//	command dbo4linkage.RelationshipItemRolesCommand,
//) (
//	target record.DataWithID[string, SpaceItemDboWithRelatedAndIDs],
//	update []update.Update,
//	err error,
//) {
//	dboFactory := getDboFactory(targetRef.ExtID, targetRef.Collection)
//	if dboFactory == nil {
//		err = validation.NewBadRequestError(fmt.Errorf("unknown moduleID or collection: %s/%s", targetRef.ExtID, targetRef.Collection))
//		return
//	}
//
//	targetKey := dbo4spaceus.NewSpaceModuleItemKey(spaceID, targetRef.ExtID, targetRef.Collection, targetRef.ItemID)
//	targetDbo := dboFactory.NewItemDbo()
//	target = record.NewDataWithID[string, SpaceItemDboWithRelatedAndIDs](targetRef.ItemID, targetKey, targetDbo)
//	if err = tx.Get(ctx, target.Record); err != nil {
//		return
//	}
//	//updateRelatedTxWorker(ctx, tx)
//	_ = targetDbo.RelatedAndIDs()
//	return
//}
