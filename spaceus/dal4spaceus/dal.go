package dal4spaceus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/facade/db"
	"time"
)

var txUpdate = func(ctx context.Context, tx dal.ReadwriteTransaction, key *dal.Key, data []update.Update, opts ...dal.Precondition) error {
	return db.TxUpdate(ctx, tx, key, data, opts...)
}

func txUpdateSpace(ctx context.Context, tx dal.ReadwriteTransaction, timestamp time.Time, space dbo4spaceus.SpaceEntry, data []update.Update, opts ...dal.Precondition) error {
	if err := space.Data.Validate(); err != nil {
		return fmt.Errorf("space record is not valid: %w", err)
	}
	space.Data.Version++
	data = append(data,
		update.ByFieldName("v", space.Data.Version),
		update.ByFieldName("timestamp", timestamp),
	)
	return txUpdate(ctx, tx, space.Key, data, opts...)
}

func txUpdateSpaceModule[D SpaceModuleDbo](ctx context.Context, tx dal.ReadwriteTransaction, _ time.Time, spaceModule record.DataWithID[string, D], data []update.Update, opts ...dal.Precondition) error {
	if !spaceModule.Record.Exists() {
		return fmt.Errorf("an attempt to update a space module record that does not exist: %s", spaceModule.Key.String())
	}
	return txUpdate(ctx, tx, spaceModule.Key, data, opts...)
}
