package dal4userus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
)

type UserExtDbo = interface {
	Validate() error
}

type UserExtWorkerParams[T any] struct {
	UserExt        record.DataWithID[string, *T]
	UserExtUpdates []update.Update
}

func RunUserExtWorker[T any](
	ctx context.Context,
	userID string,
	extID coretypes.ModuleID,
	userExtDbo *T,
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, param *UserExtWorkerParams[T]) error,
) error {
	return facade.RunReadwriteTransaction(ctx,
		func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
			params := UserExtWorkerParams[T]{
				UserExt: record.NewDataWithID(userID, NewUserExtKey(userID, extID), userExtDbo),
			}
			if err = worker(ctx, tx, &params); err != nil {
				return fmt.Errorf("failed to execute user ext worker: %w", err)
			}
			if params.UserExt.Record.Exists() {
				if len(params.UserExtUpdates) > 0 {
					if params.UserExt.Record.HasChanged() {
						return fmt.Errorf("len(params.UserExtUpdates) > 0 but params.UserExt.Record.HasChanged() == false")
					}
					if err = tx.Update(ctx, params.UserExt.Key, params.UserExtUpdates); err != nil {
						return fmt.Errorf("failed to update user's extension record: %w", err)
					}
				}
			} else if err = tx.Insert(ctx, params.UserExt.Record); err != nil {
				return fmt.Errorf("failed to insert user's extension record: %w", err)
			}
			return err
		},
	)
}
