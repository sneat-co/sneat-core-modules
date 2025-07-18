package facade4userus

import (
	"testing"
	"time"
)

func TestSetUserTimezone(t *testing.T) {
	t.Run("ctx_nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("panic expected but succeed")
			}
		}()
		tz, _ := time.LoadLocation("Europe/London")
		_, _ = SetUserTimezone(nil, tz)
	})
}
