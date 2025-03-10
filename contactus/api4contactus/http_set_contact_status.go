package api4contactus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/facade4contactus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

var setContactsStatus = facade4contactus.SetContactsStatus

func httpSetContactStatus(w http.ResponseWriter, r *http.Request) {
	var request dto4contactus.SetContactsStatusRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusOK,
		func(ctx facade.ContextWithUser) (interface{}, error) {
			return nil, setContactsStatus(ctx, request)
		})
}
