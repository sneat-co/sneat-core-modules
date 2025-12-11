package facade4contactus

import (
	"context"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-core/facade/db"
)

var txUpdate = func(ctx context.Context, tx dal.ReadwriteTransaction, key *dal.Key, data []update.Update, opts ...dal.Precondition) error {
	return db.TxUpdate(ctx, tx, key, data, opts...)
}
