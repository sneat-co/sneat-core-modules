package userus

import (
	"testing"

	"github.com/sneat-co/sneat-core-modules/userus/const4userus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func TestModule(t *testing.T) {
	m := Extension()
	extension.AssertExtension(t, m, extension.Expected{
		ExtID:         const4userus.ExtensionID,
		HandlersCount: 4,
		DelayersCount: 1,
	})
}
