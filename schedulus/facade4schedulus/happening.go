package facade4schedulus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/schedulus/dal4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/dto4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/models4schedulus"
	"github.com/sneat-co/sneat-go-core/facade"
	"log"
)

type happeningWorkerParams struct {
	*dal4schedulus.SchedulusTeamWorkerParams
	Happening        models4schedulus.HappeningContext
	HappeningUpdates []dal.Update
}

type happeningWorker = func(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	param *happeningWorkerParams,
) (err error)

func modifyHappening(ctx context.Context, userID string, request dto4schedulus.HappeningRequest, worker happeningWorker) (err error) {
	if userID == "" {
		return fmt.Errorf("not allowed to call without userID: %w", facade.ErrUnauthorized)
	}
	if err = request.Validate(); err != nil {
		return
	}
	db := facade.GetDatabase(ctx)
	err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		params := happeningWorkerParams{
			SchedulusTeamWorkerParams: dal4schedulus.NewSchedulusTeamWorkerParams(userID, request.TeamID),
			Happening:                 models4schedulus.NewHappeningContext(request.TeamID, request.HappeningID),
		}
		if err = worker(ctx, tx, &params); err != nil {
			return fmt.Errorf("failed in happening worker: %w", err)
		}
		if len(params.HappeningUpdates) > 0 {
			if err = params.Happening.Dto.Validate(); err != nil {
				return fmt.Errorf("happening record is not valid after running worker: %w", err)
			}
			log.Printf("updating happening: %s", params.Happening.Key)
			if err = tx.Update(ctx, params.Happening.Key, params.HappeningUpdates); err != nil {
				return fmt.Errorf("failed to update happening record: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update happening in transaction: %w", err)
	}
	return err
}
