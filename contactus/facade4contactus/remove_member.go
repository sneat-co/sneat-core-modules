package facade4contactus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/userus/dal4userus"
	dbo4userus2 "github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
	"slices"
)

// RemoveSpaceMember removes members from a team
func RemoveSpaceMember(ctx context.Context, userCtx facade.UserContext, request dto4contactus.ContactRequest) (err error) {
	if err = request.Validate(); err != nil {
		return err
	}
	return dal4contactus.RunContactWorker(ctx, userCtx, request,
		func(ctx context.Context, tx dal.ReadwriteTransaction,
			params *dal4contactus.ContactWorkerParams,
		) (err error) {
			return removeSpaceMemberTx(ctx, tx, request, params)
		})
}

func removeSpaceMemberTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	request dto4contactus.ContactRequest,
	params *dal4contactus.ContactWorkerParams,
) (err error) {

	if err = params.GetRecords(ctx, tx); err != nil {
		return err
	}

	if params.Contact.Record.Exists() {
		params.ContactUpdates = append(params.ContactUpdates, params.Contact.Data.RemoveRole(const4contactus.SpaceMemberRoleMember)...)
	}

	var memberUserID string
	var membersCount int

	memberUserID, membersCount, err = removeContactBrief(params)
	if err != nil {
		return
	}

	removeMemberFromSpaceRecord(params.SpaceWorkerParams, memberUserID, membersCount)

	if memberUserID != "" {
		var (
			userRef *dal.Key
		)
		memberUser := dbo4userus2.NewUserEntry(memberUserID)
		if err = dal4userus.GetUser(ctx, tx, memberUser); err != nil {
			return
		}

		update := updateUserRecordOnSpaceMemberRemoved(memberUser.Data, request.SpaceID)
		if update != nil {
			if err = txUpdate(ctx, tx, userRef, []dal.Update{*update}); err != nil {
				return err
			}
		}
	}
	return
}

func updateUserRecordOnSpaceMemberRemoved(user *dbo4userus2.UserDbo, spaceID string) *dal.Update {
	delete(user.Spaces, spaceID)
	user.SpaceIDs = slice.RemoveInPlaceByValue(user.SpaceIDs, spaceID)
	return &dal.Update{
		Field: "spaces",
		Value: user.Spaces,
	}
}

func removeMemberFromSpaceRecord(
	params *dal4spaceus.SpaceWorkerParams,
	contactUserID string,
	membersCount int,
) {
	if contactUserID != "" && slices.Contains(params.Space.Data.UserIDs, contactUserID) {
		params.Space.Data.UserIDs = slice.RemoveInPlaceByValue(params.Space.Data.UserIDs, contactUserID)
		params.SpaceUpdates = append(params.SpaceUpdates, dal.Update{Field: "userIDs", Value: params.Space.Data.UserIDs})
	}
	//if params.Space.Data.NumberOf[dbo4spaceus.NumberOfMembersFieldName] != membersCount {
	//	params.SpaceUpdates = append(params.SpaceUpdates, params.Space.Data.SetNumberOf(dbo4spaceus.NumberOfMembersFieldName, membersCount))
	//}
}

func removeContactBrief(
	params *dal4contactus.ContactWorkerParams,
) (contactUserID string, membersCount int, err error) {

	for id, contactBrief := range params.SpaceModuleEntry.Data.Contacts {
		if id == params.Contact.ID {
			params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, params.SpaceModuleEntry.Data.RemoveContact(id))
			if contactBrief.UserID != "" {
				contactUserID = contactBrief.UserID
				userIDs := slice.RemoveInPlaceByValue(params.SpaceModuleEntry.Data.UserIDs, contactBrief.UserID)
				if len(userIDs) != len(params.SpaceModuleEntry.Data.UserIDs) {
					params.SpaceModuleEntry.Data.UserIDs = userIDs
					params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, dal.Update{Field: "userIDs", Value: userIDs})
				}
			}
			break
		}
	}
	membersCount = len(params.SpaceModuleEntry.Data.GetContactBriefsByRoles(const4contactus.SpaceMemberRoleMember))
	return
}
