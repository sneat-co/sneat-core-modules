package invitus

import (
	"github.com/sneat-co/sneat-core-modules/invitus/const4invitus"
	"github.com/sneat-co/sneat-go-core/extension"
	"testing"
)

func TestModule(t *testing.T) {
	m := Extension()
	extension.AssertExtension(t, m, extension.Expected{
		ExtID:         const4invitus.ExtensionID,
		HandlersCount: 6,
		DelayersCount: 0,
	})
}
