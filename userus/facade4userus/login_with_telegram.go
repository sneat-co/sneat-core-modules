package facade4userus

import (
	"context"
	"github.com/sneat-co/sneat-core-modules/userus/dto4userus"
)

type FirebaseCustomAuthResponse struct {
	Token  string                 `json:"token"`
	Claims map[string]interface{} `json:"claims,omitempty"`
}

func LoginWithTelegram(ctx context.Context, tgAuthData dto4userus.TelegramAuthData) (response FirebaseCustomAuthResponse, err error) {
	return
}
