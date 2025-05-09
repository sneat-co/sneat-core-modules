package dal4contactus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type ContactWorkerParams struct {
	*ContactusSpaceWorkerParams
	Contact        ContactEntry
	ContactUpdates []update.Update
}

func (v ContactWorkerParams) GetRecords(ctx context.Context, tx dal.ReadSession, records ...dal.Record) error {
	return v.ContactusSpaceWorkerParams.GetRecords(ctx, tx, append(records, v.Contact.Record)...)
}

func NewContactWorkerParams(moduleParams *ContactusSpaceWorkerParams, contactID string) *ContactWorkerParams {
	return &ContactWorkerParams{
		ContactusSpaceWorkerParams: moduleParams,
		Contact:                    NewContactEntry(moduleParams.Space.ID, contactID),
	}
}

type ContactWorker = func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *ContactWorkerParams) (err error)

func RunContactWorker(
	ctx facade.ContextWithUser,
	request dto4contactus.ContactRequest,
	worker ContactWorker,
) error {
	contactWorker := func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, moduleWorkerParams *ContactusSpaceWorkerParams) (err error) {
		params := NewContactWorkerParams(moduleWorkerParams, request.ContactID)
		if err = worker(ctx, tx, params); err != nil {
			return err
		}
		if err = applyContactUpdates(ctx, tx, params); err != nil {
			return err
		}
		return err
	}
	return RunContactusSpaceWorker(ctx, request.SpaceRequest, contactWorker)
}

func applyContactUpdates(ctx context.Context, tx dal.ReadwriteTransaction, params *ContactWorkerParams) (err error) {
	if len(params.ContactUpdates) > 0 {
		if err = tx.Update(ctx, params.Contact.Record.Key(), params.ContactUpdates); err != nil {
			return err
		}
	}
	return err
}
