package api4invitus

import (
	"github.com/sneat-co/sneat-core-modules/invitus/facade4invitus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var createMassInviteFunc = facade4invitus.CreateMassInvite

// httpPostCreateMassInvite is an API endpoint to create a mass-invite
func httpPostCreateMassInvite(w http.ResponseWriter, r *http.Request) {

	var request facade4invitus.CreateMassInviteRequest
	ctx, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}

	var response facade4invitus.CreateInviteResponse
	response, err = createMassInviteFunc(ctx, request)

	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, response)
}
