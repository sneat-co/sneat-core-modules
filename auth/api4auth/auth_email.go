package api4auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/auth/models4auth"
	"github.com/sneat-co/sneat-core-modules/auth/unsorted4auth"
	"github.com/sneat-co/sneat-core-modules/common4all"
	"github.com/sneat-co/sneat-core-modules/userus/dal4userus"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/strongoapp/appuser"
	"net/http"
	"strconv"
	"time"
)

//var (
//	reEmail         = regexp.MustCompile(`.+@.+\.\w+`)
//	ErrInvalidEmail = errors.New("Invalid email")
//)

//func validateEmail(email string) error {
//	if !reEmail.MatchString(email) {
//		return ErrInvalidEmail
//	}
//	return nil
//}

//func HandleSignUpWithEmail(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//	email := strings.TrimSpace(r.PostFormValue("email"))
//	userName := strings.TrimSpace(r.PostFormValue("name"))
//
//	if email == "" {
//		common4all.BadRequestMessage(ctx, w, "Missing required value: email")
//		return
//	}
//
//	if err := validateEmail(email); err != nil {
//		common4all.ErrorAsJson(ctx, w, http.StatusBadRequest, err)
//		return
//	}
//
//	if _, err := unsorted4auth.UserEmail.GetUserEmailByID(ctx, nil, email); err != nil {
//		if !dal.IsNotFound(err) {
//			common4all.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
//			return
//		} else {
//			common4all.ErrorAsJson(ctx, w, http.StatusConflict, sneaterrors.ErrEmailAlreadyRegistered)
//			return
//		}
//	}
//
//	if user, userEmail, err := facade4debtus.User.CreateUserByEmail(ctx, email, userName); err != nil {
//		if errors.Is(err, sneaterrors.ErrEmailAlreadyRegistered) {
//			common4all.ErrorAsJson(ctx, w, http.StatusConflict, err)
//			return
//		} else {
//			common4all.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
//			return
//		}
//	} else {
//		if err = emailing.CreateConfirmationEmailAndQueueForSending(ctx, user, userEmail); err != nil {
//			common4all.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
//			return
//		}
//		ReturnToken(ctx, w, r, user.ID, r.Referer() /*, user.Data.EmailAddress == "alexander.trakhimenok@gmail.com"*/)
//	}
//}
//
//func HandleSignInWithEmail(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//	email := strings.TrimSpace(r.PostFormValue("email"))
//	password := strings.TrimSpace(r.PostFormValue("password"))
//	//logus.Debugf(ctx, "EmailAddress: %s", email)
//	if email == "" || password == "" {
//		common4all.ErrorAsJson(ctx, w, http.StatusBadRequest, errors.New("Missing required value"))
//		return
//	}
//
//	if err := validateEmail(email); err != nil {
//		common4all.JsonToResponse(ctx, w, map[string]string{"error": err.Error()})
//		return
//	}
//
//	userEmail, err := unsorted4auth.UserEmail.GetUserEmailByID(ctx, nil, email)
//	if err != nil {
//		if dal.IsNotFound(err) {
//			common4all.ErrorAsJson(ctx, w, http.StatusForbidden, errors.New("unknown email"))
//		} else {
//			common4all.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
//		}
//		return
//	} else if err = userEmail.Data.CheckPassword(password); err != nil {
//		logus.Debugf(ctx, "Invalid password: %v", err)
//		common4all.ErrorAsJson(ctx, w, http.StatusForbidden, errors.New("invalid password"))
//		return
//	}
//
//	ReturnToken(ctx, w, r, userEmail.Data.AppUserID, r.Referer() /*, userEmail.ID == "alexander.trakhimenok@gmail.com"*/)
//}

func handleRequestPasswordReset(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	userEmail, err := unsorted4auth.UserEmail.GetUserEmailByID(ctx, nil, email)
	if dal.IsNotFound(err) {
		common4all.ErrorAsJson(ctx, w, http.StatusForbidden, errors.New("unknown email"))
		return
	}

	now := time.Now()

	pwdResetEntity := models4auth.PasswordResetData{
		Email:             userEmail.ID,
		Status:            "created",
		OwnedByUserWithID: appuser.NewOwnedByUserWithID(userEmail.Data.AppUserID, now),
	}

	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		_, err = unsorted4auth.PasswordReset.CreatePasswordResetByID(ctx, tx, &pwdResetEntity)
		return err
	})
	if err != nil {
		common4all.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
		return
	}
}

