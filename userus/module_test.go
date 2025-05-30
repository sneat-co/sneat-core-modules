package userus

import (
	"github.com/sneat-co/sneat-core-modules/userus/const4userus"
	"github.com/sneat-co/sneat-go-core/module"
	"testing"
)

func TestModule(t *testing.T) {
	m := Module()
	module.AssertModule(t, m, module.Expected{
		ModuleID:      const4userus.ModuleID,
		HandlersCount: 4,
		DelayersCount: 1,
	})
}
