package facade4linkage

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"slices"
)

func updateUserWithRelatedTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userID string,
	users map[string]dbo4userus.UserEntry,
	itemRef dbo4linkage.SpaceModuleItemRef,
	relatedItem dbo4linkage.RelatedItem,
) (err error) {
	if users == nil {
		panic("users == nil")
	}

	var user dbo4userus.UserEntry
	var ok bool

	if user, ok = users[userID]; !ok {
		user := dbo4userus.NewUserEntry(userID)
		if err = tx.Get(ctx, user.Record); err != nil {
			return err
		}
		users[userID] = user
	}

	if slices.Contains(user.Data.SpaceIDs, string(itemRef.Space)) {
		return nil
	}

	userRelated := dbo4linkage.GetRelatedItemByRef(user.Data.Related, itemRef, true)

	var updates []update.Update

	for roleID, role := range relatedItem.RolesToItem {
		if userRelated.RolesOfItem[roleID] != role {
			userRelated.RolesOfItem[roleID] = role
			updates = append(updates, update.ByFieldName("related."+itemRef.ID()+".rolesOfItem."+roleID, role))
		}
	}

	return tx.Update(ctx, user.Record.Key(), updates)
}
