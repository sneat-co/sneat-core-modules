package auth

import (
	"github.com/sneat-co/sneat-core-modules/auth/const4auth"
	"github.com/sneat-co/sneat-go-core/extension"
)

func Extension() extension.Config {
	return extension.NewExtension(
		const4auth.ExtensionID,
		extension.RegisterRoutes(func(handle extension.HTTPHandleFunc) {
			// Moved to sneat-go-bots
		}),
	)
}
