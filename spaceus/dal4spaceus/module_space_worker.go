package dal4spaceus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	"slices"
)

// ModuleSpaceWorkerParams passes data to a space worker
type ModuleSpaceWorkerParams[D SpaceModuleDbo] struct {
	*SpaceWorkerParams
	SpaceModuleEntry   record.DataWithID[string, D]
	SpaceModuleUpdates []dal.Update
}

func (v *ModuleSpaceWorkerParams[D]) AddSpaceModuleUpdates(updates ...dal.Update) {
	if len(updates) > 0 {
		v.SpaceModuleUpdates = append(v.SpaceModuleUpdates, updates...)
		v.SpaceModuleEntry.Record.MarkAsChanged()
	}
}

func (v *ModuleSpaceWorkerParams[D]) GetRecords(ctx context.Context, tx dal.ReadSession, records ...dal.Record) error {
	return v.SpaceWorkerParams.GetRecords(ctx, tx, append(records, v.SpaceModuleEntry.Record)...)
}

type SpaceModuleDbo = interface {
	Validate() error
}

func RunModuleSpaceWorkerNoUpdates[D SpaceModuleDbo](
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userCtx facade.UserContext,
	spaceID, moduleID string,
	data D,
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, spaceWorkerParams *ModuleSpaceWorkerParams[D]) (err error),
) (err error) {
	if err = validateRunModuleSpaceWorkerArgs[D](spaceID, moduleID, data); err != nil {
		return err
	}
	if worker == nil {
		panic("worker is nil")
	}
	spaceWorkerParams := NewSpaceWorkerParams(userCtx, spaceID)
	params := NewSpaceModuleWorkerParams(moduleID, spaceWorkerParams, data)
	return worker(ctx, tx, params)
}

func NewSpaceModuleWorkerParams[D SpaceModuleDbo](
	moduleID string,
	spaceWorkerParams *SpaceWorkerParams,
	data D,
) *ModuleSpaceWorkerParams[D] {
	return &ModuleSpaceWorkerParams[D]{
		SpaceWorkerParams: spaceWorkerParams,
		SpaceModuleEntry: record.NewDataWithID(moduleID,
			dal.NewKeyWithParentAndID(spaceWorkerParams.Space.Key, dbo4spaceus.SpaceModulesCollection, moduleID),
			data,
		),
	}
}

func runModuleSpaceWorkerReadonlyTx[D SpaceModuleDbo](
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *ModuleSpaceWorkerParams[D],
	worker func(ctx context.Context, tx dal.ReadTransaction, spaceWorkerParams *ModuleSpaceWorkerParams[D]) (err error),
) (err error) {
	if err = worker(ctx, tx, params); err != nil {
		return fmt.Errorf("failed to execute module space worker inside runModuleSpaceWorkerReadonlyTx: %w", err)
	}
	return nil
}

func runModuleSpaceWorkerReadwriteTx[D SpaceModuleDbo](
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *ModuleSpaceWorkerParams[D],
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, spaceWorkerParams *ModuleSpaceWorkerParams[D]) (err error),
) (err error) {
	if worker == nil {
		panic("worker is nil")
	}
	if err = worker(ctx, tx, params); err != nil {
		return fmt.Errorf("failed to execute module space worker inside runModuleSpaceWorkerReadwriteTx: %w", err)
	}
	if err = applySpaceModuleUpdates(ctx, tx, params); err != nil {
		return fmt.Errorf("space module worker failed to apply space module record updates: %w", err)
	}

	// This should be called in runSpaceWorkerTx() instead
	//if err = applySpaceUpdates(ctx, tx, params.SpaceWorkerParams); err != nil {
	//	return fmt.Errorf("space module worker failed to apply space record updates: %w", err)
	//}

	return nil
}

func RunReadonlyModuleSpaceWorker[D SpaceModuleDbo](
	ctx context.Context,
	userCtx facade.UserContext,
	request dto4spaceus.SpaceRequest,
	moduleID string,
	data D,
	worker func(ctx context.Context, tx dal.ReadTransaction, spaceWorkerParams *ModuleSpaceWorkerParams[D]) (err error),
) (err error) {
	spaceWorkerParams := NewSpaceWorkerParams(userCtx, request.SpaceID)
	params := NewSpaceModuleWorkerParams(moduleID, spaceWorkerParams, data)

	return facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		return runModuleSpaceWorkerReadonlyTx(ctx, tx, params, worker)
	})
}

func RunModuleSpaceWorkerWithUserCtx[D SpaceModuleDbo](
	ctx context.Context,
	userCtx facade.UserContext,
	spaceID, moduleID string,
	data D,
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, spaceWorkerParams *ModuleSpaceWorkerParams[D]) (err error),
) (err error) {
	//const singleCall = "singleCall"
	//if ctx.Value(singleCall) != nil {
	//	panic("duplicate call")
	//}
	//ctx = context.WithValue(ctx, singleCall, true)

	spaceWorkerParams := NewSpaceWorkerParams(userCtx, spaceID)
	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return
	}
	if err = db.Get(ctx, spaceWorkerParams.Space.Record); err != nil {
		return fmt.Errorf("failed to get space record outside of transaction: %w", err)
	}
	if userCtx != nil {
		if userID := userCtx.GetUserID(); userID != "" {
			if !slices.Contains(spaceWorkerParams.Space.Data.UserIDs, userID) {
				return fmt.Errorf("user is not a member of the space")
			}
		}
	}
	return facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		moduleWorkerParams := NewSpaceModuleWorkerParams(moduleID, spaceWorkerParams, data)
		if err = runSpaceWorkerTx(ctx, tx, spaceWorkerParams, nil,
			func(ctx context.Context, tx dal.ReadwriteTransaction, spaceWorkerParams *SpaceWorkerParams) (err error) {
				if err = runModuleSpaceWorkerReadwriteTx(ctx, tx, moduleWorkerParams, worker); err != nil {
					return fmt.Errorf("failed in runModuleSpaceWorkerReadwriteTx(): %w", err)
				}
				return nil
			}); err != nil {
			return fmt.Errorf("failed in runSpaceWorkerTx(): %w", err)
		}
		//logus.Debugf(ctx, "RunModuleSpaceWorkerWithUserCtx() completed, err = nil")
		return
	})
}

func RunModuleSpaceWorkerTx[D SpaceModuleDbo](
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userCtx facade.UserContext,
	spaceID, moduleID string,
	data D,
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, spaceWorkerParams *ModuleSpaceWorkerParams[D]) (err error),
) (err error) {
	if err = validateRunModuleSpaceWorkerArgs[D](spaceID, moduleID, data); err != nil {
		return err
	}
	spaceWorkerParams := NewSpaceWorkerParams(userCtx, spaceID)
	params := NewSpaceModuleWorkerParams(moduleID, spaceWorkerParams, data)
	return runModuleSpaceWorkerReadwriteTx(ctx, tx, params, worker)
}

func validateRunModuleSpaceWorkerArgs[D SpaceModuleDbo](spaceID, moduleID string, data D) error {
	var d any
	if d = data; d == nil {
		panic("data is nil")
	}
	if moduleID == "" {
		panic("moduleID is empty")
	}
	if spaceID == "" {
		return fmt.Errorf("spaceID is empty")
	}
	return nil
}
