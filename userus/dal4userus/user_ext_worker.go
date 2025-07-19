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

func RunUserExtWorkerWithUserContext[T any](
	ctx facade.ContextWithUser,
	extID coretypes.ExtID,
	userExtDbo *T,
	worker func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, param *UserExtWorkerParams[T]) error,
) error {
	user := ctx.User()
	userID := user.GetUserID()
	return RunUserExtWorker[T](ctx, userID, extID, userExtDbo,
		func(ctx context.Context, tx dal.ReadwriteTransaction, param *UserExtWorkerParams[T]) error {
			ctxWithUser := facade.NewContextWithUser(ctx, user)
			return worker(ctxWithUser, tx, param)
		},
	)
}

func RunUserExtWorker[T any](
	ctx context.Context,
	userID string,
	extID coretypes.ExtID,
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
			if params.UserExt.Record.HasChanged() {
				if params.UserExt.Record.Exists() {
					if len(params.UserExtUpdates) > 0 {
						if err = tx.Update(ctx, params.UserExt.Key, params.UserExtUpdates); err != nil {
							return fmt.Errorf("failed to update user's extension record: %w", err)
						}
					} else if err = tx.Set(ctx, params.UserExt.Record); err != nil {
						return fmt.Errorf("failed to rewrite user's extension record: %w", err)
					}
				} else if err = tx.Insert(ctx, params.UserExt.Record); err != nil {
					return fmt.Errorf("failed to insert user's extension record: %w", err)
				}
			} else if len(params.UserExtUpdates) > 0 {
				return fmt.Errorf("len(params.UserExtUpdates) > 0 but params.UserExt.Record.HasChanged() == false")
			}
			return err
		},
	)
}
