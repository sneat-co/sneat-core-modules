package calendarium

import (
	"github.com/sneat-co/sneat-core-modules/calendarium/api4calendarium"
	"github.com/sneat-co/sneat-core-modules/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-core/modules"
)

func Module() modules.Module {
	return modules.NewModule(const4calendarium.ModuleID, api4calendarium.RegisterHttpRoutes)
}
