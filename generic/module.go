package generic

import (
	"github.com/sneat-co/sneat-core-modules/generic/api4generic"
	"github.com/sneat-co/sneat-core-modules/generic/const4generic"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	return module.NewExtension(const4generic.ModuleID, module.RegisterRoutes(api4generic.RegisterHttpRoutes))
}
