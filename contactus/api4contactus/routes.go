package api4contactus

import (
	"github.com/sneat-co/sneat-go-core/extension"
	"net/http"
)

// RegisterHttpRoutes registers contact routes
func RegisterHttpRoutes(handle extension.HTTPHandleFunc) {
	handle(http.MethodPost, "/v0/contactus/create_contact", httpPostCreateContact)
	handle(http.MethodDelete, "/v0/contactus/delete_contact", httpDeleteContact)
	handle(http.MethodPost, "/v0/contactus/set_contacts_status", httpSetContactStatus)
	handle(http.MethodPost, "/v0/contactus/update_contact", httpUpdateContact)
	handle(http.MethodPost, "/v0/contactus/archive_contact", httpPostArchiveContact)
	handle(http.MethodPost, "/v0/contactus/create_member", httpPostCreateMember)
	handle(http.MethodPost, "/v0/contactus/remove_space_member", httpPostRemoveSpaceMember)

	handle(http.MethodPost, "/v0/contactus/add_contact_comm_channel", httpAddContactCommChannel)
	handle(http.MethodPost, "/v0/contactus/update_contact_comm_channel", httpUpdateContactCommChannel)
	handle(http.MethodPost, "/v0/contactus/delete_contact_comm_channel", httpDeleteContactCommChannel)

	////
	//handle(http.MethodGet, "/v0/space/join_info", api4debtus.GetSpaceJoinInfo)
	//handle(http.MethodPost, "/v0/space/join_space", api4debtus.JoinSpace)
	//handle(http.MethodPost, "/v0/space/refuse_to_join_space", api4debtus.RefuseToJoinSpace)
	//handle(http.MethodPost, "/v0/space/leave_space", api4debtus.LeaveSpace)
	//handle(http.MethodPost, "/v0/space/create_member", api4debtus.AddMember)
	//handle(http.MethodPost, "/v0/space/add_metric", api4debtus.AddMetric)
	//handle(http.MethodPost, "/v0/space/remove_member", api4debtus.RemoveMember)
	//handle(http.MethodPost, "/v0/space/change_member_role", api4debtus.ChangeMemberRole)
	//handle(http.MethodPost, "/v0/space/remove_metrics", api4debtus.RemoveMetrics)
	//handle(http.MethodGet, "/v0/space", api4debtus.GetSpace)
}
