package facade4contactus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// ArchiveContact archives space contact - e.g., hides it from the list of contacts
func ArchiveContact(ctx facade.ContextWithUser, request dto4contactus.ContactRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	return dal4contactus.RunContactWorker(ctx, request,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4contactus.ContactWorkerParams) (err error) {
			return archiveContactTxWorker(ctx, tx, params)
		},
	)
}

func archiveContactTxWorker(
	ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactWorkerParams,
) (err error) {
	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}
	if removeContactRoles(params); len(params.ContactUpdates) > 0 {
		if err = tx.Update(ctx, params.Contact.Record.Key(), params.ContactUpdates); err != nil {
			return err
		}
	}
	return err
}
