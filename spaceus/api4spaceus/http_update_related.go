package api4spaceus

import (
	"github.com/sneat-co/sneat-core-modules/linkage/dto4linkage"
	"github.com/sneat-co/sneat-core-modules/linkage/facade4linkage"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

// httpPostUpdateRelated is an API endpoint that adds or removes related items to a space item
func httpPostUpdateRelated(w http.ResponseWriter, r *http.Request) {
	var request dto4linkage.UpdateRelatedRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired,
		http.StatusNoContent,
		func(ctx facade.ContextWithUser) (body any, err error) {
			return nil, facade4linkage.UpdateRelatedAndIDsOfSpaceItem(ctx, request)
		},
	)
}
