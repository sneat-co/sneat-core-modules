package facade4userus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/userus/dal4userus"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
)

func SetUserTimezone(
	ctx facade.ContextWithUser, ianaTimezone string, utcOffset string,
) (
	user dbo4userus.UserEntry, err error,
) {
	if ctx == nil {
		panic("ctx cannot be nil")
	}
	err = dal4userus.RunUserWorker(ctx, true, func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, userWorkerParams *dal4userus.UserWorkerParams) error {
		if user = userWorkerParams.User; user.Data.Timezone == nil || user.Data.Timezone.Iana != ianaTimezone {
			userWorkerParams.UserUpdates = user.Data.SetTimezone(ianaTimezone, utcOffset)
		}
		return nil
	})
	return
}
