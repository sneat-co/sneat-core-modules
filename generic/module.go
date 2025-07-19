package generic

import (
	"github.com/sneat-co/sneat-core-modules/generic/api4generic"
	"github.com/sneat-co/sneat-core-modules/generic/const4generic"
	"github.com/sneat-co/sneat-go-core/extension"
)

func Extension() extension.Config {
	return extension.NewExtension(const4generic.ExtensionID, extension.RegisterRoutes(api4generic.RegisterHttpRoutes))
}
