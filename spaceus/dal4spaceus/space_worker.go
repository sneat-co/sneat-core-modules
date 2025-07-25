package dal4spaceus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"slices"
	"strings"
	"time"
)

type spaceWorker = func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *SpaceWorkerParams) (err error)

func NewSpaceWorkerParams(userCtx facade.UserContext, spaceID coretypes.SpaceID) *SpaceWorkerParams {
	return &SpaceWorkerParams{
		UserCtx: userCtx,
		Space:   dbo4spaceus.NewSpaceEntry(spaceID),
		Started: time.Now(),
	}
}

// SpaceWorkerParams passes data to a space worker
type SpaceWorkerParams struct {
	UserCtx facade.UserContext // TODO: consider removing this field in favor of using facade.ContextWithUser
	Started time.Time
	//
	Space         dbo4spaceus.SpaceEntry
	SpaceUpdates  []update.Update
	RecordUpdates []record.Updates
}

func (v SpaceWorkerParams) UserID() string {
	if v.UserCtx == nil {
		return ""
	}
	return v.UserCtx.GetUserID()
}

func (v SpaceWorkerParams) hasSpaceRecord(records []dal.Record) bool {
	for _, rec := range records {
		if rec == v.Space.Record {
			return true
		}
	}
	return false
}

// GetRecords loads records from DB. If userID is passed, it will check for user permissions
func (v SpaceWorkerParams) GetRecords(ctx context.Context, tx dal.MultiGetter, records ...dal.Record) (err error) {
	userID := v.UserID()

	hasSpaceRecord := v.hasSpaceRecord(records)

	if userID != "" && !hasSpaceRecord {
		records = append(records, v.Space.Record)
		hasSpaceRecord = true
	}

	if err = tx.GetMulti(ctx, records); err != nil {
		return fmt.Errorf("failed in SpaceWorkerParams.GetMulti(len(records)=%d): %w", len(records), err)
	}

	if hasSpaceRecord {
		if err = v.Space.Data.Validate(); err != nil {
			return fmt.Errorf("space record loaded from DB is not valid: %w", err)
		}
	}

	if userID != "" {
		if !v.Space.Record.Exists() {
			return errors.New("space record does not exist")
		}
		if !slices.Contains(v.Space.Data.UserIDs, userID) {
			return fmt.Errorf("%w: space record has no current userID in UserIDs field: %s", facade.ErrUnauthorized, userID)
		}
	}
	return nil
}

// RunSpaceWorkerWithUserContext executes a space worker
var RunSpaceWorkerWithUserContext = func(ctx facade.ContextWithUser, spaceID coretypes.SpaceID, worker spaceWorker) (err error) {
	if strings.TrimSpace(string(spaceID)) == "" {
		return fmt.Errorf("required parameter `spaceID` of RunSpaceWorkerWithUserContext() is an empty string")
	}
	return runSpaceWorker(ctx, spaceID, worker)
}

// RunSpaceWorkerWithoutUserContext executes a space worker without user context
//var RunSpaceWorkerWithoutUserContext = func(ctx context.Context, spaceID coretypes.SpaceID, worker spaceWorker) (err error) {
//	if strings.TrimSpace(string(spaceID)) == "" {
//		return fmt.Errorf("required parameter `spaceID` of RunSpaceWorkerWithoutUserContext() is an empty string")
//	}
//	return runSpaceWorker(ctx, nil, spaceID, worker)
//}

var runSpaceWorker = func(ctx facade.ContextWithUser, spaceID coretypes.SpaceID, worker spaceWorker) (err error) {
	if strings.TrimSpace(string(spaceID)) == "" {
		return fmt.Errorf("required parameter `spaceID` of runSpaceWorker() is an empty string")
	}
	userCtx := ctx.User()
	return facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		ctxWithUser := facade.NewContextWithUser(ctx, userCtx)
		return RunSpaceWorkerTx(ctxWithUser, tx, spaceID, worker)
	})
}

func RunSpaceWorkerTx(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, spaceID coretypes.SpaceID, worker spaceWorker) (err error) {
	userCtx := ctx.User()
	params := NewSpaceWorkerParams(userCtx, spaceID)
	return runSpaceWorkerTx(ctx, tx, params, nil, worker)
}

func runSpaceWorkerTx(
	ctx facade.ContextWithUser,
	tx dal.ReadwriteTransaction,
	params *SpaceWorkerParams,
	beforeWorker func(ctx context.Context) error,
	worker spaceWorker,
) (err error) {
	if beforeWorker != nil {
		if err = beforeWorker(ctx); err != nil {
			return fmt.Errorf("failed to run beforeWorker: %w", err)
		}
	}
	if err = worker(ctx, tx, params); err != nil {
		return fmt.Errorf("failed to execute space worker: %w", err)
	}
	for i, rec := range params.RecordUpdates {
		if rec.Record == nil {
			panic(fmt.Sprintf("worker %v returned params.RecordUpdates[%d] == nil", worker, i))
		}
	}
	if err = applySpaceUpdates(ctx, tx, params); err != nil {
		return fmt.Errorf("space worker failed to apply space record updates: %w", err)
	}
	if err = applyRecordUpdates(ctx, tx, params.RecordUpdates); err != nil {
		return fmt.Errorf("space worker failed to apply record updates: %w", err)
	}
	return
}

