package linkage

import (
	"github.com/sneat-co/sneat-core-modules/linkage/api4linkage"
	"github.com/sneat-co/sneat-core-modules/linkage/const4linkage"
	"github.com/sneat-co/sneat-go-core/extension"
)

func Extension() extension.Config {
	return extension.NewExtension(const4linkage.ExtensionID,
		extension.RegisterRoutes(api4linkage.RegisterHttpRoutes),
	)
}
