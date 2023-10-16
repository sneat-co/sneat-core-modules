package schedulus

import (
	"github.com/sneat-co/sneat-core-modules/schedulus/api4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/const4schedulus"
	"github.com/sneat-co/sneat-go-core/modules"
)

func Module() modules.Module {
	return modules.NewModule(const4schedulus.ModuleID, api4schedulus.RegisterHttpRoutes)
}
