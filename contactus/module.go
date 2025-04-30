package contactus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/api4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/delays4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/facade4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	facade4linkage.RegisterDboFactory(const4contactus.ModuleID, const4contactus.ContactsCollection,
		facade4linkage.NewDboFactory(
			func() facade4linkage.SpaceItemDboWithRelatedAndIDs {
				return new(dbo4contactus.ContactDbo)
			},
			func() dal4spaceus.SpaceModuleDbo {
				return new(dbo4contactus.ContactusSpaceDbo)
			},
		),
	)
	return module.NewModule(const4contactus.ModuleID,
		module.RegisterRoutes(api4contactus.RegisterHttpRoutes),
		module.RegisterDelays(delays4contactus.InitDelays4contactus),
	)
}
