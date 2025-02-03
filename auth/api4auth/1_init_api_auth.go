package apimapping

import (
	"github.com/sneat-co/sneat-go/src/bots/botauth/api4botauth"
	"github.com/strongo/strongoapp"
	"net/http"
)

func InitApiForAuth(handle strongoapp.HandleHttpWithContext) {
	handle(http.MethodPost, "/api4debtus/auth/login-id", api4auth.OptionalAuth(api4auth.HandleAuthLoginId))
	handle(http.MethodPost, "/api4debtus/auth/sign-in-with-pin", api4auth.OptionalAuth(api4botauth.HandleSignInWithPin))
	handle(http.MethodPost, "/api4debtus/auth/sign-in-with-fbm", api4auth.OptionalAuth(api4auth.HandleSignInWithFbm))
	handle(http.MethodPost, "/api4debtus/auth/sign-in", api4auth.OptionalAuth(api4auth.HandleSignInWithCode))
	handle(http.MethodPost, "/api4debtus/auth/fb/signed", api4auth.OptionalAuth(api4auth.HandleSignedWithFacebook))
	handle(http.MethodPost, "/api4debtus/auth/google-plus/signed", api4auth.OptionalAuth(api4auth.HandleSignedInWithGooglePlus))
	handle(http.MethodPost, "/api4debtus/auth/vk/signed", api4auth.OptionalAuth(api4auth.HandleSignedWithVK))
	//handle(http.MethodPost, "/api4debtus/auth/email-sign-up", api4auth.HandleSignUpWithEmail)
	//handle(http.MethodPost, "/api4debtus/auth/email-sign-in", api4auth.HandleSignInWithEmail)
	handle(http.MethodPost, "/api4debtus/auth/request-password-reset", api4auth.HandleRequestPasswordReset)
	handle(http.MethodPost, "/api4debtus/auth/change-password-and-sign-in", api4auth.HandleChangePasswordAndSignIn)
	handle(http.MethodPost, "/api4debtus/auth/confirm-email-and-sign-in", api4auth.HandleConfirmEmailAndSignIn)
	handle(http.MethodPost, "/api4debtus/auth/anonymous-sign-up", api4auth.HandleSignUpAnonymously)
	handle(http.MethodPost, "/api4debtus/auth/anonymous-sign-in", api4auth.HandleSignInAnonymous)
	handle(http.MethodPost, "/api4debtus/auth/disconnect", api4auth.AuthOnly(api4auth.HandleDisconnect))
}
