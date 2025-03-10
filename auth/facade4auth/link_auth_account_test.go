package facade4auth

import (
	"testing"
)

func TestLinkAuthAccount(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("LinkAuthAccount() did not panic")
		}
	}()
	LinkAuthAccount(nil)
}
