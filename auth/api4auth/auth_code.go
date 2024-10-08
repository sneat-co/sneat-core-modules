package api4auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/sneat-co/sneat-core-modules/anybot/facade4anybot"
	"github.com/sneat-co/sneat-core-modules/auth/models4auth"
	"github.com/sneat-co/sneat-core-modules/auth/token4auth"
	"github.com/sneat-co/sneat-core-modules/auth/unsorted4auth"
	common4all2 "github.com/sneat-co/sneat-core-modules/common4all"
	"github.com/strongo/logus"
	"net/http"
	"strconv"
)

// TODO: Obsolete - migrate to HandleSignInWithPin
func HandleSignInWithCode(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
	code := r.PostFormValue("code")
	if code == "" {
		common4all2.BadRequestMessage(ctx, w, "Missing required attribute: code")
		return
	}
	if loginCode, err := strconv.Atoi(code); err != nil {
		common4all2.BadRequestMessage(ctx, w, "Parameter code is not an integer")
		return
	} else if loginCode == 0 {
		common4all2.ErrorAsJson(ctx, w, http.StatusBadRequest, errors.New("Login code should not be 0."))
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
				common4all2.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
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

func HandleSignInWithPin(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
	loginID, err := common4all2.DecodeIntID(r.PostFormValue("loginID"))
	if err != nil {
		common4all2.BadRequestError(ctx, w, fmt.Errorf("parameter 'loginID' is not an integer: %w", err))
		return
	}

	if loginCode, err := strconv.ParseInt(r.PostFormValue("loginPin"), 10, 32); err != nil {
		common4all2.BadRequestMessage(ctx, w, "Parameter 'loginCode' is not an integer")
		return
	} else if loginCode == 0 {
		common4all2.ErrorAsJson(ctx, w, http.StatusBadRequest, errors.New("Parameter 'loginCode' should not be 0."))
		return
	} else {
		if userID, err := facade4anybot.AuthFacade.SignInWithPin(ctx, loginID, int32(loginCode)); err != nil {
			switch err {
			case facade4anybot.ErrLoginExpired:
				_, _ = w.Write([]byte("expired"))
			case facade4anybot.ErrLoginAlreadySigned:
				_, _ = w.Write([]byte("claimed"))
			default:
				err = fmt.Errorf("failed to claim loginCode: %w", err)
				common4all2.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
			}
		} else {
			if authInfo.UserID != "" && userID != authInfo.UserID {
				logus.Warningf(ctx, "userID:%s != authInfo.AppUserIntID:%s", userID, authInfo.UserID)
			}
			ReturnToken(ctx, w, r, userID, r.Referer())
		}
	}
}
