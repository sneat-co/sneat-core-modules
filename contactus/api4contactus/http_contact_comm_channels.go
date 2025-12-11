package api4contactus

import (
	"net/http"

	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/facade4contactus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
)

func httpAddContactCommChannel(w http.ResponseWriter, r *http.Request) {
	var request dto4contactus.AddCommChannelRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusNoContent,
		func(ctx facade.ContextWithUser) (body any, err error) {
			err = facade4contactus.AddCommChannel(ctx, request)
			return nil, err
		})
}

func httpUpdateContactCommChannel(w http.ResponseWriter, r *http.Request) {
	var request dto4contactus.UpdateCommChannelRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusNoContent,
		func(ctx facade.ContextWithUser) (body any, err error) {
			err = facade4contactus.UpdateCommChannel(ctx, request)
			return nil, err
		})
}
func httpDeleteContactCommChannel(w http.ResponseWriter, r *http.Request) {
	var request dto4contactus.DeleteCommChannelRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusNoContent,
		func(ctx facade.ContextWithUser) (body any, err error) {
			err = facade4contactus.DeleteCommChannel(ctx, request)
			return nil, err
		})
}
