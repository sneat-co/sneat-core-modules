package userus

import (
	"github.com/sneat-co/sneat-core-modules/userus/api4userus"
	"github.com/sneat-co/sneat-core-modules/userus/const4userus"
	"github.com/sneat-co/sneat-core-modules/userus/delays4userus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func Extension() extension.Config {
	return extension.NewExtension(const4userus.ExtensionID,
		extension.RegisterRoutes(api4userus.RegisterHttpRoutes),
		extension.RegisterDelays(delays4userus.InitDelays4userus),
	)
}
