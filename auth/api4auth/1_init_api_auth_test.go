package api4auth

import (
	"slices"
	"testing"

	"github.com/strongo/strongoapp"
)

func TestInitApiForAuth(t *testing.T) {
	var registered []string
	InitApiForAuth(func(method, path string, handler strongoapp.HttpHandlerWithContext) {
		registered = append(registered, method+" "+path)
	})

	missing := false
	for _, expected := range []string{
		"POST /api4debtus/auth/login-id",
	} {
		if !slices.Contains(registered, expected) {
			missing = true
			t.Errorf("expected %q to be registered", expected)
		}
	}
	if missing {
		t.Logf("registered: %+v", registered)
	}
}
