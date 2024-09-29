package api4auth

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/auth/models4auth"
	"github.com/sneat-co/sneat-core-modules/auth/token4auth"
	"github.com/sneat-co/sneat-core-modules/auth/unsorted4auth"
	common4all2 "github.com/sneat-co/sneat-core-modules/common4all"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"io"
	"net/http"
	"strings"
)

func HandleAuthLoginId(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
	query := r.URL.Query()
	channel := query.Get("channel")
	var (
		loginID int
		err     error
	)

	loginIdStr := query.Get("id")

	if loginIdStr != "" {
		if loginID, err = common4all2.DecodeIntID(loginIdStr); err != nil {
			common4all2.BadRequestError(ctx, w, err)
			return
		}
	}

	returnLoginID := func(loginID int) {
		encoded := common4all2.EncodeIntID(loginID)
		logus.Infof(ctx, "Login ContactID: %d, Encoded: %s", loginID, encoded)
		if _, err = w.Write([]byte(encoded)); err != nil {
			logus.Criticalf(ctx, "Failed to write login ContactID to response: %v", err)
		}
	}

	if loginID != 0 {
		if loginPin, err := unsorted4auth.LoginPin.GetLoginPinByID(ctx, nil, loginID); err != nil {
			if dal.IsNotFound(err) {
				w.WriteHeader(http.StatusInternalServerError)
				logus.Errorf(ctx, err.Error())
				return
			}
		} else if loginPin.Data.IsActive(channel) {
			returnLoginID(loginID)
			return
		}
	}

	var rBody []byte
	if rBody, err = io.ReadAll(r.Body); err != nil {
		common4all2.BadRequestError(ctx, w, fmt.Errorf("failed to read request body: %w", err))
		return
	}
	gaClientID := string(rBody)

	if gaClientID != "" {
		if len(gaClientID) > 100 {
			common4all2.BadRequestMessage(ctx, w, fmt.Sprintf("Google Client ContactID is too long: %d", len(gaClientID)))
			return
		}

		if strings.Count(gaClientID, ".") != 1 {
			common4all2.BadRequestMessage(ctx, w, "Google Client ContactID has wrong format, a '.' char expected")
			return
		}
	}

	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		var loginPin models4auth.LoginPin
		if loginPin, err = unsorted4auth.LoginPin.CreateLoginPin(ctx, tx, channel, gaClientID, authInfo.UserID); err != nil {
			common4all2.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
			return
		}
		loginID = loginPin.ID
		return err
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logus.Errorf(ctx, err.Error())
		return
	}
	returnLoginID(loginID)
}
