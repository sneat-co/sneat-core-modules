package facade4userus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/userus/dal4userus"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"time"
)

func SetUserTimezone(
	ctx facade.ContextWithUser, loc *time.Location,
) (
	user dbo4userus.UserEntry, err error,
) {
	if ctx == nil {
		panic("ctx cannot be nil")
	}
	if loc == nil {
		panic("loc cannot be nil")
	}
	err = dal4userus.RunUserWorker(ctx, true,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, userWorkerParams *dal4userus.UserWorkerParams) error {
			user = userWorkerParams.User
			if userWorkerParams.UserUpdates, err = user.Data.SetTimezone(loc); err != nil {
				return err
			}
			if len(userWorkerParams.UserUpdates) > 0 {
				user.Record.MarkAsChanged()
			}
			return nil
		},
	)
	return
}
