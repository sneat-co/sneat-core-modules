package dbo4invitus

import "testing"

func TestInviteDbo_Validate(t *testing.T) {
	invite := InviteDbo{}
	if err := invite.Validate(); err == nil {
		t.Fatal("error expected for empty value")
	}
}
