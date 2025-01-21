package dal4bots

import (
	"context"
	"github.com/bots-go-framework/bots-fw/botsdal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/bots/botscore/models4bots"
)

func GetBotUserRecord(_ context.Context, botPlatformID, botID, botUserID string) (
	tgBotUser record.DataWithID[string, *models4bots.TelegramUserDbo], err error,
) {
	_ = botID // is not used
	botUserKey := botsdal.NewPlatformUserKey(botPlatformID, botUserID)
	tgBotUser = record.NewDataWithID(botUserID, botUserKey, new(models4bots.TelegramUserDbo))
	return
}
