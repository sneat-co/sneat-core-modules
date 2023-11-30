package facade4schedulus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/schedulus/dto4schedulus"
)

func UpdateSlot(ctx context.Context, userID string, request dto4schedulus.HappeningSlotRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	var worker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *happeningWorkerParams) (err error) {
		//teamKey := models4teamus.NewTeamKey(request.TeamID)
		//teamDto := new(models4teamus.TeamDto)
		//teamRecord := dal.NewRecordWithData(teamKey, teamDto)
		//
		//if err = tx.Get(ctx, teamRecord); err != nil {
		//	return nil, fmt.Errorf("failed to get team record: %w", err)
		//}

		if params.Happening.Dto.Type == "single" {
			params.Happening.Dto.Slots[0] = &request.Slot
			params.HappeningUpdates = []dal.Update{
				{
					Field: "slots",
					Value: params.Happening.Dto.Slots,
				},
			}
		}
		return
	}

	if err = modifyHappening(ctx, userID, request.HappeningRequest, worker); err != nil {
		return err
	}
	return nil
}
