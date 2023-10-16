package api4schedulus

import (
	"github.com/sneat-co/sneat-core-modules/schedulus/dto4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/facade4schedulus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var deleteSlot = facade4schedulus.DeleteSlots

// httpDeleteHappening deletes happening
func httpDeleteSlots(w http.ResponseWriter, r *http.Request) {
	var request = dto4schedulus.DeleteHappeningSlotRequest{
		HappeningSlotRefRequest: dto4schedulus.HappeningSlotRefRequest{
			HappeningRequest: getHappeningRequestParamsFromURL(r),
		},
	}
	ctx, userContext, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	err = deleteSlot(ctx, userContext, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, nil)
}
