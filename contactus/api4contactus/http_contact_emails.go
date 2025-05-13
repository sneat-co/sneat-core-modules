package api4contactus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/facade4contactus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

func httpAddEmail(w http.ResponseWriter, r *http.Request) {
	var request dto4contactus.AddEmailRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusNoContent,
		func(ctx facade.ContextWithUser) (body any, err error) {
			err = facade4contactus.AddEmail(ctx, request)
			return nil, err
		})
}

func httpDeleteEmail(w http.ResponseWriter, r *http.Request) {
	var request dto4contactus.DeleteEmailRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusNoContent,
		func(ctx facade.ContextWithUser) (body any, err error) {
			err = facade4contactus.DeleteEmail(ctx, request)
			return nil, err
		})
}
