package generic

import (
	"github.com/sneat-co/sneat-core-modules/generic/const4generic"
	"github.com/sneat-co/sneat-go-core/module"
	"testing"
)

func TestModule(t *testing.T) {
	m := Module()
	module.AssertModule(t, m, module.Expected{
		ModuleID:      const4generic.ModuleID,
		HandlersCount: 3,
		DelayersCount: 0,
	})
}
