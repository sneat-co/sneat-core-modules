package generic

import (
	"testing"

	"github.com/sneat-co/sneat-core-modules/generic/const4generic"
	"github.com/sneat-co/sneat-go-core/extension"
)

func TestModule(t *testing.T) {
	m := Extension()
	extension.AssertExtension(t, m, extension.Expected{
		ExtID:         const4generic.ExtensionID,
		HandlersCount: 3,
		DelayersCount: 0,
	})
}
