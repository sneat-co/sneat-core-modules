package facade4invitus

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"mime"

	core "github.com/sneat-co/sneat-go-core"
	"github.com/sneat-co/sneat-go-core/capturer"
	"github.com/sneat-co/sneat-go-core/emails"
)

const inviteEmailTemplateText = `
<p>
	InviteDbo from: <i>{{.fromHTML}}</i>
</p>

<p>
	To join <b>{{.space.Title}}</b> please follow the link:
</p>

<p>
	https://{{.hostPath}}/join/{{.space.Type}}?id={{.id}}#pin={{.pinCode}}
</p>

<p>You personal PIN code to join the space is: <b><code>{{.pinCode}}</code></b></p>

<p>https://sneat.app/ - A family app that saves you time & money.</p>

<p>P.S. If any issues feel free to get <a href="mailto:help@sneat.app">help@sneat.app</a></p>
`

var inviteEmailTemplate = template.Must(template.New("inviteEmail").Parse(inviteEmailTemplateText))

func sendInviteEmail(ctx context.Context, invite InviteEntry) (messageID string, err error) {
	if invite.Data.To.Address == "" {
		return "", errors.New("missing required field: invite.To.Address")
	}
	templateData := make(map[string]any)
	if core.IsInProd() {
		templateData["hostPath"] = "sneat.app/pwa"
	} else {
		templateData["hostPath"] = "localhost:4200"
	}
	templateData["id"] = invite.ID
	if invite.Data.From.Address == "" {
		templateData["fromHTML"] = invite.Data.From.Title
	} else {
		templateData["fromHTML"] = fmt.Sprintf(`<a href="mailto:%s">%s</a>`, invite.Data.From.Address, invite.Data.From.Title)
	}
	templateData["invite"] = invite
	templateData["space"] = invite.Data.Space
	templateData["pinCode"] = invite.Data.Pin
	buf := new(bytes.Buffer)
	if err := inviteEmailTemplate.Execute(buf, templateData); err != nil {
		return "", fmt.Errorf("failed to create email message body: %w", err)
	}

	msg := emails.Email{
		From:    fmt.Sprintf(`"%s" <inviter@sneat.app>`, mime.QEncoding.Encode("utf-8", invite.Data.From.Title)),
		To:      []string{invite.Data.To.Address},
		Subject: fmt.Sprintf("You are invited by %s to join %s", invite.Data.From.Title, invite.Data.Space.Title),
		HTML:    buf.String(),
		//ReplyTo: nil,
	}
	var response emails.Sent
	if response, err = emails.Send(ctx, msg); err != nil {
		err = capturer.CaptureError(ctx, err)
		return
	}
	if response != nil {
		messageID = response.MessageID()
	}
	return
}
