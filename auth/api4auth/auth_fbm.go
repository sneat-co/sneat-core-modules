package api4auth

import (
	"github.com/sneat-co/sneat-core-modules/auth/token4auth"
	"github.com/sneat-co/sneat-core-modules/common4all"
	"github.com/strongo/logus"
	"net/http"

	"context"
)

func handleSignInWithFbm(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
	logus.Debugf(ctx, "api4debtus.handleSignInWithFbm()")

	threadID := r.PostFormValue("tid")
	threadType := r.PostFormValue("thread_type")
	pageScopedID := r.PostFormValue("psid")
	fbAppID := r.PostFormValue("fbAppID")

	if fbAppID == "" {
		common4all.BadRequestMessage(ctx, w, "Missing fbAppID")
		return
	}
	if threadID == "" {
		common4all.BadRequestMessage(ctx, w, "Missing tid")
		return
	}
	if threadType == "" {
		common4all.BadRequestMessage(ctx, w, "Missing thread_type")
		return
	}
	if pageScopedID == "" {
		common4all.BadRequestMessage(ctx, w, "Missing psid")
		return
	}

	// TODO: Log FbApp Code & FbPage Code (e.g. fbAppID=12345 => code=DebtsTracker)
	logus.Debugf(ctx, "FbmContext: thread_type=%v, tid=%v, psid=%v", threadType, threadID, pageScopedID)

	user, isNewUser, _, _, _, err := signInFbUser(ctx, fbAppID, pageScopedID, r, authInfo)
	if err != nil {
		authWriteResponseForAuthFailed(ctx, w, err)
		return
	}

	authWriteResponseForUser(ctx, w, r, user, "fbm", isNewUser)
}
