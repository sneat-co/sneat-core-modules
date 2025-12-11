package delays4contactus

import (
	"context"
	"time"

	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/strongo/delaying"
)

func InitDelays4contactus(mustRegisterFunc func(key string, i any) delaying.Delayer) {
	delayerUpdateContactusSpaceDboWithContact = mustRegisterFunc("delayedUpdateContactusSpaceDboWithContact", delayedUpdateContactusSpaceDboWithContact)
}

var (
	delayerUpdateContactusSpaceDboWithContact delaying.Delayer
)

func DelayUpdateContactusSpaceDboWithContact(ctx context.Context, delay time.Duration, userID string, contactID string) error {
	return delayerUpdateContactusSpaceDboWithContact.EnqueueWork(ctx, delaying.With(const4contactus.QueueContacts, "UpdateContactusSpaceDboWithContact", delay), userID, contactID)
}
