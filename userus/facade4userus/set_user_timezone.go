package facade4userus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/userus/dal4userus"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"time"
)

func SetUserTimezone(
	ctx facade.ContextWithUser, ianaTimezone string,
) (
	user dbo4userus.UserEntry, err error,
) {
	if ctx == nil {
		panic("ctx cannot be nil")
	}
	var offsetMinutes int
	if offsetMinutes, err = getOffsetMinutes(ianaTimezone, time.Now()); err != nil {
		return
	}

	err = dal4userus.RunUserWorker(ctx, true, func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, userWorkerParams *dal4userus.UserWorkerParams) error {
		if user = userWorkerParams.User; user.Data.Timezone == nil || user.Data.Timezone.Iana != ianaTimezone {
			userWorkerParams.UserUpdates = user.Data.SetTimezone(ianaTimezone, offsetMinutes)
			if len(userWorkerParams.UserUpdates) > 0 {
				user.Record.MarkAsChanged()
			}
		}
		return nil
	})
	return
}

func getOffsetMinutes(locName string, t time.Time) (int, error) {
	loc, err := time.LoadLocation(locName)
	if err != nil {
		return 0, err
	}
	_, offsetSeconds := t.In(loc).Zone()
	return offsetSeconds / 60, nil // Convert to minutes
}
