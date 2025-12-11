package facade4contactus

import (
	"fmt"
	"slices"

	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/strongo/slice"
)

func updateContactRoles(params *dal4contactus.ContactWorkerParams, roles dto4contactus.SetRolesRequest) (updatedContactFields []string, err error) {
	var removedCount int
	var addedCount int
	params.Contact.Data.Roles, removedCount = slice.RemoveInPlace(params.Contact.Data.Roles, func(v string) bool {
		return slices.Contains(roles.Remove, v)
	})
	for _, role := range roles.Add {
		if !slices.Contains(params.Contact.Data.Roles, role) {
			addedCount++
			params.Contact.Data.Roles = append(params.Contact.Data.Roles, role)
		}
	}
	if removedCount > 0 || addedCount > 0 {
		updatedContactFields = append(updatedContactFields, "roles")
		params.ContactUpdates = append(params.ContactUpdates, update.ByFieldName("roles", params.Contact.Data.Roles))
		params.SpaceModuleUpdates = append(params.SpaceModuleUpdates,
			update.ByFieldName(fmt.Sprintf("contacts.%s.roles", params.Contact.ID), params.Contact.Data.Roles))
	}

	return updatedContactFields, err
}

func removeContactRoles(params *dal4contactus.ContactWorkerParams) {
	contact := params.Contact
	contactBrief := params.SpaceModuleEntry.Data.GetContactBriefByContactID(contact.ID)
	if contactBrief != nil {
		for _, u := range contactBrief.RemoveRole(const4contactus.SpaceMemberRoleMember) {
			params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, update.ByFieldName(
				fmt.Sprintf("contacts.%s.roles.%s", contact.ID, u.FieldName()),
				contact.Data.Roles))
		}
	}
}
