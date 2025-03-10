package api4auth

import (
	"context"
	"github.com/sneat-co/sneat-core-modules/auth/token4auth"
	"github.com/sneat-co/sneat-core-modules/common4all"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/strongo/logus"
	"github.com/strongo/strongoapp"
	"net/http"
)

type AuthHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo)

type AuthHandlerWithUser func(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo, user dbo4userus.UserEntry)

func AuthOnly(handler AuthHandler) strongoapp.HttpHandlerWithContext {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if authInfo, _, err := token4auth.Authenticate(w, r, true); err == nil {
			handler(ctx, w, r, authInfo)
		} else {
			logus.Warningf(ctx, "Failed to authenticate: %v", err)
		}
	}
}

func AuthOnlyWithUser(handler AuthHandlerWithUser) strongoapp.HttpHandlerWithContext {
	return AuthOnly(func(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
		var userID string

		if userID = common4all.GetUserID(ctx, w, r, authInfo); userID == "" {
			logus.Warningf(ctx, "userID is empty")
			return
		}

		//user, err := dal4userus.GetUserByID(ctx, nil, userID)
		var user dbo4userus.UserEntry
		var err error

		if common4all.HasError(ctx, w, err, dbo4userus.UsersCollection, userID, http.StatusInternalServerError) {
			return
		}
		handler(ctx, w, r, authInfo, user)
	})
}

func OptionalAuth(handler AuthHandler) strongoapp.HttpHandlerWithContext {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		authInfo, _, _ := token4auth.Authenticate(w, r, false)
		if authInfo.UserID == "" {
			logus.Debugf(ctx, "OptionalAuth(), anonymous")
		} else {
			logus.Debugf(ctx, "OptionalAuth(), userID=%s", authInfo.UserID)
		}
		handler(ctx, w, r, authInfo)
	}
}

func AdminOnly(handler AuthHandler) strongoapp.HttpHandlerWithContext {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if authInfo, _, err := token4auth.Authenticate(w, r, true); err == nil {
			if !authInfo.IsAdmin {
				logus.Debugf(ctx, "Not admin!")
				//hashedWriter.WriteHeader(http.StatusForbidden)
				//return
			}
			handler(ctx, w, r, authInfo)
		} else {
			logus.Errorf(ctx, "Failed to authenticate: %v", err)
		}
	}
}
