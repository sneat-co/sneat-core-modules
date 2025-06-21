package dbo4wallet

import (
	"errors"
	"fmt"
	"github.com/strongo/validation"
	"strings"
)

const EUR = "EUR"
const TelegramStarsCurrency = "XTR"

type WithWallet struct {
	Wallet map[string]WalletItem `json:"wallet,omitempty" firestore:"wallet,omitempty"`
}

type LastTransaction struct {
	MessengerChargeID       string `json:"messengerChargeID,omitempty" firestore:"messengerChargeID,omitempty"`
	ProviderPaymentChargeID string `json:"providerPaymentChargeID,omitempty" firestore:"providerPaymentChargeID,omitempty"`
}

func (v LastTransaction) Validate() error {
	if v.MessengerChargeID != "" && strings.TrimSpace(v.MessengerChargeID) != v.MessengerChargeID {
		return validation.NewErrBadRecordFieldValue("messengerChargeID", "has leading or trailing spaces")
	}
	if v.ProviderPaymentChargeID != "" && strings.TrimSpace(v.ProviderPaymentChargeID) != v.ProviderPaymentChargeID {
		return validation.NewErrBadRecordFieldValue("providerPaymentChargeID", "has leading or trailing spaces")
	}
	return nil
}

type WalletItem struct {
	Balance         int             `json:"balance" firestore:"balance"`
	LastTransaction LastTransaction `json:"lastTransaction" firestore:"lastTransaction"`
}

func (v *WalletItem) Validate() error {
	if v.Balance < 0 {
		return validation.NewErrBadRequestFieldValue("balance", "must be positive or zero")
	}
	return v.LastTransaction.Validate()
}

func (w WithWallet) GetBalance(currency string) int {
	return w.Wallet[currency].Balance
}

func (w WithWallet) Validate() error {
	for currency, value := range w.Wallet {
		if err := validateCurrencyCode(currency); err != nil {
			return validation.NewErrBadRecordFieldValue("wallet", err.Error())
		}
		if err := value.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("wallet."+currency, err.Error())
		}
	}
	return nil
}

func validateCurrencyCode(currency string) error {

	switch len(currency) {
	case 3:
		break // OK
	case 0:
		return errors.New("empty currency code")
	default:
		return fmt.Errorf("invalid currency '%s': expected to be 3 characters long", currency)
	}
	if strings.ToUpper(currency) != currency {
		return fmt.Errorf("invalid currency '%s': expected to be 3 characters long in UPPER case", currency)
	}
	switch currency {
	case TelegramStarsCurrency, EUR: // Supported currencies
	default:
		return fmt.Errorf("unsuppported currency: %s", currency)
	}
	return nil
}
