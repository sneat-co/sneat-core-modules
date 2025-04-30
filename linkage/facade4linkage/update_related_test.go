package facade4linkage

import (
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dbo4contactus"
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
				f: NewDboFactory(
					func() SpaceItemDboWithRelatedAndIDs {
						return new(dbo4contactus.ContactDbo)
					},
					func() dal4spaceus.SpaceModuleDbo {
						return new(dbo4contactus.ContactusSpaceDbo)
					},
				),
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
			if f.NewSpaceModuleDbo() == nil {
				t.Errorf("NewSpaceModuleDbo() = nil")
			}
			if f.NewItemDbo() == nil {
				t.Errorf("NewItemDbo() = nil")
			}
		})
	}
}
