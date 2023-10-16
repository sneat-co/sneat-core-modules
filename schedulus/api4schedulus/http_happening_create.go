package api4schedulus

import (
	"github.com/sneat-co/sneat-core-modules/schedulus/dto4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/facade4schedulus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var createHappening = facade4schedulus.CreateHappening

// httpPostCreateHappening creates recurring happening
func httpPostCreateHappening(w http.ResponseWriter, r *http.Request) {
	var request dto4schedulus.CreateHappeningRequest
	ctx, userContext, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	response, err := createHappening(ctx, userContext, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, &response)
}
