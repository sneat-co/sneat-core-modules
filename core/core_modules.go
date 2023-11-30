package core

import "github.com/sneat-co/sneat-go-core/modules"

import (
	"github.com/sneat-co/sneat-core-modules/calendarium"
	"github.com/sneat-co/sneat-core-modules/contactus"
	"github.com/sneat-co/sneat-core-modules/invitus"
	"github.com/sneat-co/sneat-core-modules/teamus"
	"github.com/sneat-co/sneat-core-modules/userus"
)

func Modules() []modules.Module {
	return []modules.Module{
		calendarium.Module(),
		contactus.Module(),
		invitus.Module(),
		teamus.Module(),
		userus.Module(),
	}
}
