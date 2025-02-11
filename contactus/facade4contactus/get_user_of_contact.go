package facade4contactus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/strongoapp/appuser"
	"reflect"
	"strings"
)

type UserAccountsProvider interface { // Should we move it into appuser package?
	GetUserID() string
	GetAccount(provider, app string) (userAccount *appuser.AccountKey, err error)
}

func GetUserOfContact(ctx context.Context, contact UserAccountsProvider, platform, app string) (user dbo4userus.UserEntry, err error) {
	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return
	}
	if userID := contact.GetUserID(); userID != "" {
		user = dbo4userus.NewUserEntry(userID)
		err = db.Get(ctx, user.Record)
		return
	}

	var account *appuser.AccountKey
	if account, err = contact.GetAccount(platform, app); err != nil {
		return
	}
	if account == nil {
		return
	}
	usersCollection := dal.NewCollectionRef(dbo4userus.UsersCollection, "", nil)
	qb := dal.From(usersCollection).WhereField("accounts", "array-contains", account.String())
	q := qb.SelectInto(func() dal.Record {
		return dal.NewRecordWithIncompleteKey(dbo4userus.UsersCollection, reflect.String, new(dbo4userus.UserDbo))
	})
	var records []dal.Record
	if records, err = db.QueryAllRecords(ctx, q); err != nil {
		return
	}
	switch count := len(records); count {
	case 0:
		return
	case 1:
		r := records[0]
		user = dbo4userus.NewUserEntryWithDbo(r.Key().ID.(string), r.Data().(*dbo4userus.UserDbo))
		return
	default:
		var userIDs []string
		for _, r := range records {
			userIDs = append(userIDs, r.Key().ID.(string))
		}
		err = fmt.Errorf("%w: %d records by account=%s: userID=%s",
			ErrTooManuUsersFound, len(records), account, strings.Join(userIDs, ","))
		return
	}
}

var ErrTooManuUsersFound = errors.New("too many users found")
