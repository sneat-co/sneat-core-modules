package cmds4sneatbot

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-core-modules/botscore/bothelpers"
	tghelpers2 "github.com/sneat-co/sneat-core-modules/botscore/tghelpers"
	"github.com/strongo/logus"
	"net/url"
)

var membersCommand = botsfw.Command{
	Code:     "members",
	Commands: []string{"/members"},
	InputTypes: []botinput.WebhookInputType{
		botinput.WebhookInputText,
		botinput.WebhookInputCallbackQuery,
	},
	CallbackAction: membersCallbackAction,
	Action:         membersAction,
}

func membersCallbackAction(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	if m, err = membersAction(whc); err != nil {
		return
	}

	keyboard := m.Keyboard.(*tgbotapi.InlineKeyboardMarkup)
	spaceRef := tghelpers2.GetSpaceRef(callbackUrl)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tgbotapi.InlineKeyboardButton{
		tghelpers2.BackToSpaceMenuButton(spaceRef),
	})
	if m, err = whc.NewEditMessage(m.Text, m.Format); err != nil {
		return
	}
	m.Keyboard = keyboard

	m.EditMessageUID, err = tghelpers2.GetEditMessageUID(whc)
	return
}

func membersAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()
	logus.Infof(ctx, "membersCommand.Action(): InputType=%v", whc.Input().InputType())
	m.Text = "<b>Family members</b>"
	m.Format = botsfw.MessageFormatHTML
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text: "💻 Manage in app",
				WebApp: &tgbotapi.WebappInfo{
					Url: bothelpers.GetBotWebAppUrl() + "space/family/h4qax/members", // TODO: generate URL
				},
			},
			{
				Text:         "➕ Add member",
				CallbackData: "/add-member",
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "🧑 Myself",
				CallbackData: "/contact=myself",
			},
		},
	)
	return
}
