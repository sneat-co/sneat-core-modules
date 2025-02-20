package facade4invitus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-go-core/facade"
	"reflect"
)

func GetPendingInviteForTarget(ctx context.Context, userID, spaceID, targetType, targetID string) (invite InviteEntry, err error) {
	invitesCollection := dal.NewRootCollectionRef("invites", "i")
	var q = dal.From(invitesCollection).
		WhereField("fromUserID", dal.Equal, userID).
		WhereField("spaceID", dal.Equal, spaceID).
		WhereField("targetType", dal.Equal, targetType).
		WhereInArrayField("targetIDs", targetID).
		WhereField("status", dal.Equal, dbo4invitus.InviteStatusPending).
		Limit(1).
		SelectInto(func() dal.Record {
			return dal.NewRecordWithIncompleteKey(InvitesCollection, reflect.String, new(dbo4invitus.InviteDbo))
		})
	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return
	}
	var records []dal.Record
	records, err = db.QueryAllRecords(ctx, q)
	if len(records) == 1 {
		record := records[0]
		invite.ID = record.Key().ID.(string)
		invite = NewInviteEntryWithDbo(invite.ID, record.Data().(*dbo4invitus.InviteDbo))
	}
	return
}
