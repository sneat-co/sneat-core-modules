package dal4schedulus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/schedulus/const4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/dto4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/models4schedulus"
	"github.com/sneat-co/sneat-core-modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type SchedulusTeamWorkerParams = dal4teamus.ModuleTeamWorkerParams[*models4schedulus.SchedulusTeamDto]

func NewSchedulusTeamWorkerParams(userID, teamID string) *SchedulusTeamWorkerParams {
	teamWorkerParams := dal4teamus.NewTeamWorkerParams(userID, teamID)
	return dal4teamus.NewTeamModuleWorkerParams(const4schedulus.ModuleID, teamWorkerParams, new(models4schedulus.SchedulusTeamDto))
}

type HappeningWorkerParams struct {
	SchedulusTeamWorkerParams
	Happening models4schedulus.HappeningContext
}

func RunHappeningTeamWorker(
	ctx context.Context,
	user facade.User,
	request dto4schedulus.HappeningRequest,
	moduleID string,
	happeningWorker func(ctx context.Context, tx dal.ReadwriteTransaction, params *HappeningWorkerParams) (err error),
) (err error) {
	schedulusTeamDto := new(models4schedulus.SchedulusTeamDto)

	moduleTeamWorker := func(
		ctx context.Context,
		tx dal.ReadwriteTransaction,
		moduleTeamParams *dal4teamus.ModuleTeamWorkerParams[*models4schedulus.SchedulusTeamDto],
	) (err error) {
		params := &HappeningWorkerParams{
			SchedulusTeamWorkerParams: *moduleTeamParams,
			Happening:                 models4schedulus.NewHappeningContext(request.TeamID, request.HappeningID),
		}
		if err = tx.Get(ctx, params.Happening.Record); err != nil {
			if dal.IsNotFound(err) {
				params.Happening.Dto.Type = request.HappeningType
			} else {
				return fmt.Errorf("failed to get happening: %w", err)
			}
		}

		return happeningWorker(ctx, tx, params)
	}
	return dal4teamus.RunModuleTeamWorker(ctx, user, request.TeamRequest, moduleID, schedulusTeamDto, moduleTeamWorker)
}
