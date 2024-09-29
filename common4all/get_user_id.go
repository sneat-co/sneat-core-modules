package common4all

import (
	"context"
	"github.com/sneat-co/sneat-core-modules/auth/token4auth"
	"net/http"
)

func GetUserID(_ context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) (userID string) {
	userID = authInfo.UserID

	if stringID := r.URL.Query().Get("user"); stringID != "" {
		if !authInfo.IsAdmin && userID != authInfo.UserID {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
	return
}
