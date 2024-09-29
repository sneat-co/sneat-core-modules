package api4auth

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/auth/token4auth"
	unsorted4auth2 "github.com/sneat-co/sneat-core-modules/auth/unsorted4auth"
	"github.com/sneat-co/sneat-core-modules/common4all"
	"github.com/sneat-co/sneat-core-modules/userus/const4userus"
	dal4userus2 "github.com/sneat-co/sneat-core-modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"github.com/strongo/strongoapp/appuser"
	"net/http"
)

func HandleDisconnect(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
	provider := r.URL.Query().Get("provider")

	userCtx := facade.NewUserContext(authInfo.UserID)

	if err := dal4userus2.RunUserWorker(ctx, userCtx, true, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4userus2.UserWorkerParams) error {
		appUser, err := dal4userus2.GetUserByID(ctx, tx, authInfo.UserID)
		if err != nil {
			return err
		}

		changed := false

		deleteFbUser := func(userAccount appuser.AccountKey) error {
			if userFb, err := unsorted4auth2.UserFacebook.GetFbUserByFbID(ctx, userAccount.App, userAccount.ID); err != nil {
				if !dal.IsNotFound(err) {
					return err
				}
			} else if fbUserAppID := userFb.FbUserData().GetAppUserID(); fbUserAppID == appUser.ID {
				if err = unsorted4auth2.UserFacebook.DeleteFbUser(ctx, userAccount.App, userAccount.ID); err != nil {
					return err
				}
			} else {
				logus.Warningf(ctx, "TODO: Handle case if userFb.AppUserIntID:%s != appUser.ContactID:%s", fbUserAppID, appUser.ID)
			}
			return nil
		}

		if !const4userus.IsKnownUserAccountProvider(provider) {
			common4all.ErrorAsJson(ctx, w, http.StatusBadRequest, errors.New("Unknown provider: "+provider))
			return nil
		}
		if !appUser.Data.HasAccount(provider, "") {
			return nil
		}
		var userAccount *appuser.AccountKey
		switch provider {
		case const4userus.GoogleAuthProvider:
			if userAccount, err = appUser.Data.GetAccount(provider, ""); err != nil {
				return err
			} else if userAccount != nil {
				if userGoogle, err := unsorted4auth2.UserGoogle.GetUserGoogleByID(ctx, userAccount.ID); err != nil {
					if !dal.IsNotFound(err) {
						return err
					}
				} else if userGoogle.Data.GetAppUserID() == appUser.ID {
					userGoogle.Data.SetAppUserID("")
					if err = unsorted4auth2.UserGoogle.DeleteUserGoogle(ctx, userGoogle.ID); err != nil {
						return err
					}
				} else {
					logus.Warningf(ctx, "TODO: Handle case if userGoogle.AppUserIntID:%s != appUser.ContactID:%s", userGoogle.Data.GetAppUserID(), appUser.ID)
				}
				_ = appUser.Data.RemoveAccount(*userAccount)
				changed = true
			}
		case const4userus.FacebookAuthProvider:
			if userAccount, err = appUser.Data.GetAccount("facebook", ""); err != nil {
				return err
			} else if userAccount != nil {
				if err = deleteFbUser(*userAccount); err != nil {
					return err
				}
				_ = appUser.Data.RemoveAccount(*userAccount)
				changed = true
			}
		case const4userus.FacebookMessengerAuthProvider:
			if userAccount, err = appUser.Data.GetAccount("facebook", ""); err != nil {
				return err
			} else if userAccount != nil {
				if err = deleteFbUser(*userAccount); err != nil {
					return err
				}
				_ = appUser.Data.RemoveAccount(*userAccount)
				changed = true
			}
		default:
		}

		if changed {
			appUser.Record.MarkAsChanged()
		}
		return nil
	}); err != nil {
		common4all.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
	}
}
