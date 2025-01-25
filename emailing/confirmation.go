package emailing

import (
	"context"
	"fmt"
	"github.com/sneat-co/sneat-core-modules/auth/models4auth"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
)

func CreateConfirmationEmailAndQueueForSending(ctx context.Context, user dbo4userus.UserEntry, userEmail models4auth.UserEmailEntry) error {
	emailEntity := &models4auth.EmailData{
		From:    "Alex @ DebtsTracker.io <alex@debtusbot.io>",
		To:      userEmail.ID,
		Subject: "Please confirm your account at Sneat.app",
		BodyText: fmt.Sprintf(`%v, we are thrilled to have you on board!

To keep your account secure please confirm your email by clicking this link:

  >> https://debtstracker.io/confirm?email=%v&pin=%v

If you have any questions or issue please drop me an email to alex@debtusbot.io
--
Alex
Creator of https://DebtsTracker.io

We are social:
  FB page - https://www.facebook.com/debtstracker
  Twitter - https://twitter.com/debtstracker
`, user.Data.GetFullName(), userEmail.ID, userEmail.Data.ConfirmationPin()),
	}
	_, err := CreateEmailRecordAndQueueForSending(ctx, emailEntity)
	return err
}
