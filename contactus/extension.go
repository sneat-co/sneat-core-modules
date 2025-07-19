package contactus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/api4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/delays4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/facade4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func Extension() extension.Config {
	facade4linkage.RegisterDboFactory(const4contactus.ExtensionID, const4contactus.ContactsCollection,
		facade4linkage.NewDboFactory(
			func() facade4linkage.SpaceItemDboWithRelatedAndIDs {
				return new(dbo4contactus.ContactDbo)
			},
			func() dal4spaceus.SpaceModuleDbo {
				return new(dbo4contactus.ContactusSpaceDbo)
			},
		),
	)
	return extension.NewExtension(const4contactus.ExtensionID,
		extension.RegisterRoutes(api4contactus.RegisterHttpRoutes),
		extension.RegisterDelays(delays4contactus.InitDelays4contactus),
	)
}
