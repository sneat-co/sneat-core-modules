package facade4linkage

import (
	"testing"

	"github.com/sneat-co/sneat-core-modules/contactusmodels/const4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

// fakeSpaceItemDbo and fakeSpaceModuleDbo are test doubles so this test does not
// depend on the contactus module (cycle-break). They satisfy the minimal interfaces
// the linkage dbo factory requires.
type fakeSpaceItemDbo struct{}

func (fakeSpaceItemDbo) Validate() error { return nil }

func (fakeSpaceItemDbo) RelatedAndIDs() *dbo4linkage.WithRelatedAndIDs {
	return &dbo4linkage.WithRelatedAndIDs{}
}

type fakeSpaceModuleDbo struct{}

func (fakeSpaceModuleDbo) Validate() error { return nil }

func TestRegisterDboFactory(t *testing.T) {
	type args struct {
		moduleID   coretypes.ExtID
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
				moduleID:   const4contactus.ExtensionID,
				collection: const4contactus.ContactsCollection,
				f: NewDboFactory(
					func() SpaceItemDboWithRelatedAndIDs {
						return new(fakeSpaceItemDbo)
					},
					func() dal4spaceus.SpaceModuleDbo {
						return new(fakeSpaceModuleDbo)
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
