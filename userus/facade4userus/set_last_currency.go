package facade4userus

import (
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-core/facade"
)

func SetLastCurrency(ctx facade.ContextWithUser, currencyCode money.CurrencyCode) (err error) {
	return dal4userus.RunUserWorker(ctx, true,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, userWorkerParams *dal4userus.UserWorkerParams) (err error) {
			return setLastCurrency(userWorkerParams, currencyCode)
		})
}

func setLastCurrency(params *dal4userus.UserWorkerParams, currencyCode money.CurrencyCode) (err error) {
	params.UserUpdates, err = params.User.Data.SetLastCurrency(currencyCode)
	return
}
