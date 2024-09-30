package facade4auth

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	models4auth2 "github.com/sneat-co/sneat-core-modules/auth/models4auth"
)

type UserEmailGaeDal struct {
}

func NewUserEmailGaeDal() UserEmailGaeDal {
	return UserEmailGaeDal{}
}

func (UserEmailGaeDal) GetUserEmailByID(ctx context.Context, tx dal.ReadSession, email string) (userEmail models4auth2.UserEmailEntry, err error) {
	userEmail = models4auth2.NewUserEmail(email, nil)
	return userEmail, tx.Get(ctx, userEmail.Record)
}

func (UserEmailGaeDal) SaveUserEmail(ctx context.Context, tx dal.ReadwriteTransaction, userEmail models4auth2.UserEmailEntry) (err error) {
	return tx.Set(ctx, userEmail.Record)
}