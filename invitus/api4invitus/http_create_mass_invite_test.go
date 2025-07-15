package api4invitus

import (
	"bytes"
	"encoding/json"
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-core-modules/invitus/facade4invitus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/sneatauth"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateMassInvite(t *testing.T) {
	const spaceID = "unit-test"
	var invite dbo4invitus.InviteDbo
	invite.Type = dbo4invitus.InviteTypeMass
	invite.Channel = "email"
	invite.Roles = []string{
		"contributor",
		"test-role1",
	}
	invite.From = dbo4invitus.InviteFrom{
		InviteContact: dbo4invitus.InviteContact{
			Channel:   "email",
			Address:   "from@example.com",
			Title:     "From Title",
			ContactID: "f1",
			UserID:    "u1",
		},
	}
	//invite.To = &dbo4invitus.InviteTo{
	//	Channel:      "email",
	//	Address:      "to@example.com",
	//	Title:        "To Title",
	//	ToSpaceContactID: "t1",
	//}
	invite.FromUserID = "u1"
	invite.SpaceID = spaceID
	invite.Space = &dbo4invitus.InviteSpace{
		Type:  coretypes.SpaceTypeFamily,
		Title: "Unit Test",
	}
	invite.Created.Client.HostOrApp = "unit-test"
	invite.Created.Client.RemoteAddr = "127.0.0.1"
	invite.CreatedAt = time.Now()
	invite.From.UserID = "u1"
	invite.Status = dbo4invitus.InviteStatusPending
	invite.Pin = "123456"

	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(facade4invitus.CreateMassInviteRequest{Invite: invite}); err != nil {
		t.Fatal(err)
	}
	//t.Log(buffer.String())

	req, err := http.NewRequest(http.MethodPost, "/api4meetingus/create-invite", buffer)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Origin", "http://localhost:3000")

	createMassInviteFunc = func(ctx facade.ContextWithUser, request facade4invitus.CreateMassInviteRequest) (response facade4invitus.CreateInviteResponse, err error) {
		response.Invite.ID = "test-id"
		return
	}

	apicore.GetAuthTokenFromHttpRequest = func(r *http.Request, authRequired bool) (token *sneatauth.Token, err error) {
		return &sneatauth.Token{UID: "unit-test-user"}, nil
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(httpPostCreateMassInvite)
	handler.ServeHTTP(rr, req)

	responseBody := rr.Body.String()

	if expected := http.StatusCreated; rr.Code != expected {
		t.Fatalf(
			"unexpected status: got (%d) expects (%d): %s",
			rr.Code,
			expected,
			responseBody,
		)
	}

	var response facade4invitus.CreateInviteResponse
	if err = json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err, responseBody)
	}
	if response.Invite.ID == "" {
		t.Fatal("Response is missing invite ID")
	}
}
