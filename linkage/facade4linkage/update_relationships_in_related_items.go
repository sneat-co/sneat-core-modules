package facade4linkage

import (
	"time"

	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

// UpdateRelationshipsInRelatedItems TODO: This should be a generic function in the linkage module
func UpdateRelationshipsInRelatedItems(
	now time.Time,
	userID string,
	spaceID coretypes.SpaceID,
	dbo *dbo4linkage.WithRelatedAndIDs,
	related dbo4linkage.RelatedModules,
) (err error) {
	// TODO: either use the parameters or remove them or document why not using but needs keeping
	//_, _ = ctx, tx

	if len(related) == 0 {
		return nil // nothing to do
	}

	for moduleID, relatedCollections := range related {
		for collection, relatedItems := range relatedCollections {
			for itemID, relatedItem := range relatedItems {
				itemRef := dbo4linkage.ItemRef{
					ExtID:      coretypes.ExtID(moduleID),
					Collection: collection,
					ItemID:     itemID,
				}

				command := dbo4linkage.RelationshipItemRolesCommand{
					ItemRef: itemRef,
				}
				if len(relatedItem.RolesOfItem) > 0 {
					if command.Add == nil {
						command.Add = &dbo4linkage.RolesCommand{}
					}
					for roleID := range relatedItem.RolesOfItem {
						command.Add.RolesOfItem = append(command.Add.RolesOfItem, roleID)
					}
				}
				if len(relatedItem.RolesToItem) > 0 {
					if command.Add == nil {
						command.Add = &dbo4linkage.RolesCommand{}
					}
					for roleID := range relatedItem.RolesToItem {
						command.Add.RolesToItem = append(command.Add.RolesOfItem, roleID)
					}
				}
				if _, err = dbo.AddRelationshipAndID(now, userID, spaceID, command); err != nil {
					return err
				}
			}
		}
	}
	dbo4linkage.UpdateRelatedIDs(spaceID, &dbo.WithRelated, &dbo.WithRelatedIDs)
	return nil
}
