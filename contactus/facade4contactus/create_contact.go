package facade4contactus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/models4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage"
	"github.com/sneat-co/sneat-core-modules/teamus/core4teamus"
	"github.com/sneat-co/sneat-core-modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/slice"
)

// CreateContact creates team contact
func CreateContact(
	ctx context.Context,
	userContext facade.User,
	userCanBeNonTeamMember bool,
	request dto4contactus.CreateContactRequest,
) (
	response dto4contactus.CreateContactResponse,
	err error,
) {
	if err = request.Validate(); err != nil {
		return response, fmt.Errorf("invalid CreateContactRequest: %w", err)
	}

	err = dal4teamus.CreateTeamItem(ctx, userContext, const4contactus.ContactsCollection, request.TeamRequest, const4contactus.ModuleID, new(models4contactus.ContactusTeamDto),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.ModuleTeamWorkerParams[*models4contactus.ContactusTeamDto]) (err error) {
			var contact dal4contactus.ContactEntry
			if contact, err = CreateContactTx(ctx, tx, userCanBeNonTeamMember, request, params); err != nil {
				return err
			}
			response = dto4contactus.CreateContactResponse{
				ID:   contact.ID,
				Data: contact.Data,
			}
			return err
		},
	)
	if err != nil {
		err = fmt.Errorf("failed to create a new contact: %w", err)
		return
	}
	return
}

func CreateContactTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userCanBeNonTeamMember bool,
	request dto4contactus.CreateContactRequest,
	params *dal4teamus.ModuleTeamWorkerParams[*models4contactus.ContactusTeamDto],
) (
	contact dal4contactus.ContactEntry,
	err error,
) {
	if err = request.Validate(); err != nil {
		return
	}
	if err = params.GetRecords(ctx, tx, params.UserID); err != nil {
		return
	}
	userContactID, userContactBrief := params.TeamModuleEntry.Data.GetContactBriefByUserID(params.UserID)
	if !userCanBeNonTeamMember && (userContactBrief == nil || !userContactBrief.IsTeamMember()) {
		err = errors.New("user is not a member of the team")
		return
	}
	switch userContactBrief.AgeGroup {
	case "", dbmodels.AgeGroupUnknown:
		if request.RelatedTo != nil {
			for _, relatedAs := range request.RelatedTo.RelatedAs {
				switch relatedAs {
				case dbmodels.RelationshipSpouse, dbmodels.RelationshipChild:
					userContactBrief.AgeGroup = dbmodels.AgeGroupAdult
					userContactKey := dal4contactus.NewContactKey(request.TeamID, userContactID)
					if err = tx.Update(ctx, userContactKey, []dal.Update{
						{
							Field: "ageGroup",
							Value: userContactBrief.AgeGroup,
						},
					}); err != nil {
						err = fmt.Errorf("failed to update member record: %w", err)
						return
					}
				}
			}
		}
	}

	parentContactID := request.ParentContactID

	var parent dal4contactus.ContactEntry
	if parentContactID != "" {
		parent = dal4contactus.NewContactEntry(request.TeamID, parentContactID)
		if err = tx.Get(ctx, parent.Record); err != nil {
			return contact, fmt.Errorf("failed to get parent contact with ID=[%s]: %w", parentContactID, err)
		}
	}

	var contactDto models4contactus.ContactDto
	contactDto.CreatedAt = params.Started
	contactDto.CreatedBy = params.UserID
	contactDto.Status = "active"
	contactDto.ParentID = parentContactID
	contactDto.WithRoles = request.WithRoles
	if request.Person != nil {
		contactDto.ContactBase = request.Person.ContactBase
		contactDto.Type = briefs4contactus.ContactTypePerson
		if contactDto.AgeGroup == "" {
			contactDto.AgeGroup = "unknown"
		}
		if contactDto.Gender == "" {
			contactDto.Gender = "unknown"
		}
		contactDto.ContactBase = request.Person.ContactBase
		for _, role := range request.Roles {
			if !slice.Contains(contactDto.Roles, role) {
				contactDto.Roles = append(contactDto.Roles, role)
			}
		}
	} else if request.Company != nil {
		contactDto.Type = briefs4contactus.ContactTypeCompany
		contactDto.Title = request.Company.Title
		contactDto.VATNumber = request.Company.VATNumber
		contactDto.Address = request.Company.Address
	} else if request.Location != nil {
		contactDto.Type = briefs4contactus.ContactTypeLocation
		contactDto.Title = request.Location.Title
		contactDto.Address = &request.Location.Address
	} else if request.Basic != nil {
		contactDto.Type = request.Type
		contactDto.Title = request.Basic.Title
	} else {
		return contact, errors.New("contact type is not specified")
	}
	if contactDto.Address != nil {
		contactDto.CountryID = contactDto.Address.CountryID
	}
	contactDto.ShortTitle = contactDto.DetermineShortTitle(request.Person.Title, params.TeamModuleEntry.Data.Contacts)
	var contactID string
	if request.ContactID == "" {
		contactID, err = dbmodels.NewUniqueRandomID(params.TeamModuleEntry.Data.ContactIDs(), 3)
		if err != nil {
			return contact, fmt.Errorf("failed to generate new contact ItemID: %w", err)
		}
	} else {
		contactID = request.ContactID
	}
	if contactDto.CountryID == "" && params.Team.Data.CountryID != "" && params.Team.Data.Type == core4teamus.TeamTypeFamily {
		contactDto.CountryID = params.Team.Data.CountryID
	}
	params.TeamModuleEntry.Data.AddContact(contactID, &contactDto.ContactBrief)
	if params.TeamModuleEntry.Record.Exists() {
		if err = tx.Update(ctx, params.TeamModuleEntry.Key, []dal.Update{
			{
				Field: const4contactus.ContactsField,
				Value: params.TeamModuleEntry.Data.Contacts,
			},
		}); err != nil {
			return contact, fmt.Errorf("failed to update team contact briefs: %w", err)
		}
	} else {
		if err = tx.Insert(ctx, params.TeamModuleEntry.Record); err != nil {
			return contact, fmt.Errorf("faield to insert team contacts brief record: %w", err)
		}
	}

	params.TeamUpdates = append(params.TeamUpdates, params.Team.Data.UpdateNumberOf(const4contactus.ContactsField, len(params.TeamModuleEntry.Data.Contacts)))

	if request.RelatedTo != nil {
		if request.RelatedTo.ItemID == "" {
			request.RelatedTo.ItemID, _ = params.TeamModuleEntry.Data.GetContactBriefByUserID(params.UserID)
			if request.RelatedTo.ItemID == "" {
				err = errors.New("user does not have a contact brief in contactus team record")
				return
			}
		}
		if _, err = contactDto.SetRelationshipToItem(
			params.UserID,
			linkage.TeamModuleDocRef{
				TeamID:     request.TeamID,
				ModuleID:   const4contactus.ModuleID,
				Collection: const4contactus.ContactsCollection,
				ItemID:     contactID,
			}, *request.RelatedTo, params.Started,
		); err != nil {
			return contact, err
		}
	}

	contact = dal4contactus.NewContactEntryWithData(request.TeamID, contactID, &contactDto)

	//contact.Data.UserIDs = params.Team.Data.UserIDs
	if err := contact.Data.Validate(); err != nil {
		return contact, fmt.Errorf("contact record is not valid: %w", err)
	}
	if err = tx.Insert(ctx, contact.Record); err != nil {
		return contact, fmt.Errorf("failed to insert contact record: %w", err)
	}

	if parent.ID != "" {
		if err = updateParentContact(ctx, tx, contact, parent); err != nil {
			return contact, fmt.Errorf("failed to update parent contact: %w", err)
		}
	}
	return contact, err
}
