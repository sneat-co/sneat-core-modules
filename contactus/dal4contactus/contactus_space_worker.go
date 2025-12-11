package dal4contactus

import (
	"context"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-go-core/coretypes"

	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type ContactusSpaceWorkerParams = dal4spaceus.ModuleSpaceWorkerParams[*dbo4contactus.ContactusSpaceDbo]

func NewContactusSpaceWorkerParams(userCtx facade.UserContext, spaceID coretypes.SpaceID) *ContactusSpaceWorkerParams {
	spaceWorkerParams := dal4spaceus.NewSpaceWorkerParams(userCtx, spaceID)
	return dal4spaceus.NewSpaceModuleWorkerParams(const4contactus.ExtensionID, spaceWorkerParams, new(dbo4contactus.ContactusSpaceDbo))
}

func RunReadonlyContactusSpaceWorker(
	ctx context.Context,
	userCtx facade.UserContext,
	request dto4spaceus.SpaceRequest,
	worker func(ctx context.Context, tx dal.ReadTransaction, params *ContactusSpaceWorkerParams) (err error),
) error {
	return dal4spaceus.RunReadonlyModuleSpaceWorker(ctx, userCtx, request, const4contactus.ExtensionID, new(dbo4contactus.ContactusSpaceDbo), worker)
}

type ContactusModuleWorker = func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *ContactusSpaceWorkerParams) (err error)

func RunContactusSpaceWorker(
	ctx facade.ContextWithUser,
	request dto4spaceus.SpaceRequest,
	worker ContactusModuleWorker,
) error {
	return dal4spaceus.RunModuleSpaceWorkerWithUserCtx(ctx, request.SpaceID, const4contactus.ExtensionID, new(dbo4contactus.ContactusSpaceDbo), worker)
}

func RunContactusSpaceWorkerTx(
	ctx facade.ContextWithUser,
	tx dal.ReadwriteTransaction,
	spaceID coretypes.SpaceID,
	worker ContactusModuleWorker,
) error {
	return dal4spaceus.RunModuleSpaceWorkerTx(ctx, tx, spaceID, const4contactus.ExtensionID, new(dbo4contactus.ContactusSpaceDbo), worker)
}

func RunContactusSpaceWorkerNoUpdate(
	ctx facade.ContextWithUser,
	tx dal.ReadwriteTransaction,
	spaceID coretypes.SpaceID,
	worker ContactusModuleWorker,
) error {
	return dal4spaceus.RunModuleSpaceWorkerNoUpdates(ctx, tx, spaceID, const4contactus.ExtensionID, new(dbo4contactus.ContactusSpaceDbo), worker)
}
