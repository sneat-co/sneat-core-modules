package facade4contactus

import (
	"context"
	"fmt"
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/checks4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// CreateMember adds members to a space
func CreateMember(
	ctx context.Context,
	userCtx facade.UserContext,
	request dal4contactus.CreateMemberRequest,
) (
	contact dal4contactus.ContactEntry,
	err error,
) {
	if err = request.Validate(); err != nil {
		return contact, fmt.Errorf("invalid CreateMemberRequest: %w", err)
	}
	createContactRequest := dto4contactus.CreateContactRequest{
		SpaceRequest: request.SpaceRequest,
		WithRelated:  request.WithRelated,
		Status:       request.Status,
		Type:         briefs4contactus.ContactTypePerson,
		Person:       &request.CreatePersonRequest,
	}
	if !checks4contactus.IsSpaceMember(request.Roles) {
		createContactRequest.Roles = append(createContactRequest.Roles, const4contactus.SpaceMemberRoleMember)
	}
	if contact, err = CreateContact(ctx, userCtx, false, createContactRequest); err != nil {
		return contact, err
	}
	if contact.Data == nil {
		return contact, fmt.Errorf("CreateContact returned nil response data")
	}
	if !checks4contactus.IsSpaceMember(contact.Data.Roles) {
		err = fmt.Errorf("created contact does not have space member role")
		return contact, err
	}
	return contact, err
}
