package api4userus

import (
	"github.com/sneat-co/sneat-core-modules/userus/facade4userus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

func httpSetUserCountry(w http.ResponseWriter, r *http.Request) {
	var request facade4userus.SetUserCountryRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusNoContent,
		func(ctx facade.ContextWithUser) (response interface{}, err error) {
			return nil, facade4userus.SetUserCountry(ctx, request)
		})
}
