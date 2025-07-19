package spaceus

import (
	"github.com/sneat-co/sneat-core-modules/spaceus/const4spaceus"
	"github.com/sneat-co/sneat-go-core/extension"
	"testing"
)

func TestModule(t *testing.T) {
	m := Extension()
	extension.AssertExtension(t, m, extension.Expected{
		ExtID:         const4spaceus.ExtensionID,
		HandlersCount: 8,
		DelayersCount: 0,
	})
}
