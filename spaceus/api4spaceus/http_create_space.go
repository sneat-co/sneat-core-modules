package api4spaceus

import (
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/facade4spaceus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

// httpPostCreateSpace is an API endpoint that creates a new space
func httpPostCreateSpace(w http.ResponseWriter, r *http.Request) {
	var request dto4spaceus.CreateSpaceRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusCreated,
		func(ctx facade.ContextWithUser) (any, error) {
			result, err := facade4spaceus.CreateSpace(ctx, request)
			if err != nil {
				return nil, err
			}
			var apiResponse dto4spaceus.CreateSpaceResponse
			apiResponse.Space.ID = result.Space.ID
			apiResponse.Space.Dbo = *result.Space.Data
			return apiResponse, err
		})
}
