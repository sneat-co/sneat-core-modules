package auth

import (
	"github.com/sneat-co/sneat-core-modules/auth/const4auth"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	return module.NewExtension(
		const4auth.ExtensionID,
		module.RegisterRoutes(func(handle module.HTTPHandleFunc) {
			// Moved to sneat-go-bots
		}),
	)
}
