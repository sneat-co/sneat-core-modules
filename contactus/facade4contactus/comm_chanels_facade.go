package facade4contactus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/strongoapp/with"
	"strings"
)

func AddCommChannel(ctx facade.ContextWithUser, request dto4contactus.AddCommChannelRequest) (err error) {
	return dal4contactus.RunContactWorker(ctx, request.ContactRequest, func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4contactus.ContactWorkerParams) (err error) {
		return addCommChannelWorker(ctx, tx, params, request)
	})
}

func addCommChannelWorker(
	ctx facade.ContextWithUser,
	tx dal.ReadwriteTransaction,
	params *dal4contactus.ContactWorkerParams,
	request dto4contactus.AddCommChannelRequest,
) (err error) {
	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}
	channelID := strings.ToLower(request.ChannelID)
	channels, fieldName := params.Contact.Data.GetCommChannels(request.ChannelType)
	if _, ok := channels[channelID]; ok {
		return nil
	}
	props := with.CommunicationChannelProps{
		Type: request.Type,
		CreatedFields: with.CreatedFields{
			CreatedAtField: with.CreatedAtField{
				CreatedAt: params.Started,
			},
			CreatedByField: with.CreatedByField{
				CreatedBy: ctx.User().GetUserID(),
			},
		},
	}

	if channelID != request.ChannelID {
		props.Original = request.ChannelID
	}

	channels[channelID] = &props
	params.ContactUpdates = append(params.ContactUpdates,
		update.ByFieldPath([]string{fieldName, channelID}, props))
	params.Contact.Record.MarkAsChanged()
	return err
}

func UpdateCommChannel(ctx facade.ContextWithUser, request dto4contactus.UpdateCommChannelRequest) (err error) {
	return dal4contactus.RunContactWorker(ctx, request.ContactRequest, func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4contactus.ContactWorkerParams) (err error) {
		if err = params.GetRecords(ctx, tx); err != nil {
			return err
		}
		return updateCommChannelWorker(params, request)
	})
}

func updateCommChannelWorker(
	params *dal4contactus.ContactWorkerParams,
	request dto4contactus.UpdateCommChannelRequest,
) (err error) {
	channelID := strings.ToLower(request.ChannelID)
	channels, fieldName := params.Contact.Data.GetCommChannels(request.ChannelType)

	props := channels[channelID]
	if props == nil {
		return fmt.Errorf("contact has no %s %s", request.ChannelType, channelID)
	}

	if request.Type != nil {
		if t := *request.Type; t != props.Type {
			props.Type = *request.Type
			params.ContactUpdates = append(params.ContactUpdates, update.ByFieldPath([]string{fieldName, channelID, "type"}, props.Type))
			params.Contact.Record.MarkAsChanged()
		}
	}
	if request.Note != nil {
		if note := *request.Note; note != props.Note {
			props.Note = *request.Note
			params.ContactUpdates = append(params.ContactUpdates, update.ByFieldPath([]string{fieldName, channelID, "note"}, props.Note))
			params.Contact.Record.MarkAsChanged()
		}
	}
	if request.NewChannelID != nil {
		if newChannelID := strings.ToLower(*request.NewChannelID); newChannelID != channelID {
			if newChannelID != request.ChannelID {
				props.Original = request.ChannelID
			}
			// Intentionally overrider slice with a single update
			params.ContactUpdates = []update.Update{
				update.ByFieldPath([]string{fieldName, channelID}, update.DeleteField),
				update.ByFieldPath([]string{fieldName, *request.NewChannelID}, props),
			}
			params.Contact.Record.MarkAsChanged()
		}
	}
	return err
}

func DeleteCommChannel(ctx facade.ContextWithUser, request dto4contactus.DeleteCommChannelRequest) (err error) {
	return dal4contactus.RunContactWorker(ctx, request.ContactRequest, func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4contactus.ContactWorkerParams) (err error) {
		if err = params.GetRecords(ctx, tx); err != nil {
			return err
		}
		return deleteCommChannelWorker(params, request)
	})
}

func deleteCommChannelWorker(
	params *dal4contactus.ContactWorkerParams,
	request dto4contactus.DeleteCommChannelRequest,
) (err error) {
	channelID := strings.ToLower(request.ChannelID)
	channels, fieldName := params.Contact.Data.GetCommChannels(request.ChannelType)
	if _, ok := channels[channelID]; !ok {
		return nil
	}

	delete(channels, channelID)
	if len(channels) == 0 {
		params.ContactUpdates = append(params.ContactUpdates,
			update.ByFieldName(fieldName, update.DeleteField))
	} else {
		params.ContactUpdates = append(params.ContactUpdates,
			update.ByFieldPath([]string{fieldName, channelID}, update.DeleteField))
	}
	params.Contact.Record.MarkAsChanged()
	return err
}
