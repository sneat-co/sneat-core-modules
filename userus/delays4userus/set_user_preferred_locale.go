package delays4userus

import (
	"context"

	"github.com/sneat-co/sneat-core-modules/userus/facade4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
)

func delayedSetUserPreferredLocale(ctx context.Context, userID string, localeCode5 string) (err error) {
	logus.Debugf(ctx, "delayedSetUserPreferredLocale(userID=%v, localeCode5=%v)", userID, localeCode5)
	ctxWithUser := facade.NewContextWithUserID(ctx, userID)
	return facade4userus.SetUserPreferredLocale(ctxWithUser, localeCode5)
}
