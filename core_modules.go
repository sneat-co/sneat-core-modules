package sneat_core_modules

import (
	"github.com/sneat-co/sneat-core-modules/auth"
	"github.com/sneat-co/sneat-core-modules/contactus"
	"github.com/sneat-co/sneat-core-modules/generic"
	"github.com/sneat-co/sneat-core-modules/invitus"
	"github.com/sneat-co/sneat-core-modules/linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus"
	"github.com/sneat-co/sneat-core-modules/userus"
	"github.com/sneat-co/sneat-go-core/module"
)

func CoreModules() []module.Module {
	return []module.Module{
		auth.Module(),
		contactus.Module(),
		generic.Module(),
		invitus.Module(),
		linkage.Module(),
		spaceus.Module(),
		userus.Module(),
	}
}
