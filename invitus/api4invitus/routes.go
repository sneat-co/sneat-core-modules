package api4invitus

import "net/http"

// RegisterHttpRoutes registers invites routes
func RegisterHttpRoutes(handle func(method, path string, handler http.HandlerFunc)) {
	handle(http.MethodPost, "/v0/invites/create_invite_for_member", httpPostCreateOrReuseInviteForMember)
	handle(http.MethodGet, "/v0/invites/invite_link_for_member", httpGetOrCreateInviteLink)
	handle(http.MethodGet, "/v0/invites/personal_invite", httpGetPersonal)
	handle(http.MethodPost, "/v0/invites/create_mass_invite", httpPostCreateMassInvite)
	handle(http.MethodPost, "/v0/invites/accept_personal_invite", httpPostClaimPersonalInvite)
	handle(http.MethodPost, "/v0/invites/decline_personal_invite", httpPostClaimPersonalInvite)
}
