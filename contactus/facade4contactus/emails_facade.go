package facade4contactus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/dbo4all"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/strongoapp/with"
	"strings"
)

func AddEmail(ctx facade.ContextWithUser, request dto4contactus.AddEmailRequest) (err error) {
	return dal4contactus.RunContactWorker(ctx, request.ContactRequest, func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4contactus.ContactWorkerParams) (err error) {
		return addEmailWorker(ctx, tx, params, request)
	})
}

func addEmailWorker(
	ctx facade.ContextWithUser,
	tx dal.ReadwriteTransaction,
	params *dal4contactus.ContactWorkerParams,
	request dto4contactus.AddEmailRequest,
) (err error) {
	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}
	emailKey := strings.ToLower(request.EmailAddress)
	if _, ok := params.Contact.Data.Emails[emailKey]; ok {
		return nil
	}
	emailProps := dbo4all.EmailProps{
		Type:  request.Type,
		Title: request.Type,
		CreatedFields: with.CreatedFields{
			CreatedAtField: with.CreatedAtField{
				CreatedAt: params.Started,
			},
			CreatedByField: with.CreatedByField{
				CreatedBy: ctx.User().GetUserID(),
			},
		},
	}

	if emailKey != request.EmailAddress {
		emailProps.OriginalEmail = request.EmailAddress
	}
	params.Contact.Data.Emails[emailKey] = emailProps
	params.ContactUpdates = append(params.ContactUpdates,
		update.ByFieldPath([]string{dbo4all.EmailsField, emailKey}, emailProps))
	params.Contact.Record.MarkAsChanged()
	return err
}

func DeleteEmail(ctx facade.ContextWithUser, request dto4contactus.DeleteEmailRequest) (err error) {
	return dal4contactus.RunContactWorker(ctx, request.ContactRequest, func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4contactus.ContactWorkerParams) (err error) {
		return deleteEmailWorker(ctx, tx, params, request)
	})
}

func deleteEmailWorker(
	ctx facade.ContextWithUser,
	tx dal.ReadwriteTransaction,
	params *dal4contactus.ContactWorkerParams,
	request dto4contactus.DeleteEmailRequest,
) (err error) {
	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}
	emailKey := strings.ToLower(request.EmailAddress)
	if _, ok := params.Contact.Data.Emails[emailKey]; !ok {
		return nil
	}
	delete(params.Contact.Data.Emails, emailKey)
	if len(params.Contact.Data.Emails) == 0 {
		params.Contact.Data.Emails = nil
		params.ContactUpdates = append(params.ContactUpdates,
			update.ByFieldName(dbo4all.EmailsField, update.DeleteField))
	} else {
		params.ContactUpdates = append(params.ContactUpdates,
			update.ByFieldPath([]string{dbo4all.EmailsField, emailKey}, update.DeleteField))
	}
	params.Contact.Record.MarkAsChanged()
	return err
}
