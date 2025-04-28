package api4spaceus

import (
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/facade4spaceus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"net/http"
)

//var getSpaceByID = facade4spaceus.GetSpaceByID

// httpGetSpace is an API endpoint that return space data
func httpGetSpace(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id := coretypes.SpaceID(q.Get("id"))
	verifyOptions := verify.Request(verify.AuthenticationRequired(true))
	ctx, err := apicore.VerifyRequestAndCreateUserContext(w, r, verifyOptions)
	if err != nil {
		return
	}
	var space dbo4spaceus.SpaceEntry
	var response any
	if space, err = facade4spaceus.GetSpace(ctx, id); err == nil {
		response = space.Data
	}
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, response)
}
