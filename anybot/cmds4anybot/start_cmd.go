package cmds4anybot

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	facade4anybot2 "github.com/sneat-co/sneat-core-modules/anybot/facade4anybot"
	models4auth2 "github.com/sneat-co/sneat-core-modules/auth/models4auth"
	"github.com/sneat-co/sneat-core-modules/common4all"
	"github.com/sneat-co/sneat-core-modules/tgsharedcommands"
	"github.com/sneat-co/sneat-core-modules/userus/dal4userus"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/strongo/logus"
	"net/url"
	"strings"
)

var ErrUnknownStartParam = errors.New("unknown start parameter")

func StartBotLink(botID, command string, params ...string) string {
	var buf bytes.Buffer
	_, _ = fmt.Fprintf(&buf, "https://t.me/%v?start=%v", botID, command)
	for _, p := range params {
		buf.WriteString("__")
		buf.WriteString(p)
	}
	return buf.String()
}

const StartCommandCode = "start"

func createStartCommand(
	startInBotAction StartInBotActionFunc,
	startInGroupAction botsfw.CommandAction,
	getWelcomeMessageText WelcomeMessageProvider,
	mainMenuAction SetMainMenuFunc,
) botsfw.Command {
	return botsfw.Command{
		Code:     StartCommandCode,
		Commands: []string{"/start"},
		InputTypes: []botinput.WebhookInputType{
			botinput.WebhookInputText,
			botinput.WebhookInputReferral,            // FBM
			botinput.WebhookInputConversationStarted, // Viber
		},
		Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
			return sharedStartCommandAction(whc /*startInBotAction,*/, startInGroupAction, mainMenuAction)
		},
		CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
			return sharedStartCommandCallbackAction(whc, callbackUrl, getWelcomeMessageText, startInBotAction, mainMenuAction)
		},
	}
}

func sharedStartCommandCallbackAction(
	whc botsfw.WebhookContext,
	callbackUrl *url.URL,
	getWelcomeMessageText WelcomeMessageProvider,
	startInBotAction StartInBotActionFunc,
	mainMenuAction SetMainMenuFunc,
) (
	m botsfw.MessageFromBot, err error,
) {
	q := callbackUrl.Query()
	if localeCode := q.Get("locale"); localeCode != "" {
		if m, err = setPreferredLocaleAction(whc, localeCode, setPreferredLocaleModeStart, mainMenuAction); err != nil {
			return m, fmt.Errorf("failed to setPreferredLocaleAction(): %w", err)
		}
		m.IsEdit = true
	} else {
		m.Text = fmt.Sprintf("Unknown callback parameters: %s", callbackUrl)
		return
	}

	// Sends a message over HTTPS as we should ensure that we do not send it multiple times
	if err = runBotSpecificStartCommand(whc, startInBotAction, nil, getWelcomeMessageText); err != nil {
		return
	}

	var welcomeText string
	if welcomeText, err = getWelcomeMessageText(whc); err != nil {
		return m, fmt.Errorf("failed to get welcome message text: %w", err)
	}
	m.Text = welcomeText + "\n" + strings.Repeat("-", len(m.Text)) + "\n" + m.Text
	m.Format = botsfw.MessageFormatHTML
	return
}

func sharedStartCommandAction(
	whc botsfw.WebhookContext,
	//startInBotAction StartInBotActionFunc,
	startInGroupAction botsfw.CommandAction,
	mainMenuAction SetMainMenuFunc,
) (
	m botsfw.MessageFromBot, err error,
) {
	whc.Input().LogRequest()
	ctx := whc.Context()
	text := whc.Input().(botinput.WebhookTextMessage).Text()
	logus.Debugf(ctx, "createStartCommand.Action() => text: "+text)

	startParam, _ := tgsharedcommands.ParseStartCommand(whc)

	var isInGroup bool
	if isInGroup, err = whc.IsInGroup(); err != nil {
		return
	} else if isInGroup {
		return startInGroupAction(whc)
	}
	chatEntity := whc.ChatData()
	chatEntity.SetAwaitingReplyTo("")

	switch {
	case startParam == "help_inline":
		return startInlineHelp(whc)
	case strings.HasPrefix(startParam, "login-"):
		loginID, err := common4all.DecodeIntID(startParam[len("login-"):])
		if err != nil {
			return m, err
		}
		return startLoginGac(whc, loginID)
		//case strings.HasPrefix(textToMatchNoStart, JOIN_BILL_COMMAND):
		//	return JoinBillCommand.Action(whc)
	case strings.HasPrefix(startParam, "refbytguser-") && startParam != "refbytguser-YOUR_CHANNEL":
		facade4anybot2.Referer.AddTelegramReferrer(ctx, whc.AppUserID(), strings.TrimPrefix(startParam, "refbytguser-"), whc.GetBotCode())
	}
	//if m.Text, err = getWelcomeMessage(whc); err != nil {
	//	return
	//} else if m.Text != "" {
	//	responder := whc.Responder()
	//	if _, err = responder.SendMessage(ctx, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
	//		return
	//	}
	//}
	/*
		var user dbo4userus.UserEntry
		if user, err = GetUser(whc); err != nil {
			return
		}
		if user.Data.PreferredLocale == ""
	*/
	{
		var localesMsg botsfw.MessageFromBot
		if localesMsg, err = onStartAskLocaleAction(whc, mainMenuAction); err != nil {
			return
		}
		if localesMsg.Text = strings.TrimSpace(localesMsg.Text); localesMsg.Text != "" {
			m.Text += "\n" + localesMsg.Text
			m.Keyboard = localesMsg.Keyboard
			m.Format = botsfw.MessageFormatHTML
		}
		return
	}
	//if m, err = runBotSpecificStartCommand(whc, startInBotAction, startParams, getWelcomeMessage); err != nil {
	//	return
	//}
	//return
}

