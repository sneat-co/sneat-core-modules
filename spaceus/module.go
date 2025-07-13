package spaceus

import (
	"github.com/sneat-co/sneat-core-modules/spaceus/api4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/const4spaceus"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	return module.NewExtension(const4spaceus.ModuleID, module.RegisterRoutes(api4spaceus.RegisterHttpRoutes))
}
