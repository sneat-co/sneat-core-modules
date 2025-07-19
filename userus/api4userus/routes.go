package api4userus

import (
	"github.com/sneat-co/sneat-go-core/extension"
	"net/http"
)

// RegisterHttpRoutes initiates users module
func RegisterHttpRoutes(handle extension.HTTPHandleFunc) {
	handle(http.MethodPost, "/v0/users/init_user_record", httpInitUserRecord)
	handle(http.MethodPost, "/v0/users/link_auth_account", httpLinkAuthAccount)
	handle(http.MethodPost, "/v0/users/unlink_auth_account", httpUnlinkAuthAccount) // duplicate of /api4debtus/auth/disconnect
	handle(http.MethodPost, "/v0/users/set_user_country", httpSetUserCountry)
	//handle(http.MethodPost, "/v0/users/create_user", httpPostCreateUser)
}
