package facade4linkage

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-core-modules/linkage/dto4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	"reflect"
	"testing"
)

func TestUpdateItemRelationships(t *testing.T) {
	type args struct {
		ctx     facade.ContextWithUser
		request dto4linkage.UpdateItemRequest
	}
	const space1ID = "space_1"
	const item1ID = "item_1"
	const collection1ID = "collection_1"
	const module1ID = "module_1"

	facade.GetSneatDB = func(ctx context.Context) (dal.DB, error) {
		return nil, nil
	}

	tests := []struct {
		name      string
		args      args
		wantItem  record.DataWithID[string, *dbo4linkage.WithRelatedAndIDsAndUserID]
		wantErr   bool
		wantPanic bool
	}{
		{
			name:      "should_update_contact_with_reciprocal_role",
			wantPanic: true, // TODO: Fix this test
			args: args{
				ctx: facade.NewContextWithUserID(context.Background(), "123"),
				request: dto4linkage.UpdateItemRequest{
					SpaceRequest: dto4spaceus.SpaceRequest{
						SpaceID: space1ID,
					},
					ItemRef: dbo4linkage.ItemRef{
						ExtID:      const4contactus.ExtensionID,
						Collection: const4contactus.ContactsCollection,
						ItemID:     item1ID,
					},
					UpdateRelatedFieldRequest: dto4linkage.UpdateRelatedFieldRequest{
						Related: []dbo4linkage.RelationshipItemRolesCommand{
							{
								ItemRef: dbo4linkage.NewItemRefSameSpace(module1ID, collection1ID, item1ID),
								Add: &dbo4linkage.RolesCommand{
									RolesOfItem: []dbo4linkage.RelationshipRoleID{
										dbo4linkage.RelationshipRoleSpouse,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("UpdateItemRelationships() did not panic")
					}
				}()
			}
			gotItem, err := UpdateItemRelationships(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateItemRelationships() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotItem, tt.wantItem) {
				t.Errorf("UpdateItemRelationships() gotItem = %v, want %v", gotItem, tt.wantItem)
			}
		})
	}
}
