package api4linkage

import (
	"github.com/sneat-co/sneat-core-modules/linkage/dto4linkage"
	"github.com/sneat-co/sneat-core-modules/linkage/facade4linkage"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

func httpUpdateItemRelationships(w http.ResponseWriter, r *http.Request) {
	var request dto4linkage.UpdateItemRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusNoContent,
		func(ctx facade.ContextWithUser) (any, error) {
			_, err := facade4linkage.UpdateItemRelationships(ctx, request)
			return nil, err
		})
}
