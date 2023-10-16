package api4schedulus

import (
	"github.com/sneat-co/sneat-core-modules/schedulus/dto4schedulus"
	"net/http"
)

func getHappeningRequestParamsFromURL(r *http.Request) (request dto4schedulus.HappeningRequest) {
	query := r.URL.Query()
	request.TeamID = query.Get("teamID")
	request.HappeningID = query.Get("happeningID")
	request.HappeningType = query.Get("happeningType")
	return
}
