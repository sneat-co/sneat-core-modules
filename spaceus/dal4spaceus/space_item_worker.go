package dal4spaceus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-core/coretypes"

	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// SpaceItemDbo is an interface for space item DB data object. Examples: ContactDbo, ListDbo, etc.
type SpaceItemDbo = interface {
	Validate() error
}

// SpaceItemRequest DTO

// SliceIndexes DTO
type SliceIndexes struct {
	Start int
	End   int
}

// Brief is an interface for a brief of a DB data object
type Brief = interface {
	Validate()
}

// BriefsAdapter is an interface for a briefs adapter
type BriefsAdapter[ModuleDbo SpaceModuleDbo] interface {
	DeleteBrief(space ModuleDbo, id string) ([]update.Update, error)
	GetBriefsCount(space ModuleDbo) int
}

type mapBriefsAdapter[ModuleDbo SpaceModuleDbo] struct {
	getBriefsCount func(space ModuleDbo) int
	deleteBrief    func(space ModuleDbo, id string) ([]update.Update, error)
}

func (v mapBriefsAdapter[ModuleDbo]) DeleteBrief(spaceModuleDbo ModuleDbo, id string) ([]update.Update, error) {
	return v.deleteBrief(spaceModuleDbo, id)
}

func (v mapBriefsAdapter[ModuleDbo]) GetBriefsCount(spaceModuleDbo ModuleDbo) int {
	return v.getBriefsCount(spaceModuleDbo)
}

func NewMapBriefsAdapter[ModuleDbo SpaceModuleDbo](
	getBriefsCount func(spaceModuleDbo ModuleDbo) int,
	deleteBrief func(spaceModuleDbo ModuleDbo, id string) ([]update.Update, error),
) BriefsAdapter[ModuleDbo] {
	return mapBriefsAdapter[ModuleDbo]{
		getBriefsCount: getBriefsCount,
		deleteBrief:    deleteBrief,
	}
}

// SpaceItemWorkerParams defines params for space item worker
type SpaceItemWorkerParams[ModuleDbo SpaceModuleDbo, ItemDbo SpaceItemDbo] struct {
	*ModuleSpaceWorkerParams[ModuleDbo]
	SpaceItem        record.DataWithID[string, ItemDbo]
	SpaceItemUpdates []update.Update
}

// RunSpaceItemWorker runs space item worker
func RunSpaceItemWorker[ModuleDbo SpaceModuleDbo, ItemDbo SpaceItemDbo](
	ctx facade.ContextWithUser,
	request dto4spaceus.SpaceItemRequest,
	moduleID coretypes.ModuleID,
	spaceModuleDbo ModuleDbo,
	spaceItemCollection string,
	spaceItemDbo ItemDbo,
	worker func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *SpaceItemWorkerParams[ModuleDbo, ItemDbo]) (err error),
) (err error) {
	return RunModuleSpaceWorkerWithUserCtx(ctx, request.SpaceID, moduleID, spaceModuleDbo,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, moduleSpaceWorkerParams *ModuleSpaceWorkerParams[ModuleDbo]) (err error) {
			spaceItemKey := dal.NewKeyWithParentAndID(moduleSpaceWorkerParams.SpaceModuleEntry.Key, spaceItemCollection, request.ID)
			params := SpaceItemWorkerParams[ModuleDbo, ItemDbo]{
				ModuleSpaceWorkerParams: moduleSpaceWorkerParams,
				SpaceItem:               record.NewDataWithID(request.ID, spaceItemKey, spaceItemDbo),
			}
			if err = worker(ctx, tx, &params); err != nil {
				return err
			}
			if len(params.SpaceItemUpdates) > 0 {
				if err = tx.Update(ctx, spaceItemKey, params.SpaceItemUpdates); err != nil {
					return fmt.Errorf("failed to update space item record: %w", err)
				}
			}
			return nil
		},
	)
}

// DeleteSpaceItem deletes space item
func DeleteSpaceItem[ModuleDbo SpaceModuleDbo, ItemDbo SpaceItemDbo](
	ctx facade.ContextWithUser,
	request dto4spaceus.SpaceItemRequest,
	moduleID coretypes.ModuleID,
	moduleData ModuleDbo,
	spaceItemCollection string,
	spaceItemDbo ItemDbo,
	briefsAdapter BriefsAdapter[ModuleDbo],
	worker func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *SpaceItemWorkerParams[ModuleDbo, ItemDbo]) (err error),
) (err error) {
	return RunSpaceItemWorker(ctx, request, moduleID, moduleData, spaceItemCollection, spaceItemDbo,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, spaceItemWorkerParams *SpaceItemWorkerParams[ModuleDbo, ItemDbo]) (err error) {
			return deleteSpaceItemTxWorker[ModuleDbo](ctx, tx, spaceItemWorkerParams, briefsAdapter, worker)
		},
	)
}

func deleteSpaceItemTxWorker[ModuleDbo SpaceModuleDbo, ItemDbo SpaceItemDbo](
	ctx facade.ContextWithUser,
	tx dal.ReadwriteTransaction,
	params *SpaceItemWorkerParams[ModuleDbo, ItemDbo],
	briefsAdapter BriefsAdapter[ModuleDbo],
	worker func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, p *SpaceItemWorkerParams[ModuleDbo, ItemDbo]) error,
) (err error) {
	if err = tx.Get(ctx, params.Space.Record); err != nil {
		return
	}
	if err = tx.Get(ctx, params.SpaceItem.Record); err != nil && !dal.IsNotFound(err) {
		return err
	}
	if err = tx.Get(ctx, params.SpaceModuleEntry.Record); err != nil {
		return
	}
	if worker != nil {
		if err = worker(ctx, tx, params); err != nil {
			return fmt.Errorf("failed to execute spaceItemWorker: %w", err)
		}
	}
	//if err = decrementCounter(&params); err != nil {
	//	return err
	//}
	if len(params.SpaceUpdates) > 0 {
		if err = TxUpdateSpace(ctx, tx, params.Started, params.Space, params.SpaceUpdates); err != nil {
			return fmt.Errorf("failed to update space record: %w", err)
		}
	}
	var spaceModuleUpdates []update.Update
	if spaceModuleUpdates, err = deleteBrief[ModuleDbo](params.SpaceModuleEntry, params.SpaceItem.ID, briefsAdapter, params.SpaceModuleUpdates); err != nil {
		return err
	} else {
		params.AddSpaceModuleUpdates(spaceModuleUpdates...)
	}

	if params.SpaceItem.Record.Exists() {
		if err = tx.Delete(ctx, params.SpaceItem.Key); err != nil {
			return fmt.Errorf("failed to delete space item record by key=%v: %w", params.SpaceItem.Key, err)
		}
	}
	return err
}

func deleteBrief[D SpaceModuleDbo](spaceModuleEntry record.DataWithID[coretypes.ModuleID, D], itemID string, adapter BriefsAdapter[D], updates []update.Update) ([]update.Update, error) {
	if adapter == nil {
		return updates, nil
	}
	return adapter.DeleteBrief(spaceModuleEntry.Data, itemID)
}
