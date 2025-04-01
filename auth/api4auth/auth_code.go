package api4auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/sneat-co/sneat-core-modules/auth/models4auth"
	"github.com/sneat-co/sneat-core-modules/auth/token4auth"
	"github.com/sneat-co/sneat-core-modules/auth/unsorted4auth"
	"github.com/sneat-co/sneat-core-modules/common4all"
	"github.com/strongo/logus"
	"net/http"
	"strconv"
)

// TODO: Obsolete - migrate to HandleSignInWithPin
func handleSignInWithCode(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
	code := r.PostFormValue("code")
	if code == "" {
		common4all.BadRequestMessage(ctx, w, "Missing required attribute: code")
		return
	}
	if loginCode, err := strconv.Atoi(code); err != nil {
		common4all.BadRequestMessage(ctx, w, "Parameter code is not an integer")
		return
	} else if loginCode == 0 {
		common4all.ErrorAsJson(ctx, w, http.StatusBadRequest, errors.New("login code should not be 0"))
		return
	} else {
		if userID, err := unsorted4auth.LoginCode.ClaimLoginCode(ctx, loginCode); err != nil {
			switch err {
			case models4auth.ErrLoginCodeExpired:
				_, _ = w.Write([]byte("expired"))
			case models4auth.ErrLoginCodeAlreadyClaimed:
				_, _ = w.Write([]byte("claimed"))
			default:
				err = fmt.Errorf("failed to claim code: %w", err)
				common4all.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
			}
		} else {
			if authInfo.UserID != "" && userID != authInfo.UserID {
				logus.Warningf(ctx, "userID:%s != authInfo.AppUserIntID:%s", userID, authInfo.UserID)
			}
			ReturnToken(ctx, w, r, userID, r.Referer())
			return
		}
	}
}
