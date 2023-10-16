package api4schedulus

import (
	"github.com/sneat-co/sneat-core-modules/schedulus/facade4schedulus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var deleteHappening = facade4schedulus.DeleteHappening

// httpDeleteHappening deletes happening
func httpDeleteHappening(w http.ResponseWriter, r *http.Request) {
	var request = getHappeningRequestParamsFromURL(r)
	ctx, userContext, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	err = deleteHappening(ctx, userContext, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, nil)
}