func applyRecordUpdates(ctx context.Context, tx dal.ReadwriteTransaction, recordUpdates []record.Updates) error {
	for i, rec := range recordUpdates {
		if rec.Record == nil {
			panic(fmt.Sprintf("recordUpdates[%d] == nil", i))
		}
		key := rec.Record.Key()
		if err := tx.Update(ctx, key, rec.Updates); err != nil {
			updateFieldNames := make([]string, len(rec.Updates))
			for _, u := range rec.Updates {
				fieldName := u.FieldName()
				if fieldName == "" {
					fieldName = strings.Join(u.FieldPath(), ".")
				}
				updateFieldNames = append(updateFieldNames, fieldName)
			}
			return fmt.Errorf(
				"failed to apply record updates (key=%s, updateFieldNames: %s): %w",
				key, strings.Join(updateFieldNames, ","), err)
		}
	}
	return nil
}

func applySpaceUpdates(ctx context.Context, tx dal.ReadwriteTransaction, params *SpaceWorkerParams) (err error) {
	if len(params.SpaceUpdates) == 0 && !params.Space.Record.HasChanged() {
		return
	}
	if spaceRecErr := params.Space.Record.Error(); spaceRecErr != nil {
		return fmt.Errorf("an attempt to update a space record with an error: %w", spaceRecErr)
	}
	if !params.Space.Record.HasChanged() {
		return fmt.Errorf("space record should be marked as changed before applying updates")
	}
	if err = params.Space.Data.Validate(); err != nil {
		return fmt.Errorf("space record is not valid before applying %d space updates: %w", len(params.SpaceUpdates), err)
	}
	if !params.Space.Record.Exists() {
		return tx.Insert(ctx, params.Space.Record)
	} else if len(params.SpaceUpdates) == 0 {
		return tx.Set(ctx, params.Space.Record)
	} else if err = TxUpdateSpace(ctx, tx, params.Started, params.Space, params.SpaceUpdates); err != nil {
		return fmt.Errorf("failed to update space record: %w", err)
	}
	return
}

func applySpaceModuleUpdates[D SpaceModuleDbo](
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *ModuleSpaceWorkerParams[D],
) (err error) {
	if len(params.SpaceModuleUpdates) == 0 && !params.SpaceModuleEntry.Record.HasChanged() {
		return nil
	}
	if err = params.SpaceModuleEntry.Record.Error(); err != nil && !dal.IsNotFound(err) {
		return fmt.Errorf("an attempt to update a space module record that has an error: %w", err)
	}
	if !params.SpaceModuleEntry.Record.HasChanged() {
		return fmt.Errorf("space module record should be marked as changed before applying updates")
	}
	if err = params.SpaceModuleEntry.Data.Validate(); err != nil {
		return fmt.Errorf("space module record is not valid before applying space module updates: %w", err)
	}

	if params.SpaceModuleEntry.Record.Exists() {
		if len(params.SpaceModuleUpdates) == 0 {
			if err = tx.Set(ctx, params.SpaceModuleEntry.Record); err != nil {
				return fmt.Errorf("failed to set space module record: %w", err)
			}
		} else if err = txUpdateSpaceModule(ctx, tx, params.Started, params.SpaceModuleEntry, params.SpaceModuleUpdates); err != nil {
			return fmt.Errorf("failed to update space module record: %w", err)
		}
	} else if err = tx.Insert(ctx, params.SpaceModuleEntry.Record); err != nil {
		return fmt.Errorf("failed to insert space module record: %w", err)
	}
	return
}

// CreateSpaceItem creates a space item
func CreateSpaceItem[D SpaceModuleDbo](
	ctx facade.ContextWithUser,
	spaceRequest dto4spaceus.SpaceRequest,
	moduleID coretypes.ExtID,
	data D,
	worker func(
		ctx facade.ContextWithUser,
		tx dal.ReadwriteTransaction,
		spaceWorkerParams *ModuleSpaceWorkerParams[D],
	) (err error),
) (err error) {
	if worker == nil {
		panic("worker is nil")
	}
	if err = spaceRequest.Validate(); err != nil {
		return err
	}
	err = RunModuleSpaceWorkerWithUserCtx(ctx, spaceRequest.SpaceID, moduleID, data,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *ModuleSpaceWorkerParams[D]) (err error) {
			if err = worker(ctx, tx, params); err != nil {
				return fmt.Errorf("failed to execute space worker passed to CreateSpaceItem: %w", err)
			}
			if err = params.Space.Data.Validate(); err != nil {
				return fmt.Errorf("space record is not valid after performing worker: %w", err)
			}
			return
		})
	if err != nil {
		return fmt.Errorf("failed to create a space item: %w", err)
	}
	return nil
}
