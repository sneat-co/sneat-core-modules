package facade4linkage

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
)

func updateUserRelated( // TODO: Document use case when this is needed
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	item record.DataWithID[string, *dbo4linkage.WithRelatedAndIDsAndUserID],
) (
	userUpdates record.Updates, err error,
) {
	user := dbo4userus.NewUserEntry(item.Data.UserID)
	if err = tx.Get(ctx, user.Record); err != nil {
		return
	}

	return
}
