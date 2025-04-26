package dbo4linkage

import (
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"reflect"
	"testing"
	"time"
)

func TestAddRelationshipAndID(t *testing.T) {
	type args struct {
		now            time.Time
		userID         string
		spaceID        coretypes.SpaceID
		withRelated    *WithRelated
		withRelatedIDs *WithRelatedIDs
		command        RelationshipItemRolesCommand
	}
	now := time.Now()
	tests := []struct {
		name        string
		args        args
		wantUpdates []update.Update
		wantErr     bool
	}{
		{
			name: "add",
			args: args{
				now:            now,
				userID:         "user1",
				spaceID:        "space1",
				withRelated:    &WithRelated{},
				withRelatedIDs: &WithRelatedIDs{},
				command: RelationshipItemRolesCommand{
					ItemRef: ItemRef{
						Module:     "module1",
						Collection: "collection1",
						ItemID:     "item1",
					},
					Add: &RolesCommand{
						RolesOfItem: []RelationshipRoleID{
							"role1",
						},
					},
				},
			},
			wantUpdates: []update.Update{ // TODO(help-wanted): Fix this test
				update.ByFieldName("related.module1.collection1", RelatedItems{
					"item1": {},
				}),
				update.ByFieldName("relatedIds", []string{
					"*",
					"module1.collection1",
					"module1.collection1.item1",
				}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUpdates, err := AddRelationshipAndID(tt.args.now, tt.args.userID, tt.args.spaceID, tt.args.withRelated, tt.args.withRelatedIDs, tt.args.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddRelationshipAndID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUpdates == nil {
				t.Errorf("AddRelationshipAndID() gotUpdates = nil, want %v", tt.wantUpdates)
				return
			}
			//if !reflect.DeepEqual(gotUpdates, tt.wantUpdates) {
			//	t.Errorf("AddRelationshipAndID() gotUpdates = %v, want %v", gotUpdates, tt.wantUpdates)
			//}
		})
	}
}

func TestRemoveRelatedAndID(t *testing.T) {
	type args struct {
		spaceID        coretypes.SpaceID
		withRelated    *WithRelated
		withRelatedIDs *WithRelatedIDs
		ref            ItemRef
	}
	tests := []struct {
		name        string
		args        args
		wantUpdates []update.Update
	}{
		{
			name: "remove_non_existing_item",
			args: args{
				spaceID:        "space1",
				withRelated:    &WithRelated{},
				withRelatedIDs: &WithRelatedIDs{},
				ref: ItemRef{
					Module:     "module1",
					Collection: "collection1",
					ItemID:     "item1",
				},
			},
			wantUpdates: []update.Update{
				update.ByFieldName("relatedIDs", []string{"-"}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotUpdates := RemoveRelatedAndID(tt.args.spaceID, tt.args.withRelated, tt.args.withRelatedIDs, tt.args.ref); !reflect.DeepEqual(gotUpdates, tt.wantUpdates) {
				t.Errorf("RemoveRelatedAndID() = %v, want %v", gotUpdates, tt.wantUpdates)
			}
		})
	}
}
