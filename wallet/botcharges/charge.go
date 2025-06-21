package botcharges

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/strongo/validation"
	"strings"
	"time"
)

type ChargeDbo struct {
	Currency       string `json:"currency" firestore:"currency"`
	TotalAmount    int    `json:"totalAmount" firestore:"totalAmount"`
	InvoicePayload string `json:"invoicePayload" firestore:"invoicePayload"`
	//
	MessengerChargeID string `json:"messengerChargeID,omitempty" firestore:"messengerChargeID,omitempty"`
	ProviderChargeID  string `json:"providerChargeID,omitempty" firestore:"providerChargeID,omitempty"`

	//
	UserID string `json:"userID" firestore:"userID"`
	// WalletTransactionID is a unique identifier for a transaction associated with the wallet, used for tracking and auditing.
	WalletTransactionID string     `json:"walletTransactionID,omitempty" firestore:"walletTransactionID,omitempty"`
	IsRefunded          bool       `json:"isRefunded" firestore:"isRefunded"` // intentionally do not omitempty
	RefundedAt          *time.Time `json:"refundedAt,omitempty" firestore:"refundedAt,omitempty"`
}

func (v *ChargeDbo) Validate() error {
	if strings.TrimSpace(v.Currency) == "" {
		return validation.NewErrRecordIsMissingRequiredField("currency")
	}
	if v.TotalAmount == 0 {
		return validation.NewErrRecordIsMissingRequiredField("totalAmount")
	}
	if v.MessengerChargeID == "" && v.ProviderChargeID == "" {
		return validation.NewErrRecordIsMissingRequiredField("messengerChargeID|providerChargeID")
	}
	if strings.TrimSpace(v.UserID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("userID")
	}
	return nil
}

const chargesCollection = "charges"
const botPlatformsCollection = "botPlatforms"
const botsCollection = "bots"

func newMessengerChargeKey(platform, bot, chargeID string) *dal.Key {
	botPlatformKey := dal.NewKeyWithID(botPlatformsCollection, platform)
	botKey := dal.NewKeyWithParentAndID(botPlatformKey, botsCollection, bot)
	return dal.NewKeyWithParentAndID(botKey, chargesCollection, chargeID)
}

func IsProcessedCharge(ctx context.Context, tx dal.Getter, platform, bot, chargeID string) (bool, error) {
	chargeKey := newMessengerChargeKey(platform, bot, chargeID)
	return tx.Exists(ctx, chargeKey)
}

func SaveCharge(ctx context.Context, tx dal.ReadwriteTransaction, platform, bot string, charge *ChargeDbo) error {
	chargeKey := newMessengerChargeKey(platform, bot, charge.ProviderChargeID)
	record := dal.NewRecordWithData(chargeKey, charge)
	return tx.Insert(ctx, record)
}
