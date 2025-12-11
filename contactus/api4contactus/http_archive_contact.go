package api4contactus

import (
	"net/http"

	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/facade4contactus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
)

// httpPostArchiveContact is an API endpoint that archive contact - e.g., hides it from the list of contacts
func httpPostArchiveContact(w http.ResponseWriter, r *http.Request) {
	var request dto4contactus.ContactRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusOK,
		func(ctx facade.ContextWithUser) (any, error) {
			return nil, facade4contactus.ArchiveContact(ctx, request)
		})
}
