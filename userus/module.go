package userus

import (
	"github.com/sneat-co/sneat-core-modules/userus/api4userus"
	"github.com/sneat-co/sneat-core-modules/userus/const4userus"
	"github.com/sneat-co/sneat-core-modules/userus/delays4userus"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	return module.NewExtension(const4userus.ModuleID,
		module.RegisterRoutes(api4userus.RegisterHttpRoutes),
		module.RegisterDelays(delays4userus.InitDelays4userus),
	)
}
