package dal4schedulus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/schedulus/models4schedulus"
	"github.com/sneat-co/sneat-core-modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-core-modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
)

func RunSchedulusTeamWorker(
	ctx context.Context,
	user facade.User,
	request dto4teamus.TeamRequest,
	worker func(ctx context.Context, tx dal.ReadwriteTransaction, params *SchedulusTeamWorkerParams) (err error),
) error {
	return dal4teamus.RunModuleTeamWorker(ctx, user, request, const4contactus.ModuleID, new(models4schedulus.SchedulusTeamDto), worker)
}
