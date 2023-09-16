package userus

import (
	"github.com/sneat-co/sneat-core-modules/userus/api4userus"
	"github.com/sneat-co/sneat-core-modules/userus/const4userus"
	"github.com/sneat-co/sneat-go-core/modules"
)

func Module() modules.Module {
	return modules.NewModule(const4userus.ModuleID, api4userus.RegisterHttpRoutes)
}
