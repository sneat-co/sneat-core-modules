package api4invitus

import (
	"github.com/sneat-co/sneat-core-modules/invitus/facade4invitus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

// httpPostClaimPersonalInvite is an API endpoint that marks a personal invite as accepted
func httpPostClaimPersonalInvite(w http.ResponseWriter, r *http.Request) {
	ctx, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}

	var request facade4invitus.ClaimPersonalInviteRequest
	if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}

	request.RemoteClient = apicore.GetRemoteClientInfo(r)

	_, err = facade4invitus.ClaimPersonalInvite(ctx, request)

	apicore.IfNoErrorReturnOK(ctx, w, r, err)
}
