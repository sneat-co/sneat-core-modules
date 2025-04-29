package facade4spaceus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"testing"
)

func TestRegisterDboFactory(t *testing.T) {
	type args struct {
		moduleID   coretypes.ModuleID
		collection string
		f          RelatedDboFactory
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "contactus/contact",
			args: args{
				moduleID:   const4contactus.ModuleID,
				collection: const4contactus.ContactsCollection,
				f: func() (dal4spaceus.SpaceModuleDbo, dal4spaceus.SpaceItemDbo, *dbo4linkage.WithRelatedAndIDs) {
					dbo := new(dbo4contactus.ContactDbo)
					return new(dbo4contactus.ContactusSpaceDbo), dbo, &dbo.WithRelatedAndIDs
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterDboFactory(tt.args.moduleID, tt.args.collection, tt.args.f)
			f := getDboFactory(tt.args.moduleID, tt.args.collection)
			if f == nil {
				t.Errorf("getDboFactory() = nil")
			}
			spaceModuleDbo, itemDbo, withRelatedAndIDs := f()
			if spaceModuleDbo == nil {
				t.Errorf("getDboFactory() spaceModuleDbo = nil")
			}
			if itemDbo == nil {
				t.Errorf("getDboFactory() itemDbo = nil")
			}
			if withRelatedAndIDs == nil {
				t.Errorf("getDboFactory() withRelatedAndIDs = nil")
			}
		})
	}
}
