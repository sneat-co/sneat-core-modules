package contactus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/api4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/delays4contactus"
	"github.com/sneat-co/sneat-core-modules/contactusmodels/const4contactus"
	"github.com/sneat-co/sneat-core-modules/invitus/facade4invitus"
	"github.com/sneat-co/sneat-core-modules/linkage/facade4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/facade4spaceus"
	"github.com/sneat-co/sneat-core-modules/userus/facade4userus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func Extension() extension.Config {
	facade4spaceus.RegisterContactusSpaceContributor(spaceusContactusContributor{})
	facade4userus.RegisterContactusCountryUpdater(userusContactusContributor{})
	facade4invitus.RegisterContactusAccess(invitusContactusAccess{})
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
