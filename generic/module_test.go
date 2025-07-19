package generic

import (
	"github.com/sneat-co/sneat-core-modules/generic/const4generic"
	"github.com/sneat-co/sneat-go-core/extension"
	"testing"
)

func TestModule(t *testing.T) {
	m := Extension()
	extension.AssertExtension(t, m, extension.Expected{
		ExtID:         const4generic.ExtensionID,
		HandlersCount: 3,
		DelayersCount: 0,
	})
}
