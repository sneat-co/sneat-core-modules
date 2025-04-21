package api4contactus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/facade4contactus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

// httpDeleteContact deletes contact
func httpDeleteContact(w http.ResponseWriter, r *http.Request) {
	var request dto4contactus.ContactRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusCreated,
		func(ctx facade.ContextWithUser) (any, error) {
			return nil, facade4contactus.DeleteContact(ctx, request)
		})
}