func runBotSpecificStartCommand(whc botsfw.WebhookContext, startInBotAction StartInBotActionFunc, startParams []string, getWelcomeMessage WelcomeMessageProvider) (err error) {
	var m botsfw.MessageFromBot
	if m, err = startInBotAction(whc, startParams); err != nil {
		return
	}
	responder := whc.Responder()
	ctx := whc.Context()
	if _, err = responder.SendMessage(ctx, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
		return
	}

	//if m.Text, err = getWelcomeMessage(whc); err != nil {
	//	return
	//}
	//m.Format = botsfw.MessageFormatHTML
	//m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
	//	[]tgbotapi.InlineKeyboardButton{
	//		{
	//			Text:         "🛠 Settings",
	//			CallbackData: SettingsCommandCode,
	//		},
	//		{
	//			Text:         whc.CommandText(trans.COMMAND_TEXT_LANGUAGE, emoji.EARTH_ICON),
	//			CallbackData: UserSettingsCommandCode,
	//		},
	//	},
	//)

	return
}

func startLoginGac(whc botsfw.WebhookContext, loginID int) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()
	var loginPin models4auth2.LoginPin
	if loginPin, err = facade4anybot2.AuthFacade.AssignPinCode(ctx, loginID, whc.AppUserID()); err != nil {
		return
	}
	return whc.NewMessageByCode(trans.MESSAGE_TEXT_LOGIN_CODE, models4auth2.LoginCodeToString(loginPin.Data.Code)), nil
}

func startInlineHelp(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	m = whc.NewMessage("<b>Help: How to use this bot in chats</b>\n\nExplain here how to use bot's inline mode.")
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{{Text: "Button 1", URL: "https://debtstracker.io/#btn=1"}},
		[]tgbotapi.InlineKeyboardButton{{Text: "Button 2", URL: "https://debtstracker.io/#btn=2"}},
		//[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonSwitch("Back to chat 1", "1")},
		//[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonSwitch("Back to chat 2", "2")},
		[]tgbotapi.InlineKeyboardButton{{Text: "Button 3", CallbackData: "help-3"}},
		[]tgbotapi.InlineKeyboardButton{{Text: "Button 4", CallbackData: "help-4"}},
		[]tgbotapi.InlineKeyboardButton{{Text: "Button 5", CallbackData: "help-5"}},
	)
	return m, err
}

func GetUser(whc botsfw.WebhookContext) (user dbo4userus.UserEntry, err error) { // TODO: Make library and use across app
	appUserID := whc.AppUserID()
	if appUserID == "" {
		return user, fmt.Errorf("%w: app user ID is empty", dal.ErrRecordNotFound)
	}
	user = dbo4userus.NewUserEntry(appUserID)
	ctx := whc.Context()
	tx := whc.Tx()
	return user, dal4userus.GetUser(ctx, tx, user)
}

//var LangKeyboard = tgbotapi.NewInlineKeyboardMarkup(
//	[]tgbotapi.InlineKeyboardButton{
//		{
//			Text:         i18n.LocaleEnUS.TitleWithIcon(),
//			CallbackData: onStartCallbackCommandCode + "?lang=" + i18n.LocaleCodeEnUS,
//		},
//		{
//			Text:         i18n.LocaleRuRu.TitleWithIcon(),
//			CallbackData: onStartCallbackCommandCode + "?lang=" + i18n.LocalCodeRuRu,
//		},
//	},
//)
