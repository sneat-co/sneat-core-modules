package spaceus

import (
	"github.com/sneat-co/sneat-core-modules/spaceus/api4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/const4spaceus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func Extension() extension.Config {
	return extension.NewExtension(const4spaceus.ExtensionID,
		extension.RegisterRoutes(api4spaceus.RegisterHttpRoutes),
	)
}
