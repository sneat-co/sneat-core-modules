package dal4contactus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dbo4contactus"
	dal4spaceus2 "github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type ContactusSpaceWorkerParams = dal4spaceus2.ModuleSpaceWorkerParams[*dbo4contactus.ContactusSpaceDbo]

func NewContactusSpaceWorkerParams(userCtx facade.UserContext, spaceID string) *ContactusSpaceWorkerParams {
	teamWorkerParams := dal4spaceus2.NewSpaceWorkerParams(userCtx, spaceID)
	return dal4spaceus2.NewSpaceModuleWorkerParams(const4contactus.ModuleID, teamWorkerParams, new(dbo4contactus.ContactusSpaceDbo))
}

func RunReadonlyContactusSpaceWorker(
	ctx context.Context,
	userCtx facade.UserContext,
	request dto4spaceus.SpaceRequest,
	worker func(ctx context.Context, tx dal.ReadTransaction, params *ContactusSpaceWorkerParams) (err error),
) error {
	return dal4spaceus2.RunReadonlyModuleSpaceWorker(ctx, userCtx, request, const4contactus.ModuleID, new(dbo4contactus.ContactusSpaceDbo), worker)
}

type ContactusModuleWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *ContactusSpaceWorkerParams) (err error)

func RunContactusSpaceWorker(
	ctx context.Context,
	userCtx facade.UserContext,
	request dto4spaceus.SpaceRequest,
	worker ContactusModuleWorker,
) error {
	return dal4spaceus2.RunModuleSpaceWorkerWithUserCtx(ctx, userCtx, request.SpaceID, const4contactus.ModuleID, new(dbo4contactus.ContactusSpaceDbo), worker)
}

func RunContactusSpaceWorkerTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userCtx facade.UserContext,
	spaceID string,
	worker ContactusModuleWorker,
) error {
	return dal4spaceus2.RunModuleSpaceWorkerTx(ctx, tx, userCtx, spaceID, const4contactus.ModuleID, new(dbo4contactus.ContactusSpaceDbo), worker)
}

func RunContactusSpaceWorkerNoUpdate(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userCtx facade.UserContext,
	spaceID string,
	worker ContactusModuleWorker,
) error {
	return dal4spaceus2.RunModuleSpaceWorkerNoUpdates(ctx, tx, userCtx, spaceID, const4contactus.ModuleID, new(dbo4contactus.ContactusSpaceDbo), worker)
}
