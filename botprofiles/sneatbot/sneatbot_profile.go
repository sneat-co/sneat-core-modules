package sneatbot

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-core-modules/anybot"
	"github.com/sneat-co/sneat-core-modules/anybot/cmds4anybot"
	cmds4sneatbot2 "github.com/sneat-co/sneat-core-modules/botprofiles/sneatbot/cmds4sneatbot"
	//"github.com/sneat-co/sneat-go-backend/src/modules/listus/listusbot/cmds4listusbot"
)

const ProfileID = "sneat_bot"

var sneatBotProfile botsfw.BotProfile

func GetProfile(errFooterText func() string) botsfw.BotProfile {
	if sneatBotProfile == nil {
		sneatBotProfile = createSneatBotProfile(errFooterText)
	}
	return sneatBotProfile
}

func createSneatBotProfile(errFooterText func() string) botsfw.BotProfile {
	router := botsfw.NewWebhookRouter(errFooterText)

	botParams := cmds4sneatbot2.GetBotParams()
	cmds4anybot.AddSharedCommands(router, botParams)

	commandsByType := make(map[botinput.WebhookInputType][]botsfw.Command) // TODO: get rid of `commandsByType`

	cmds4sneatbot2.AddSneatBotOnlyCommands(commandsByType)
	cmds4sneatbot2.AddSneatSharedCommands(commandsByType)
	//cmds4listusbot.AddListusSharedCommands(commandsByType)

	router.AddCommandsGroupedByType(commandsByType)

	return anybot.NewProfile(ProfileID, &router)
}
