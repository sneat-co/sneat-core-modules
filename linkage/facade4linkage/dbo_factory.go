package facade4linkage

import (
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
)

type SpaceItemDboWithRelatedAndIDs interface {
	dal4spaceus.SpaceItemDbo
	RelatedAndIDs() *dbo4linkage.WithRelatedAndIDs
}

type RelatedDboFactory = interface {
	NewItemDbo() SpaceItemDboWithRelatedAndIDs
	NewSpaceModuleDbo() dal4spaceus.SpaceModuleDbo
}

func NewDboFactory(
	newItemDbo func() SpaceItemDboWithRelatedAndIDs,
	newSpaceModuleDbo func() dal4spaceus.SpaceModuleDbo,
) RelatedDboFactory {
	return relatedDboFactory{
		newItemDbo:        newItemDbo,
		newSpaceModuleDbo: newSpaceModuleDbo,
	}
}

type relatedDboFactory struct {
	newItemDbo        func() SpaceItemDboWithRelatedAndIDs
	newSpaceModuleDbo func() dal4spaceus.SpaceModuleDbo
}

func (v relatedDboFactory) NewItemDbo() SpaceItemDboWithRelatedAndIDs {
	return v.newItemDbo()
}

func (v relatedDboFactory) NewSpaceModuleDbo() dal4spaceus.SpaceModuleDbo {
	return v.newSpaceModuleDbo()
}
