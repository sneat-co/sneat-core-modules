package facade4contactus

import (
	"context"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
)

// CreateMember adds members to a team
func CreateMember(
	ctx context.Context,
	user facade.User,
	request dal4contactus.CreateMemberRequest,
) (
	response dto4contactus.CreateContactResponse,
	err error,
) {
	createContactRequest := dto4contactus.CreateContactRequest{
		TeamRequest: request.TeamRequest,
		RelatedTo:   request.RelatedTo,
		Person:      &request.CreatePersonRequest,
	}
	if !slice.Contains(request.Roles, const4contactus.TeamMemberRoleMember) {
		createContactRequest.Roles = append(createContactRequest.Roles, const4contactus.TeamMemberRoleMember)
	}
	return CreateContact(ctx, user, false, createContactRequest)
}
