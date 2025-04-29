package contactus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/api4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/delays4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/facade4spaceus"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	facade4spaceus.RegisterDboFactory(const4contactus.ModuleID, const4contactus.ContactsCollection,
		func() (dal4spaceus.SpaceModuleDbo, dal4spaceus.SpaceItemDbo, *dbo4linkage.WithRelatedAndIDs) {
			dbo := new(dbo4contactus.ContactDbo)
			return new(dbo4contactus.ContactusSpaceDbo), dbo, &dbo.WithRelatedAndIDs
		},
	)
	return module.NewModule(const4contactus.ModuleID,
		module.RegisterRoutes(api4contactus.RegisterHttpRoutes),
		module.RegisterDelays(delays4contactus.InitDelays4contactus),
	)
}
