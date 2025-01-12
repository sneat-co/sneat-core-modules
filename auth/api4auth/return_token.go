package api4auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/sneat-co/sneat-core-modules/auth/token4auth"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/strongo/validation"
	"net/http"
)

type TokenClaim func(claims *TokenClaims)

type TokenClaims struct {
	isAdmin bool
}

func (t *TokenClaims) IsAdmin() bool {
	return t.isAdmin
}

func IsAdminClaim() func(claim *TokenClaims) {
	return func(claim *TokenClaims) {
		claim.isAdmin = true
	}
}

// ReturnToken returns token
func ReturnToken(ctx context.Context, w http.ResponseWriter, r *http.Request, userID, issuer string, options ...TokenClaim) {
	claims := TokenClaims{}
	for _, o := range options {
		o(&claims)
	}
	if claims.isAdmin {
		apicore.ReturnError(ctx, w, r, validation.NewBadRequestError(errors.New("issuing admin token is not implemented yet")))
		return
	}
	token, err := token4auth.IssueAuthToken(ctx, userID, issuer)
	if err != nil {
		err = fmt.Errorf("failed to issue Firebase token: %w", err)
		apicore.ReturnError(ctx, w, r, err)
		return
	}
	header := w.Header()

	// If decided to remove or add the Access-Control-Allow-Origin header - comment  the reason for doing that.
	// Reason for removing: Apparently, we add the header in an HTTP handler wrapper using the request's ORIGIN header.
	// Header.Add("Access-Control-Allow-Origin", "*")

	header.Add("Content-Type", "application/json")
	responseBody := fmt.Sprintf(`{"token":"%s"}`, token)
	_, _ = w.Write([]byte(responseBody))
}
