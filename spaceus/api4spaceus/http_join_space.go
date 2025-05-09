package api4spaceus

import (
	"github.com/sneat-co/sneat-core-modules/invitus/facade4invitus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

// httpPostJoinSpace joins a members to a space
func httpPostJoinSpace(w http.ResponseWriter, r *http.Request) {
	var request facade4invitus.JoinSpaceRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusOK,
		func(ctx facade.ContextWithUser) (response any, err error) {
			return facade4invitus.JoinSpace(ctx, request)
		})
}
