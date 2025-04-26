package facade4contactus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-core-modules/linkage/facade4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

// UpdateContact sets contact fields
func UpdateContact(
	ctx facade.ContextWithUser,
	request dto4contactus.UpdateContactRequest,
) (
	contact dal4contactus.ContactEntry,
	contactusSpace dal4contactus.ContactusSpaceEntry,
	space dbo4spaceus.SpaceEntry,
	err error,
) {
	user := ctx.User()
	updateContactWorker := func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactWorkerParams) (err error) {
		contact = params.Contact
		contactusSpace = params.SpaceModuleEntry
		space = params.Space
		return UpdateContactTx(ctx, tx, request, params)
	}
	err = dal4contactus.RunContactWorker(ctx, user, request.ContactRequest, updateContactWorker)
	return
}

// UpdateContactTx sets contact fields
func UpdateContactTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	request dto4contactus.UpdateContactRequest,
	params *dal4contactus.ContactWorkerParams,
) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	return updateContactTxWorker(ctx, tx, request, params)
}

func updateContactTxWorker(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	request dto4contactus.UpdateContactRequest,
	params *dal4contactus.ContactWorkerParams,
) (err error) {

	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}

	contact := params.Contact

	if err := contact.Data.Validate(); err != nil {
		return fmt.Errorf("contact DBO is not valid after loading from DB: %w", err)
	}

	contactBrief := params.SpaceModuleEntry.Data.Contacts[request.ContactID]

	updateContactBriefField := func(field string, value any) {
		params.SpaceModuleEntry.Record.MarkAsChanged()
		params.SpaceModuleUpdates = append(params.SpaceModuleUpdates,
			update.ByFieldName(
				fmt.Sprintf("contacts.%s.%s", request.ContactID, field),
				value))
	}

	var updatedContactFields []string

	if request.Address != nil {
		if *request.Address != *contact.Data.Address {
			updatedContactFields = append(updatedContactFields, "address")
			contact.Data.Address = request.Address
			params.ContactUpdates = append(params.ContactUpdates, update.ByFieldName("address", request.Address))
		}
	}

	if request.VatNumber != nil {
		if vat := *request.VatNumber; vat != contact.Data.VATNumber {
			updatedContactFields = append(updatedContactFields, "vatNumber")
			contact.Data.VATNumber = vat
			params.ContactUpdates = append(params.ContactUpdates, update.ByFieldName("vatNumber", vat))
		}
	}

	if request.Gender != "" {
		updatedContactFields = append(updatedContactFields, "gender")
		contact.Data.Gender = request.Gender
		params.ContactUpdates = append(params.ContactUpdates, update.ByFieldName("gender", request.Gender))
		if contactBrief != nil && contactBrief.Gender != request.Gender {
			updateContactBriefField("gender", contact.Data.Gender)
		}
	}

	if request.AgeGroup != "" {
		if request.AgeGroup != contact.Data.AgeGroup {
			updatedContactFields = append(updatedContactFields, "ageGroup")
			contact.Data.AgeGroup = request.AgeGroup
			params.ContactUpdates = append(params.ContactUpdates, update.ByFieldName("ageGroup", contact.Data.AgeGroup))
		}
		if contactBrief != nil && contactBrief.AgeGroup != request.AgeGroup {
			updateContactBriefField("ageGroup", contact.Data.AgeGroup)
		}
	}

	if request.Roles != nil {
		var contactFieldsUpdated []string
		if contactFieldsUpdated, err = updateContactRoles(params, *request.Roles); err != nil {
			return err
		}
		updatedContactFields = append(updatedContactFields, contactFieldsUpdated...)
	}

	if request.Related != nil {
		itemRef := dbo4linkage.ItemRef{
			Module:     const4contactus.ModuleID,
			Collection: const4contactus.ContactsCollection,
			ItemID:     request.ContactID,
		}
		var recordsUpdates []record.Updates
		userID := params.UserID()
		recordsUpdates, err = facade4linkage.UpdateRelatedFields(ctx, tx,
			params.Started,
			userID,
			request.SpaceID,
			itemRef, request.UpdateRelatedFieldRequest,
			&dbo4linkage.WithRelatedAndIDsAndUserID{
				WithUserID: dbmodels.WithUserID{
					UserID: params.Contact.Data.UserID,
				},
				WithRelatedAndIDs: &params.Contact.Data.WithRelatedAndIDs,
			},
			func(updates []update.Update) {
				params.ContactUpdates = append(params.ContactUpdates, updates...)
			})
		if err != nil {
			return err
		}
		params.RecordUpdates = append(params.RecordUpdates, recordsUpdates...)
	}

	if len(params.ContactUpdates) > 0 {
		contact.Data.IncreaseVersion(params.Started, params.UserID())
		params.ContactUpdates = append(params.ContactUpdates, contact.Data.GetUpdates()...)
		if err := contact.Data.Validate(); err != nil {
			return fmt.Errorf("contact DBO is not valid after updating %d fields (%+v) and before storing changes DB: %w",
				len(updatedContactFields), updatedContactFields, err)
		}
		if err := tx.Update(ctx, contact.Key, params.ContactUpdates); err != nil {
			return err
		}
	}

	return nil
}
