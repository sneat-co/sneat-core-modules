package api4calendarium

import (
	"github.com/sneat-co/sneat-core-modules/calendarium/dto4calendarium"
	"net/http"
)

func getHappeningRequestParamsFromURL(r *http.Request) (request dto4calendarium.HappeningRequest) {
	query := r.URL.Query()
	request.TeamID = query.Get("teamID")
	request.HappeningID = query.Get("happeningID")
	request.HappeningType = query.Get("happeningType")
	return
}
