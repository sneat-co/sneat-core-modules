package facade4schedulus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/schedulus/dal4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/dto4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/models4schedulus"
	"github.com/sneat-co/sneat-go-core/facade"
	"log"
)

type happeningWorkerParams struct {
	*dal4contactus.ContactusTeamWorkerParams
	SchedulusTeam        record.DataWithID[string, *models4schedulus.SchedulusTeamDto]
	SchedulusTeamUpdates []dal.Update
	Happening            models4schedulus.HappeningContext
	HappeningUpdates     []dal.Update
}

type happeningWorker = func(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	param happeningWorkerParams,
) (err error)

func modifyHappening(ctx context.Context, userID string, request dto4schedulus.HappeningRequest, worker happeningWorker) (err error) {
	if userID == "" {
		return fmt.Errorf("not allowed to delete happening: %w", facade.ErrUnauthorized)
	}

	db := facade.GetDatabase(ctx)
	err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		params := happeningWorkerParams{
			ContactusTeamWorkerParams: dal4contactus.NewContactusTeamWorkerParams(userID, request.TeamID),
			SchedulusTeam:             dal4schedulus.NewSchedulusTeamContext(request.TeamID),
			Happening:                 models4schedulus.NewHappeningContext(request.TeamID, request.HappeningID),
		}
		if err = tx.Get(ctx, params.Happening.Record); err != nil {
			return fmt.Errorf("failed to get happening by ContactID=%v: %w", params.Happening.ID, err)
		}
		log.Printf("happening: %+v", *params.Happening.Dto)
		if err = worker(ctx, tx, params); err != nil {
			return fmt.Errorf("failed in happening worker: %w", err)
		}
		if len(params.HappeningUpdates) > 0 {
			if err = params.Happening.Dto.Validate(); err != nil {
				return fmt.Errorf("happening record is not valid after running worker: %w", err)
			}
			if err = tx.Update(ctx, params.Happening.Key, params.HappeningUpdates); err != nil {
				return fmt.Errorf("failed to update happening record: %w", err)
			}
		}
		return err
	})
	if err != nil {
		return fmt.Errorf("failed in transaction: %w", err)
	}
	return err
}
