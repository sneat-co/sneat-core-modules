package facade4contactus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/models4contactus"
	"github.com/sneat-co/sneat-core-modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
)

// UpdateContactTx sets contact fields
func UpdateContactTx(ctx context.Context, tx dal.ReadwriteTransaction, user facade.User, request dto4contactus.UpdateContactRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4teamus.RunModuleTeamWorkerTx(ctx, tx, user, request.TeamRequest, const4contactus.ModuleID, new(models4contactus.ContactusTeamDto),
		func(ctx context.Context, tx dal.ReadwriteTransaction, teamWorkerParams *dal4teamus.ModuleTeamWorkerParams[*models4contactus.ContactusTeamDto]) (err error) {
			return updateContactTxWorker(ctx, tx, teamWorkerParams, request)
		},
	)
	if err != nil {
		return fmt.Errorf("failed to set contact status: %w", err)
	}
	return err
}

// UpdateContact sets contact fields
func UpdateContact(ctx context.Context, user facade.User, request dto4contactus.UpdateContactRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	db := facade.GetDatabase(ctx)
	return db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		return UpdateContactTx(ctx, tx, user, request)
	})
}

func updateContactTxWorker(
	ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.ModuleTeamWorkerParams[*models4contactus.ContactusTeamDto],
	request dto4contactus.UpdateContactRequest,
) (err error) {
	contact := dal4contactus.NewContactEntry(params.Team.ID, request.ContactID)
	if err = tx.Get(ctx, contact.Record); err != nil {
		return fmt.Errorf("failed to get contact record: %w", err)
	}

	if err := contact.Data.Validate(); err != nil {
		return fmt.Errorf("contact DTO is not valid after loading from DB: %w", err)
	}

	contactBrief := params.TeamModuleEntry.Data.Contacts[request.ContactID]

	var updatedContactFields []string
	var contactUpdates []dal.Update

	if request.Address != nil {
		if *request.Address != *contact.Data.Address {
			updatedContactFields = append(updatedContactFields, "address")
			contact.Data.Address = request.Address
			contactUpdates = append(contactUpdates, dal.Update{Field: "address", Value: request.Address})
		}
	}
	if request.AgeGroup != "" {
		if request.AgeGroup != contact.Data.AgeGroup {
			updatedContactFields = append(updatedContactFields, "ageGroup")
			contact.Data.AgeGroup = request.AgeGroup
			contactUpdates = append(contactUpdates, dal.Update{Field: "ageGroup", Value: contact.Data.AgeGroup})
		}
		if contactBrief != nil && contactBrief.AgeGroup != request.AgeGroup {
			params.TeamModuleUpdates = append(params.TeamModuleUpdates,
				dal.Update{
					Field: fmt.Sprintf("contacts.%s.ageGroup", request.ContactID),
					Value: contact.Data.AgeGroup,
				})
		}
	}
	if request.Roles != nil {
		for _, role := range request.Roles.Remove {
			contact.Data.Roles = slice.RemoveInPlace(role, contact.Data.Roles)
		}
		contact.Data.Roles = append(contact.Data.Roles, request.Roles.Add...)
		updatedContactFields = append(updatedContactFields, "roles")
		contactUpdates = append(contactUpdates, dal.Update{Field: "roles", Value: contact.Data.Roles})
		params.TeamModuleUpdates = append(params.TeamModuleUpdates,
			dal.Update{
				Field: fmt.Sprintf("contacts.%s.roles", request.ContactID),
				Value: contact.Data.Roles,
			})
	}
	if err := contact.Data.Validate(); err != nil {
		return fmt.Errorf("contact DTO is not valid after updating %d fields %+v: %w", len(updatedContactFields), updatedContactFields, err)
	}
	if len(contactUpdates) > 0 {
		if err := tx.Update(ctx, contact.Key, contactUpdates); err != nil {
			return err
		}
	}
	return nil
}
