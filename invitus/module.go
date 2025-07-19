package invitus

import (
	"github.com/sneat-co/sneat-core-modules/invitus/api4invitus"
	"github.com/sneat-co/sneat-core-modules/invitus/const4invitus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func Extension() extension.Config {
	return extension.NewExtension(const4invitus.ExtensionID, extension.RegisterRoutes(api4invitus.RegisterHttpRoutes))
}
