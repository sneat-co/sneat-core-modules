package cmds4sneatbot

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	cmds4anybot2 "github.com/sneat-co/sneat-core-modules/anybot/cmds4anybot"
)

func GetBotParams() cmds4anybot2.BotParams {
	return cmds4anybot2.BotParams{
		GetWelcomeMessageText: sneatBotWelcomeMessage,
		StartInBotAction:      startActionWithStartParams,
		StartInGroupAction: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
			m.Text = "Start in group is not implemented yet for @SneatBot"
			return
		},
		HelpCommandAction: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
			m.Text = "Help is not implemented yet for @SneatBot"
			return
		},
		SetMainMenu: func(whc botsfw.WebhookContext, messageText string, showHint bool) (m botsfw.MessageFromBot, err error) {
			m.Text = "SneatBot main menu"
			return
		},
	}
}
