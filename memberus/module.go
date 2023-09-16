package memberus

import (
	"github.com/sneat-co/sneat-core-modules/memberus/api4memberus"
	"github.com/sneat-co/sneat-core-modules/memberus/const4memberus"
	"github.com/sneat-co/sneat-go-core/modules"
)

func Module() modules.Module {
	return modules.NewModule(const4memberus.ModuleID, api4memberus.RegisterHttpRoutes)
}
