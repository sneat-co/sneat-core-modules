package facade4contactus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
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
	updateContactWorker := func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4contactus.ContactWorkerParams) (err error) {
		contact = params.Contact
		contactusSpace = params.SpaceModuleEntry
		space = params.Space
		return UpdateContactTx(ctx, tx, request, params)
	}
	err = dal4contactus.RunContactWorker(ctx, request.ContactRequest, updateContactWorker)
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

	if request.Names != nil { // TODO: move this into `person` package.
		fullName := request.Names.FullName
		if fullName != "" {
			names := *request.Names
			names.FullName = ""
			if fullName == request.Names.GetFullName() {
				request.Names.FullName = ""
			}
		}
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

	if err = contact.Data.Validate(); err != nil {
		return fmt.Errorf("contact DBO is not valid after loading from DB: %w", err)
	}

	contactBrief := params.SpaceModuleEntry.Data.GetContactBriefByContactID(request.ContactID)

	updateContactBriefField := func(field string, value any) {
		params.SpaceModuleEntry.Record.MarkAsChanged()
		params.SpaceModuleUpdates = append(params.SpaceModuleUpdates,
			update.ByFieldName(
				fmt.Sprintf("contacts.%s.%s", request.ContactID, field),
				value))
	}

	var updatedContactFields []string

	if request.Names != nil {
		names := *request.Names
		if contact.Data.Names == nil || *contact.Data.Names != names {
			updatedContactFields = append(updatedContactFields, "names")
			contact.Data.Names = request.Names
			params.ContactUpdates = append(params.ContactUpdates, update.ByFieldName("names", request.Names))
		}
		if contactBrief != nil && *contactBrief.Names != names {
			contactBrief.Names = &names
			updateContactBriefField("names", contact.Data.Names)
		}
	}
	if request.Address != nil && (contact.Data.Address == nil || *request.Address != *contact.Data.Address) {
		updatedContactFields = append(updatedContactFields, "address")
		contact.Data.Address = request.Address
		params.ContactUpdates = append(params.ContactUpdates, update.ByFieldName("address", request.Address))
	}

	if request.VatNumber != nil {
		if vat := *request.VatNumber; vat != contact.Data.VATNumber {
			updatedContactFields = append(updatedContactFields, "vatNumber")
			contact.Data.VATNumber = vat
			params.ContactUpdates = append(params.ContactUpdates, update.ByFieldName("vatNumber", vat))
		}
	}

	if request.Gender != "" {
		if request.Gender != contact.Data.Gender {
			updatedContactFields = append(updatedContactFields, "gender")
			contact.Data.Gender = request.Gender
			params.ContactUpdates = append(params.ContactUpdates, update.ByFieldName("gender", request.Gender))
		}
		if contactBrief != nil && contactBrief.Gender != request.Gender {
			contactBrief.Gender = request.Gender
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
			contactBrief.AgeGroup = request.AgeGroup
			updateContactBriefField("ageGroup", contact.Data.AgeGroup)
		}
	}

	if request.DateOfBirth != nil {
		dob := *request.DateOfBirth
		if dob != contact.Data.DoB {
			updatedContactFields = append(updatedContactFields, "dob")
			contact.Data.DoB = dob
			params.ContactUpdates = append(params.ContactUpdates,
				update.ByFieldName("dob", dob))
		}
		if contactBrief != nil && contactBrief.DoB != dob {
			contactBrief.DoB = dob
			updateContactBriefField("dob", contact.Data.DoB)
		}
	}

	if request.Roles != nil {
		var contactFieldsUpdated []string
		if contactFieldsUpdated, err = updateContactRoles(params, *request.Roles); err != nil {
			return err
		}
		if len(contactFieldsUpdated) > 0 {
			updatedContactFields = append(updatedContactFields, contactFieldsUpdated...)
			if contactBrief != nil {
				contactBrief.Roles = contact.Data.Roles
				updateContactBriefField("roles", contact.Data.Roles)
			}
		}
	}

	if len(params.ContactUpdates) > 0 {
		contact.Data.IncreaseVersion(params.Started, params.UserID())
		params.ContactUpdates = append(params.ContactUpdates, contact.Data.GetUpdates()...)
		if err = contact.Data.Validate(); err != nil {
			return fmt.Errorf("contact DBO is not valid after updating %d fields (%+v) and before storing changes DB: %w",
				len(updatedContactFields), updatedContactFields, err)
		}
		if err = tx.Update(ctx, contact.Key, params.ContactUpdates); err != nil {
			return err
		}
		if contactBrief == nil {
			params.SpaceModuleUpdates = append(params.SpaceModuleUpdates,
				params.SpaceModuleEntry.Data.SetContactBrief(contact.ID, &contact.Data.ContactBrief)...)
			params.SpaceModuleEntry.Record.MarkAsChanged()
		}
	}

	return nil
}
