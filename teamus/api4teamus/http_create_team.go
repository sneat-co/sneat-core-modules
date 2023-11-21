package api4teamus

import (
	"context"
	"github.com/sneat-co/sneat-core-modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-core-modules/teamus/facade4teamus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

// httpPostCreateTeam is an API endpoint that creates a new team
func httpPostCreateTeam(w http.ResponseWriter, r *http.Request) {
	var request dto4teamus.CreateTeamRequest
	handler := func(ctx context.Context, userCtx facade.User) (interface{}, error) {
		return facade4teamus.CreateTeam(ctx, userCtx, request)
	}
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, handler, http.StatusCreated, verify.DefaultJsonWithAuthRequired)
}
