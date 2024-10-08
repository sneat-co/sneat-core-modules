package facade4auth

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/auth/models4auth"
)

func NewUserGoogleKey(id string) *dal.Key {
	return dal.NewKeyWithID(models4auth.UserGoogleCollection, id)
}

type UserGoogleDalGae struct {
}

func NewUserGoogleDalGae() UserGoogleDalGae {
	return UserGoogleDalGae{}
}

func (UserGoogleDalGae) GetUserGoogleByID(_ context.Context, googleUserID string) (userGoogle models4auth.UserAccountEntry, err error) {
	//userGoogle.ContactID = googleUserID
	//userGoogle.Data = new(models.UserGoogleData)
	//if err = gaedb.Get(ctx, NewUserGoogleKey(googleUserID), userGoogle.Data); err != nil {
	//	if err == datastore.ErrNoSuchEntity {
	//		err = dal.ErrRecordNotFound
	//	}
	//	return
	//}
	err = errors.New("not implemented")
	return
}

func (UserGoogleDalGae) DeleteUserGoogle(_ context.Context, googleUserID string) (err error) {
	//if err = gaedb.Delete(ctx, NewUserGoogleKey(googleUserID)); err != nil {
	//	return
	//}
	return errors.New("not implemented")
}
