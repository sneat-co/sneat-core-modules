package api4linkage

import (
	"context"
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
		func(ctx context.Context, userCtx facade.UserContext) (interface{}, error) {
			_, err := facade4linkage.UpdateItemRelationships(ctx, userCtx, request)
			return nil, err
		})
}
