package linkage

import (
	"github.com/sneat-co/sneat-core-modules/linkage/api4linkage"
	"github.com/sneat-co/sneat-core-modules/linkage/const4linkage"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	return module.NewExtension(const4linkage.ModuleID, module.RegisterRoutes(api4linkage.RegisterHttpRoutes))
}
