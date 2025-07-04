package dal4wallet

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/userus/dal4userus"
	"github.com/sneat-co/sneat-core-modules/wallet/const4wallet"
	"github.com/sneat-co/sneat-core-modules/wallet/dbo4wallet"
	"time"
)

type TransactionEntry = record.DataWithID[int64, *dbo4wallet.WalletTransactionDbo]

//func NewTransactionRecordWithIncompleteKey(userID string, dbo *dbo4wallet.WalletTransactionDbo) dal.Record {
//	userWalletKey := dal4userus.NewUserExtKey(userID, const4wallet.ModuleID)
//	transactionKey := dal.NewIncompleteKey(const4wallet.TransactionsCollection, reflect.String, userWalletKey)
//	return dal.NewRecordWithData(transactionKey, dbo)
//}

func NewTransactionRecordWithTimestampID(userID string, dbo *dbo4wallet.WalletTransactionDbo) dal.Record {
	userWalletKey := dal4userus.NewUserExtKey(userID, const4wallet.ModuleID)
	id := time.Now().UnixMicro()
	transactionKey := dal.NewKeyWithParentAndID(userWalletKey, const4wallet.TransactionsCollection, id)
	return dal.NewRecordWithData(transactionKey, dbo)
}

func NewTransactionEntryFromRecord(r dal.Record) TransactionEntry {
	key := r.Key()
	id := key.ID.(int64)
	return record.NewDataWithID(id, key, r.Data().(*dbo4wallet.WalletTransactionDbo))
}