func handleChangePasswordAndSignIn2(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var (
		err           error
		passwordReset models4auth.PasswordReset
	)

	if passwordReset.ID, err = strconv.Atoi(r.PostFormValue("pin")); err != nil {
		common4all.ErrorAsJson(ctx, w, http.StatusBadRequest, err)
		return
	}

	pwd := r.PostFormValue("pwd")
	if pwd == "" {
		common4all.ErrorAsJson(ctx, w, http.StatusBadRequest, errors.New("empty password"))
		return
	}

	if passwordReset, err = unsorted4auth.PasswordReset.GetPasswordResetByID(ctx, nil, passwordReset.ID); err != nil {
		if dal.IsNotFound(err) {
			common4all.ErrorAsJson(ctx, w, http.StatusForbidden, errors.New("unknown pin"))
			return
		}
		common4all.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
		return
	}

	if err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {

		now := time.Now()
		appUser := dbo4userus.NewUserEntry(passwordReset.Data.AppUserID)
		userEmail := models4auth.NewUserEmail(passwordReset.Data.Email, nil)

		records := []dal.Record{appUser.Record, userEmail.Record, passwordReset.Record}

		//var db dal.DB
		//if db, err = facade.GetSneatDB(ctx); err != nil {
		//	return err
		//}
		if err = tx.GetMulti(ctx, records); err != nil {
			return err
		}

		if err = userEmail.Data.SetPassword(pwd); err != nil {
			return err
		}

		passwordReset.Data.Status = "changed"
		passwordReset.Data.Email = "" // Clean email as we don't need it anymore
		passwordReset.Data.UpdatedAt = now
		if changed := userEmail.Data.AddProvider("password-reset"); changed {
			userEmail.Data.UpdatedAt = now
		}
		userEmail.Data.SetLastLoginAt(now)
		appUser.Data.SetLastLoginAt(now)

		if err = tx.SetMulti(ctx, records); err != nil {
			return err
		}
		return err
	}); err != nil {
		common4all.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
		return
	}

	ReturnToken(ctx, w, r, passwordReset.Data.AppUserID, r.Referer())
}

var errInvalidEmailConformationPin = errors.New("email confirmation pin is not valid")

func handleConfirmEmailAndSignIn2(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		userEmail models4auth.UserEmailEntry
		pin       string
	)

	userEmail.ID, pin = r.PostFormValue("email"), r.PostFormValue("pin")

	if userEmail.ID == "" {
		common4all.ErrorAsJson(ctx, w, http.StatusBadRequest, errors.New("empty email"))
		return
	}
	if pin == "" {
		common4all.ErrorAsJson(ctx, w, http.StatusBadRequest, errors.New("empty pin"))
		return
	}

	if err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		now := time.Now()

		if userEmail, err = unsorted4auth.UserEmail.GetUserEmailByID(ctx, tx, userEmail.ID); err != nil {
			return err
		}

		var appUser dbo4userus.UserEntry
		if appUser, err = dal4userus.GetUserByID(ctx, tx, userEmail.Data.AppUserID); err != nil {
			return err
		}

		if userEmail.Data.ConfirmationPin() != pin {
			return errInvalidEmailConformationPin
		}

		userEmail.Data.IsConfirmed = true
		if err = userEmail.Data.SetUpdatedTime(now); err != nil {
			return fmt.Errorf("failed to set update time stamp: %w", err)
		}
		userEmail.Data.PasswordBcryptHash = []byte{}
		userEmail.Data.SetLastLoginAt(now)
		appUser.Data.SetLastLoginAt(now)

		entities := []dal.Record{appUser.Record, userEmail.Record}
		if err = tx.SetMulti(ctx, entities); err != nil {
			return err
		}
		return err
	}); err != nil {
		if dal.IsNotFound(err) {
			common4all.ErrorAsJson(ctx, w, http.StatusBadRequest, err)
			return
		} else if errors.Is(err, errInvalidEmailConformationPin) {
			common4all.ErrorAsJson(ctx, w, http.StatusForbidden, err)
			return
		}
		common4all.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
		return
	}

	ReturnToken(ctx, w, r, userEmail.Data.AppUserID, r.Referer())
}
