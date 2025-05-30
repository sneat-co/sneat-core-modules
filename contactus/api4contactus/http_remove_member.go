package api4contactus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/facade4contactus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

// httpPostRemoveSpaceMember is an API endpoint that removes a members from a space
func httpPostRemoveSpaceMember(w http.ResponseWriter, r *http.Request) {
	ctx, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}
	var request dto4contactus.ContactRequest
	if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	err = facade4contactus.RemoveSpaceMember(ctx, request)
	apicore.IfNoErrorReturnOK(ctx, w, r, err)
}
