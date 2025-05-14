package facade4contactus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/strongoapp/with"
)

func AddPhone(ctx facade.ContextWithUser, request dto4contactus.AddPhoneRequest) (err error) {
	return dal4contactus.RunContactWorker(ctx, request.ContactRequest, func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4contactus.ContactWorkerParams) (err error) {
		return addPhoneWorker(ctx, tx, params, request)
	})
}

func addPhoneWorker(
	ctx facade.ContextWithUser,
	tx dal.ReadwriteTransaction,
	params *dal4contactus.ContactWorkerParams,
	request dto4contactus.AddPhoneRequest,
) (err error) {
	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}
	phoneKey := request.PhoneNumber
	if _, ok := params.Contact.Data.Phones[phoneKey]; ok {
		return nil
	}
	phoneProps := with.PhoneProps{
		Type: request.Type,
	}
	params.Contact.Data.Phones[phoneKey] = phoneProps
	params.Contact.Record.MarkAsChanged()
	return err
}

func DeletePhone(ctx facade.ContextWithUser, request dto4contactus.DeletePhoneRequest) (err error) {
	return dal4contactus.RunContactWorker(ctx, request.ContactRequest, func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4contactus.ContactWorkerParams) (err error) {
		return deletePhoneWorker(ctx, tx, params, request)
	})
}

func deletePhoneWorker(
	ctx facade.ContextWithUser,
	tx dal.ReadwriteTransaction,
	params *dal4contactus.ContactWorkerParams,
	request dto4contactus.DeletePhoneRequest,
) (err error) {
	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}
	phoneKey := request.PhoneNumber
	if _, ok := params.Contact.Data.Phones[phoneKey]; !ok {
		return nil
	}
	delete(params.Contact.Data.Phones, phoneKey)
	params.Contact.Record.MarkAsChanged()
	return err
}
