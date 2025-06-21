package facade4wallet

import (
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/userus/dal4userus"
	"github.com/sneat-co/sneat-core-modules/wallet/botcharges"
	"github.com/sneat-co/sneat-core-modules/wallet/const4wallet"
	"github.com/sneat-co/sneat-core-modules/wallet/dbo4wallet"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
	"reflect"
	"time"
)

type TopupUserWalletRequest struct {
	Currency                string `json:"currency"`
	Amount                  int    `json:"amount"`
	InvoicePayload          string `json:"invoicePayload,omitempty"`
	MessengerChargeID       string `json:"messengerChargeID,omitempty"`
	PaymentProviderChargeID string `json:"paymentProviderChargeID,omitempty"`
	BotPlatform             string `json:"botPlatform,omitempty"`
	BotCode                 string `json:"botCode,omitempty"`
}

func (v TopupUserWalletRequest) Validate() error {
	if v.Currency == "" {
		return validation.NewErrRequestIsMissingRequiredField("currency")
	}
	if v.Amount == 0 {
		return validation.NewErrRequestIsMissingRequiredField("amount")
	}
	if v.PaymentProviderChargeID == "" {
		return validation.NewErrRequestIsMissingRequiredField("paymentProviderChargeID")
	}
	return nil
}

var ErrAlreadyProcessed = errors.New("transaction is already processed")

func TopupUserWallet(ctx facade.ContextWithUser, request TopupUserWalletRequest) (transactionID string, balance int, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4userus.RunUserWorker(ctx, true, func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, userWorkerParams *dal4userus.UserWorkerParams) (err error) {
		if transactionID, balance, err = topupUserWalletTxWorker(ctx, tx, userWorkerParams, request); err != nil {
			return
		}
		return err
	})
	return
}

func topupUserWalletTxWorker(
	ctx facade.ContextWithUser,
	tx dal.ReadwriteTransaction,
	userWorkerParams *dal4userus.UserWorkerParams,
	request TopupUserWalletRequest,
) (transactionID string, balance int, err error) {
	if userWorkerParams.User.Data.Wallet == nil {
		userWorkerParams.User.Data.Wallet = make(map[string]dbo4wallet.WalletItem, 1)
	}
	walletItem := userWorkerParams.User.Data.Wallet[request.Currency]

	lastTransaction := walletItem.LastTransaction

	if request.MessengerChargeID != "" && lastTransaction.MessengerChargeID == request.MessengerChargeID {
		err = fmt.Errorf("%w: MessengerChargeID=%s", ErrAlreadyProcessed, request.MessengerChargeID)
		return
	}
	if request.PaymentProviderChargeID != "" && lastTransaction.ProviderPaymentChargeID == request.PaymentProviderChargeID {
		err = fmt.Errorf("%w: PaymentProviderChargeID=%s", ErrAlreadyProcessed, request.PaymentProviderChargeID)
		return
	}

	chargeID := request.PaymentProviderChargeID
	if chargeID == "" {
		chargeID = request.MessengerChargeID
	}
	if request.BotPlatform != "" {
		if request.BotCode == "" {
			err = validation.NewErrRequestIsMissingRequiredField("botCode")
			return
		}
		var isProcessedCharge bool
		if isProcessedCharge, err = botcharges.IsProcessedCharge(ctx, tx, request.BotPlatform, request.BotCode, chargeID); err != nil && !dal.IsNotFound(err) {
			err = fmt.Errorf("failed to check if charge is processed: %w", err)
			return
		}
		if isProcessedCharge {
			err = ErrAlreadyProcessed
			return
		}
	}

	walletItem.Balance = walletItem.Balance + request.Amount
	walletItem.LastTransaction.MessengerChargeID = request.MessengerChargeID
	walletItem.LastTransaction.ProviderPaymentChargeID = request.PaymentProviderChargeID
	userWorkerParams.User.Data.Wallet[request.Currency] = walletItem
	userWorkerParams.UserUpdates = append(userWorkerParams.UserUpdates,
		update.ByFieldPath([]string{"wallet", request.Currency}, walletItem))

	walletTransaction := dbo4wallet.WalletTransactionDbo{
		Type:                    dbo4wallet.TransactionTypeCredit,
		Operation:               dbo4wallet.TopupOperation,
		Currency:                request.Currency,
		Amount:                  request.Amount,
		Balance:                 walletItem.Balance,
		BotPlatform:             request.BotPlatform,
		BotCode:                 request.BotCode,
		MessengerChargeID:       request.MessengerChargeID,
		ProviderPaymentChargeID: request.PaymentProviderChargeID,
	}
	walletTransaction.CreatedAt = time.Now()

	userID := ctx.User().GetUserID()
	userWalletKey := dal4userus.NewUserExtKey(userID, const4wallet.ModuleID)
	transactionKey := dal.NewIncompleteKey(const4wallet.TransactionsCollection, reflect.String, userWalletKey)
	transactionRecord := dal.NewRecordWithData(transactionKey, &walletTransaction)
	if err = tx.Insert(ctx, transactionRecord, dal.WithRandomStringKeyPrefixedByUnixTime(1, 3)); err != nil {
		return
	}
	transactionID = transactionKey.ID.(string)

	if request.BotCode != "" {
		if err = botcharges.SaveCharge(ctx, tx, request.BotPlatform, request.BotCode, &botcharges.ChargeDbo{
			Currency:            request.Currency,
			TotalAmount:         request.Amount,
			InvoicePayload:      request.InvoicePayload,
			MessengerChargeID:   request.MessengerChargeID,
			ProviderChargeID:    request.PaymentProviderChargeID,
			UserID:              userWorkerParams.User.ID,
			WalletTransactionID: transactionID,
		}); err != nil {
			return
		}
	}
	balance = walletItem.Balance
	userWorkerParams.User.Record.MarkAsChanged()
	return
}
