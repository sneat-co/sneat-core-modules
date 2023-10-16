package facade4schedulus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/schedulus/models4schedulus"
)

// GetByID returns RecurringHappeningDto record
func GetByID(ctx context.Context, getter dal.ReadSession, id string, dto models4schedulus.HappeningDto) (record dal.Record, err error) {
	record = dal.NewRecordWithData(models4schedulus.NewHappeningKey(id), dto)
	return record, getter.Get(ctx, record)
}

// GetForUpdate returns TeamIDs record in transaction
func GetForUpdate(ctx context.Context, tx dal.ReadwriteTransaction, id string, dto models4schedulus.HappeningDto) (record dal.Record, err error) {
	return GetByID(ctx, tx, id, dto)
}
