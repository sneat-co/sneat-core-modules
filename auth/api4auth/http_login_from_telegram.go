package api4auth

import (
	"context"
	"fmt"
	telegram "github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw-telegram-webapp/twainitdata"
	facade4auth2 "github.com/sneat-co/sneat-core-modules/auth/facade4auth"
	"github.com/sneat-co/sneat-go-core/apicore"
	"net/http"
)

func signInWithTelegram(ctx context.Context, w http.ResponseWriter, r *http.Request, botID string, initData twainitdata.InitData) {
	var (
		err       error
		tgBotUser facade4auth2.BotUserEntry
	)
	remoteClientInfo := apicore.GetRemoteClientInfo(r)
	if tgBotUser, _, err = facade4auth2.SignInWithTelegram(ctx, botID, initData, remoteClientInfo); err != nil {
		apicore.ReturnError(ctx, w, r, fmt.Errorf("failed to sign in with Telegram: %w", err))
		return
	}

	appUserID := tgBotUser.Data.GetAppUserID()
	ReturnToken(ctx, w, r, appUserID, telegram.PlatformID)
}
