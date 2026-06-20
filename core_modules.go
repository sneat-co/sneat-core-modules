package sneat_core_modules

import (
	"github.com/sneat-co/sneat-core-modules/auth"
	"github.com/sneat-co/sneat-core-modules/generic"
	"github.com/sneat-co/sneat-core-modules/invitus"
	"github.com/sneat-co/sneat-core-modules/linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus"
	"github.com/sneat-co/sneat-core-modules/userus"
	"github.com/sneat-co/sneat-go-core/extension"
)

// CoreExtensions returns the core module extensions.
//
// NOTE: contactus is intentionally NOT included here. It has been extracted into its own
// module (github.com/sneat-co/contactus/backend); consumers must register
// contactusext.Extension() themselves at the application composition root. Including it here
// would force core-modules to import the contactus module and re-create a module cycle.
func CoreExtensions() []extension.Config {
	return []extension.Config{
		auth.Extension(),
		generic.Extension(),
		invitus.Extension(),
		linkage.Extension(),
		spaceus.Extension(),
		userus.Extension(),
	}
}
