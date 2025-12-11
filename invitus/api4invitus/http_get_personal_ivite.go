package api4invitus

import (
	"github.com/sneat-co/sneat-core-modules/invitus/facade4invitus"
	"github.com/sneat-co/sneat-go-core/coretypes"

	"net/http"
	"strings"

	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
)

// httpGetPersonal is an API endpoint that returns personal invite data
func httpGetPersonal(w http.ResponseWriter, r *http.Request) {
	ctx, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}
	q := r.URL.Query()
	request := facade4invitus.GetPersonalInviteRequest{
		SpaceRequest: dto4spaceus.SpaceRequest{
			SpaceID: coretypes.SpaceID(strings.TrimSpace(q.Get("spaceID"))),
		},
		InviteID: strings.TrimSpace(q.Get("inviteID")),
	}
	response, err := facade4invitus.GetPersonal(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, response)
}
