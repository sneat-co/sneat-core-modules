package api4auth

import (
	"github.com/strongo/strongoapp"
	"net/http"
)

func InitApiForAuth(handle strongoapp.HandleHttpWithContext) {
	handle(http.MethodPost, "/api4debtus/auth/login-id", OptionalAuth(handleAuthLoginId))
	handle(http.MethodPost, "/api4debtus/auth/sign-in-with-fbm", OptionalAuth(handleSignInWithFbm))
	handle(http.MethodPost, "/api4debtus/auth/sign-in", OptionalAuth(handleSignInWithCode))
	handle(http.MethodPost, "/api4debtus/auth/fb/signed", OptionalAuth(handleSignedWithFacebook))
	handle(http.MethodPost, "/api4debtus/auth/google-plus/signed", OptionalAuth(handleSignedInWithGooglePlus))
	handle(http.MethodPost, "/api4debtus/auth/vk/signed", OptionalAuth(handleSignedWithVK))
	//handle(http.MethodPost, "/api4debtus/auth/email-sign-up", HandleSignUpWithEmail)
	//handle(http.MethodPost, "/api4debtus/auth/email-sign-in", HandleSignInWithEmail)
	handle(http.MethodPost, "/api4debtus/auth/request-password-reset", handleRequestPasswordReset)
	handle(http.MethodPost, "/api4debtus/auth/change-password-and-sign-in", handleChangePasswordAndSignIn2)
	handle(http.MethodPost, "/api4debtus/auth/confirm-email-and-sign-in", handleConfirmEmailAndSignIn2)
	handle(http.MethodPost, "/api4debtus/auth/anonymous-sign-up", handleSignUpAnonymously)
	handle(http.MethodPost, "/api4debtus/auth/anonymous-sign-in", handleSignInAnonymous)
	handle(http.MethodPost, "/api4debtus/auth/disconnect", AuthOnly(handleDisconnect))
}
