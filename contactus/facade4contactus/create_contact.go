package facade4contactus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/facade4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/logus"
	"github.com/strongo/strongoapp/person"
	"reflect"
	"slices"
)

var ErrContactWithSameAccountKeyAlreadyExists = errors.New("contact with the same account key already exists")

// CreateContact creates space contact
func CreateContact(
	ctx facade.ContextWithUser,
	userCanBeNonSpaceMember bool,
	request dto4contactus.CreateContactRequest,
) (
	contact dal4contactus.ContactEntry,
	err error,
) {
	// De-normalize & sanitize request if required
	if request.Type == briefs4contactus.ContactTypePerson && request.Person != nil {
		if request.Person.Joined {
			request.Person.Joined = false
		}
		if request.Person.Type == "" {
			request.Person.Type = request.Type
		}
		if request.Person.Gender == "" {
			request.Person.Gender = dbmodels.GenderUnknown
		}
		if request.Person.AgeGroup == "" {
			request.Person.AgeGroup = dbmodels.AgeGroupUnknown
		}
	}

	if err = request.Validate(); err != nil {
		return contact, fmt.Errorf("invalid CreateContactRequest: %w", err)
	}

	err = dal4spaceus.CreateSpaceItem(ctx, request.SpaceRequest, const4contactus.ModuleID, new(dbo4contactus.ContactusSpaceDbo),
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4spaceus.ModuleSpaceWorkerParams[*dbo4contactus.ContactusSpaceDbo]) (err error) {
			if contact, err = CreateContactTx(ctx, tx, userCanBeNonSpaceMember, request, params); err != nil {
				return fmt.Errorf("failed in CreateContactTx(): %w", err)
			}
			if contact.ID == "" {
				return errors.New("function CreateContactTx returned empty contact.ID")
			}
			logus.Debugf(ctx, "Created contact: %s", contact.Key.String())
			if contact.Data == nil {
				return errors.New("function CreateContactTx returned nil contact data")
			}
			return err
		},
	)
	if err != nil {
		err = fmt.Errorf("failed in dal4spaceus.CreateSpaceItem(): %w", err)
		return
	}
	return
}

func CreateContactTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userCanBeNonSpaceMember bool,
	request dto4contactus.CreateContactRequest,
	params *dal4spaceus.ModuleSpaceWorkerParams[*dbo4contactus.ContactusSpaceDbo],
) (
	contact dal4contactus.ContactEntry,
	err error,
) {
	if err = request.Validate(); err != nil {
		return
	}
	if err = params.GetRecords(ctx, tx); err != nil {
		return
	}
	userID := params.UserID()
	now := params.Started
	userContactID, userContactBrief := params.SpaceModuleEntry.Data.GetContactBriefByUserID(userID)
	if !userCanBeNonSpaceMember && (userContactBrief == nil || !userContactBrief.IsSpaceMember()) {
		err = errors.New("user is not a member of the space")
		return
	}
	if len(request.Accounts) > 0 {
		spaceContactusModuleKey := dbo4spaceus.NewSpaceModuleKey(params.Space.ID, const4contactus.ModuleID)
		recordMaker := func() dal.Record {
			return dal.NewRecordWithData(dal.NewIncompleteKey(const4contactus.ContactsCollection, reflect.String, spaceContactusModuleKey), new(dbo4contactus.ContactDbo))
		}
		query := dal.
			From(dal.NewCollectionRef(const4contactus.ContactsCollection, "c", spaceContactusModuleKey)).
			WhereInArrayField("accounts", request.Accounts[0]).
			Limit(1).
			SelectInto(recordMaker)
		var reader dal.Reader
		if reader, err = tx.QueryReader(ctx, query); err != nil {
			err = fmt.Errorf("failed to query contacts by account: %w", err)
			return
		}
		var contactRecords []dal.Record
		if contactRecords, err = dal.ReadAll(ctx, reader, 0); err != nil {
			err = fmt.Errorf("failed to load contacts records by account ID: %w", err)
			return
		}
		if len(contactRecords) > 0 { // TODO: Handle gracefully?
			err = fmt.Errorf("%w: %s", ErrContactWithSameAccountKeyAlreadyExists, request.Accounts[0])
			return
		}
	}
	if request.Related != nil {
		relatedByCollection := request.Related[string(const4contactus.ModuleID)]
		if relatedByCollection != nil {
			relatedItems := relatedByCollection[const4contactus.ContactsCollection]
			if len(relatedItems) > 0 {
				var isRelatedByUserID bool
				for _, relatedItem := range relatedItems {
					if _, isRelatedByUserID = relatedItems[userContactID]; !isRelatedByUserID {
						if contactBrief := params.SpaceModuleEntry.Data.GetContactBriefByContactID(userContactID); contactBrief == nil {
							return contact, fmt.Errorf("contact with ContactID=[%s] is not found", userContactID)
						}
					}
					switch userContactBrief.AgeGroup {
					case "", dbmodels.AgeGroupUnknown:
						for relatedAs := range relatedItem.RolesOfItem {
							switch relatedAs {
							case dbmodels.RelationshipSpouse, dbmodels.RelationshipChild:
								userContactBrief.AgeGroup = dbmodels.AgeGroupAdult
								userContactKey := dal4contactus.NewContactKey(request.SpaceID, userContactID)
								if err = tx.Update(ctx, userContactKey, []update.Update{
									update.ByFieldName("ageGroup", userContactBrief.AgeGroup),
								}); err != nil {
									err = fmt.Errorf("failed to update member record: %w", err)
									return
								}
							}
						}
					}
				}
			}
		}
	}

	parentContactID := request.ParentContactID

	var parent dal4contactus.ContactEntry
	if parentContactID != "" {
		parent = dal4contactus.NewContactEntry(request.SpaceID, parentContactID)
		if err = tx.Get(ctx, parent.Record); err != nil {
			return contact, fmt.Errorf("failed to get parent contact with ContactID=[%s]: %w", parentContactID, err)
		}
	}

	contactDbo := new(dbo4contactus.ContactDbo)
	contactDbo.CreatedAt = now
	contactDbo.CreatedBy = userID
	contactDbo.Status = "active"
	contactDbo.ParentID = parentContactID
	contactDbo.RolesField = request.RolesField
	contactDbo.EmailsField = request.EmailsField
	contactDbo.PhonesField = request.PhonesField
	if request.Person != nil {
		contactDbo.ContactBase = request.Person.ContactBase
		contactDbo.Type = briefs4contactus.ContactTypePerson
		if contactDbo.AgeGroup == "" {
			contactDbo.AgeGroup = "unknown"
		}
		if contactDbo.Gender == "" {
			contactDbo.Gender = "unknown"
		}
		contactDbo.ContactBase = request.Person.ContactBase
		for _, role := range request.Roles {
			if !slices.Contains(contactDbo.Roles, role) {
				contactDbo.Roles = append(contactDbo.Roles, role)
			}
		}
	} else if request.Company != nil {
		contactDbo.Type = briefs4contactus.ContactTypeCompany
		contactDbo.Title = request.Company.Title
		contactDbo.VATNumber = request.Company.VATNumber
		contactDbo.Address = request.Company.Address
	} else if request.Location != nil {
		contactDbo.Type = briefs4contactus.ContactTypeLocation
		contactDbo.Title = request.Location.Title
		contactDbo.Address = &request.Location.Address
	} else if request.Basic != nil {
		contactDbo.Type = request.Type
		contactDbo.Title = request.Basic.Title
	} else {
		return contact, errors.New("contact type is not specified")
	}
	if contactDbo.Address != nil {
		contactDbo.CountryID = contactDbo.Address.CountryID
	}
	if len(request.Accounts) > 0 {
		contactDbo.Accounts = request.Accounts
	}
	contactDbo.ShortTitle = contactDbo.DetermineShortTitle(request.Person.Title, params.SpaceModuleEntry.Data.Contacts)
	var contactID string
	if request.ContactID == "" {
		contactIDs := params.SpaceModuleEntry.Data.ContactIDs()
		if contactID, err = person.GenerateIDFromNameOrRandom(request.Person.Names, contactIDs); err != nil {
			return contact, fmt.Errorf("failed to generate contact ContactID: %w", err)
		}
	} else {
		contactID = request.ContactID
	}
	if contactDbo.CountryID == "" && params.Space.Data.CountryID != "" && params.Space.Data.Type == coretypes.SpaceTypeFamily {
		contactDbo.CountryID = params.Space.Data.CountryID
	}
	params.SpaceModuleEntry.Data.AddContact(contactID, &contactDbo.ContactBrief)
	if params.SpaceModuleEntry.Record.Exists() {
		if err = tx.Update(ctx, params.SpaceModuleEntry.Key, []update.Update{
			update.ByFieldName(const4contactus.ContactsField, params.SpaceModuleEntry.Data.Contacts),
		}); err != nil {
			return contact, fmt.Errorf("failed to update space contact briefs: %w", err)
		}
	} else {
		if err = tx.Insert(ctx, params.SpaceModuleEntry.Record); err != nil {
			return contact, fmt.Errorf("faield to insert space contacts brief record: %w", err)
		}
	}

	//params.SpaceUpdates = append(params.SpaceUpdates, params.SpaceID.Data.UpdateNumberOf(const4contactus.ContactsField, len(params.SpaceModuleEntry.Data.Contacts)))

	if len(request.Related) > 0 {
		if err = facade4linkage.UpdateRelationshipsInRelatedItems(
			now, userID, params.Space.ID, &contactDbo.WithRelatedAndIDs, request.Related,
		); err != nil {
			err = fmt.Errorf("failed to update relationships in related items: %w", err)
			return
		}
	}

	contact = dal4contactus.NewContactEntryWithData(request.SpaceID, contactID, contactDbo)

	if err = contact.Data.Validate(); err != nil {
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
