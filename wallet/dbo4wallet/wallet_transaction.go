package dbo4wallet

import (
	"strconv"
	"time"

	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
)

//type OwnerType string
//
//const (
//	OwnerTypeUser  = OwnerType("user")
//	OwnerTypeSpace = OwnerType("space")
//)

type TransactionType string

const (
	TransactionTypeCredit TransactionType = "credit"
	TransactionTypeDebit  TransactionType = "debit"
)

const TopupOperation = "topup"

type WalletTransactionDbo struct {
	with.CreatedAtField
	RelatedTransactionID int64           `json:"relatedTransactionID,omitempty" firestore:"relatedTransactionID,omitempty"`
	Type                 TransactionType `json:"type" firestore:"type"`
	Operation            string          `json:"operation" firestore:"operation"`
	Currency             string          `json:"currency" firestore:"currency"`
	Amount               int             `json:"amount" firestore:"amount"`
	Balance              int             `json:"balance" firestore:"balance"`
	//
	AppPlatform string `json:"appPlatform,omitempty" firestore:"appPlatform,omitempty"`
	AppID       string `json:"appID,omitempty" firestore:"appID,omitempty"`
	//
	MessengerChargeID       string `json:"messengerChargeID,omitempty" firestore:"messengerChargeID,omitempty"`
	ProviderPaymentChargeID string `json:"providerPaymentChargeID,omitempty" firestore:"providerPaymentChargeID,omitempty"`
	//
	IsRefunded bool       `json:"isRefunded" firestore:"isRefunded"` // intentionally do not omitempty
	RefundedAt *time.Time `json:"refundedAt,omitempty" firestore:"refundedAt,omitempty"`
}

func (v *WalletTransactionDbo) Validate() error {
	if v.Amount == 0 {
		return validation.NewErrRecordIsMissingRequiredField("amount")
	}
	switch v.Type {
	case TransactionTypeCredit:
		if v.Amount < 0 {
			return validation.NewErrBadRecordFieldValue("amount", "negative value for credit transaction")
		}
	case TransactionTypeDebit: // OK
		if v.Amount > 0 {
			return validation.NewErrBadRecordFieldValue("amount", "positive value for debit transaction")
		}
	case "":
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	switch currencyLength := len(v.Currency); currencyLength {
	case 0:
		return validation.NewErrRecordIsMissingRequiredField("currency")
	case 3:
		switch v.Currency {
		case // Supported currencies
			"XTR", // XTR is for Telegram stars
			"EUR":
		default:
			return validation.NewErrBadRecordFieldValue("currency", "unsupported value: "+v.Currency)
		}
	default:
		return validation.NewErrBadRecordFieldValue("currency", "expected to be 3 characters long, got "+strconv.Itoa(currencyLength))
	}
	//switch v.Operation {
	//case TopupOperation:
	//	if v.Type == TransactionTypeDebit {
	//		return validation.NewErrBadRecordFieldValue("type|operation", "operation=topup should be a credit transaction")
	//	}
	//case "":
	//	return validation.NewErrRecordIsMissingRequiredField("operation")
	//default:
	//	return validation.NewErrBadRecordFieldValue("operation", "expected to be topup or topup")
	//}
	return v.CreatedAtField.Validate()
}
