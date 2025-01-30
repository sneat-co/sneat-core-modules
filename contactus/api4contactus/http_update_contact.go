package api4contactus

import (
	"context"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/facade4contactus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

func httpUpdateContact(w http.ResponseWriter, r *http.Request) {
	var request dto4contactus.UpdateContactRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusNoContent,
		func(ctx context.Context, userCtx facade.UserContext) (body interface{}, err error) {
			_, err = facade4contactus.UpdateContact(ctx, userCtx, request)
			return nil, err
		})
}
