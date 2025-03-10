package api4invitus

import (
	"fmt"
	"github.com/sneat-co/sneat-core-modules/invitus/facade4invitus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/httpserver"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"net/http"
)

// httpPostCreateOrReuseInviteForMember supports both POST & GET methods
func httpPostCreateOrReuseInviteForMember(w http.ResponseWriter, r *http.Request) {
	var request facade4invitus.InviteContactRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusCreated,
		func(ctx facade.ContextWithUser) (any, error) {
			if request.To.Channel == "link" {
				return nil, fmt.Errorf("%w: link invites should be requested via GET", facade.ErrBadRequest)
			}
			var resp facade4invitus.CreateInviteResponse
			var err error
			resp, err = facade4invitus.CreateOrReuseInviteToContact(ctx, request, func() dbmodels.RemoteClientInfo {
				return apicore.GetRemoteClientInfo(r)
			})
			return resp.Invite.ID, err
		})
}

// httpGetOrCreateInviteLink gets or creates an invitation link
func httpGetOrCreateInviteLink(w http.ResponseWriter, r *http.Request) {
	var request facade4invitus.InviteContactRequest
	q := r.URL.Query()

	if request.SpaceID = coretypes.SpaceID(q.Get("space")); request.SpaceID == "" {
		apicore.ReturnError(r.Context(), w, r, validation.NewErrRequestIsMissingRequiredField("space"))
		// TODO(deprecate): httpserver.HandleError(nil, validation.NewErrRequestIsMissingRequiredField("space"), "httpGetOrCreateInviteLink", w, r)
		return
	}
	if request.To.ContactID = q.Get("contact"); request.To.ContactID == "" {
		apicore.ReturnError(r.Context(), w, r, validation.NewErrRequestIsMissingRequiredField("contact"))
		return
	}

	request.To.Channel = "link"
	ctx, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.Request(
		verify.AuthenticationRequired(true),
		verify.MaximumContentLength(0),
	))
	if err != nil {
		httpserver.HandleError(ctx, err, "VerifyRequestAndCreateUserContext", w, r)
		return
	}
	var resp facade4invitus.CreateInviteResponse
	resp, err = facade4invitus.CreateOrReuseInviteToContact(ctx, request, func() dbmodels.RemoteClientInfo {
		return apicore.GetRemoteClientInfo(r)
	})
	if err != nil {
		apicore.ReturnError(ctx, w, r, err)
		return
	}
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, resp.Invite.ID)
}
