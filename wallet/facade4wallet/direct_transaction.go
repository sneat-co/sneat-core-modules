package facade4wallet

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/userus/dal4userus"
	"github.com/sneat-co/sneat-core-modules/wallet/dal4wallet"
	"github.com/sneat-co/sneat-core-modules/wallet/dbo4wallet"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
	"time"
)

type DirectPaymentRequest struct {
	Operation               string `json:"operation"`
	Currency                string `json:"currency"`
	Amount                  int    `json:"amount"`
	InvoicePayload          string `json:"invoicePayload"`
	AppPlatform             string `json:"appPlatform"`
	AppID                   string `json:"appID"`
	MessengerChargeID       string `json:"messengerChargeID,omitempty"`
	PaymentProviderChargeID string `json:"paymentProviderChargeID,omitempty"`
}

func (v DirectPaymentRequest) Validate() error {
	if v.Currency == "" {
		return validation.NewErrRequestIsMissingRequiredField("currency")
	}
	if v.Amount == 0 {
		return validation.NewErrRequestIsMissingRequiredField("amount")
	}
	if v.Amount < 0 {
		return validation.NewErrBadRequestFieldValue("amount", "must be positive")
	}
	if v.AppPlatform == "" {
		return validation.NewErrRequestIsMissingRequiredField("appPlatform")
	}
	if v.AppID == "" {
		return validation.NewErrRequestIsMissingRequiredField("appID")
	}
	return nil
}

type DirectPaymentResponse struct {
}

func RecordDirectPayment(
	ctx facade.ContextWithUser,
	request DirectPaymentRequest,
) (
	resp DirectPaymentResponse, err error,
) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4userus.RunUserWorker(ctx, true, func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, userWorkerParams *dal4userus.UserWorkerParams) (err error) {
		if resp, err = recordDirectPaymentTx(ctx, tx, userWorkerParams, request); err != nil {
			return
		}
		return err
	})
	return
}

func recordDirectPaymentTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userWorkerParams *dal4userus.UserWorkerParams,
	request DirectPaymentRequest,
) (
	resp DirectPaymentResponse, err error,
) {
	walletItem := userWorkerParams.User.Data.Wallet[request.Currency]

	topupTransactionDbo, topupTransactionRecord := createTransactionForDirectPayment(dbo4wallet.TransactionTypeCredit, userWorkerParams.User.ID, request, request.Amount, walletItem.Balance+request.Amount)

	time.Sleep(time.Microsecond) // Guarantees IDs will be different

	withdrawTransactionDbo, withdrawTransactionRecord := createTransactionForDirectPayment(dbo4wallet.TransactionTypeDebit, userWorkerParams.User.ID, request, -request.Amount, walletItem.Balance)

	topupTransactionDbo.RelatedTransactionID = withdrawTransactionRecord.Key().ID.(int64)
	withdrawTransactionDbo.RelatedTransactionID = topupTransactionRecord.Key().ID.(int64)

	transactions := []dal.Record{
		topupTransactionRecord,
		withdrawTransactionRecord,
	}
	if err = tx.InsertMulti(ctx, transactions); err != nil {
		return
	}
	return
}

func createTransactionForDirectPayment(
	transactionType dbo4wallet.TransactionType,
	userID string,
	request DirectPaymentRequest,
	amount, balance int,
) (dbo *dbo4wallet.WalletTransactionDbo, r dal.Record) {
	dbo = &dbo4wallet.WalletTransactionDbo{
		Type:                    transactionType,
		Operation:               request.Operation,
		Currency:                request.Currency,
		Amount:                  amount,
		Balance:                 balance,
		AppPlatform:             request.AppPlatform,
		AppID:                   request.AppID,
		MessengerChargeID:       request.MessengerChargeID,
		ProviderPaymentChargeID: request.PaymentProviderChargeID,
	}
	r = dal4wallet.NewTransactionRecordWithTimestampID(userID, dbo)
	dbo.CreatedAt = time.UnixMicro(r.Key().ID.(int64))
	return
}
