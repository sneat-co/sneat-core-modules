package facade4userus

import (
	"testing"
)

func TestSetUserTimezone(t *testing.T) {
	t.Run("ctx_nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("panic expected but succeed")
			}
		}()
		_, _ = SetUserTimezone(nil, "Europe/London")
	})
}
