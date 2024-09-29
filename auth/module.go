package auth

import (
	"github.com/sneat-co/sneat-core-modules/auth/api4auth"
	"github.com/sneat-co/sneat-core-modules/auth/const4auth"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	return module.NewModule(
		const4auth.ModuleID,
		module.RegisterRoutes(api4auth.RegisterHttpRoutes),
	)
}
