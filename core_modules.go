package sneat_core_modules

import (
	"github.com/sneat-co/sneat-core-modules/auth"
	"github.com/sneat-co/sneat-core-modules/contactus"
	"github.com/sneat-co/sneat-core-modules/generic"
	"github.com/sneat-co/sneat-core-modules/invitus"
	"github.com/sneat-co/sneat-core-modules/linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus"
	"github.com/sneat-co/sneat-core-modules/userus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func CoreExtensions() []extension.Config {
	return []extension.Config{
		auth.Extension(),
		contactus.Extension(),
		generic.Extension(),
		invitus.Extension(),
		linkage.Extension(),
		spaceus.Extension(),
		userus.Extension(),
	}
}
