package api4memberus

import (
	"github.com/sneat-co/sneat-core-modules/memberus/facade4memberus"
	"github.com/sneat-co/sneat-core-modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var leaveTeam = facade4memberus.LeaveTeam

// httpPostLeaveTeam is an API endpoint that removes user from a team by his/here request
func httpPostLeaveTeam(w http.ResponseWriter, r *http.Request) {
	verifyOptions := verify.Request(verify.MinimumContentLength(apicore.MinJSONRequestSize), verify.MaximumContentLength(10*apicore.KB), verify.AuthenticationRequired(true))
	ctx, userContext, err := apicore.VerifyRequestAndCreateUserContext(w, r, verifyOptions)
	if err != nil {
		return
	}
	var request dto4teamus.LeaveTeamRequest
	if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	err = leaveTeam(ctx, userContext, request)
	apicore.IfNoErrorReturnOK(ctx, w, r, err)
}
