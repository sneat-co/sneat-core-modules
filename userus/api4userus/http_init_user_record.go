package api4userus

import (
	"github.com/sneat-co/sneat-core-modules/auth/facade4auth"
	"github.com/sneat-co/sneat-core-modules/userus/dto4userus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/dto4auth"
	"github.com/strongo/strongoapp/appuser"
	"net/http"
	"strings"
)

// httpInitUserRecord sets user title
func httpInitUserRecord(w http.ResponseWriter, r *http.Request) {
	ctx, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}
	var request dto4userus.InitUserRecordRequest
	if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	request.RemoteClient = apicore.GetRemoteClientInfo(r)
	var params facade4auth.CreateUserWorkerParams
	userToCreate := dto4auth.DataToCreateUser{
		AuthAccount: appuser.AccountKey{
			Provider: request.AuthProvider,
			ID:       strings.ToLower(strings.TrimSpace(request.Email)),
		},
		Email:           request.Email,
		EmailIsVerified: request.EmailIsVerified,
		IanaTimezone:    request.IanaTimezone,
		RemoteClient:    request.RemoteClient,
	}
	if request.Names != nil {
		userToCreate.Names = *request.Names
	}
	if params, err = facade4auth.CreateUserRecords(ctx, userToCreate, true); err != nil {
		apicore.ReturnError(ctx, w, r, err)
		return
	}
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, params.User.Data)
}
