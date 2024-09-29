package cmds4anybot

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-core-modules/botscore/bothelpers"
	tghelpers2 "github.com/sneat-co/sneat-core-modules/botscore/tghelpers"
	"net/url"
)

const SpaceSettingsCommandCode = "spaceSettings"

var spaceSettingsCommand = botsfw.Command{
	Code:     SpaceSettingsCommandCode,
	Commands: []string{"/spaceSettings"},
	InputTypes: []botinput.WebhookInputType{
		botinput.WebhookInputText,
		botinput.WebhookInputCallbackQuery,
	},
	CallbackAction: settingsCallbackAction,
	Action:         settingsAction,
}

func settingsCallbackAction(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	if m, err = settingsAction(whc); err != nil {
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

func settingsAction(_ botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	m.Format = botsfw.MessageFormatHTML
	m.Text = "<b>Family space settings</b>"
	m.Text += "\n\n<i>Not implemented yet</i>"
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text: "💻 Manage in app",
				WebApp: &tgbotapi.WebappInfo{
					Url: bothelpers.GetBotWebAppUrl() + "space/family/h4qax/settings", // TODO: generate URL
				},
			},
		},
	)
	return
}
