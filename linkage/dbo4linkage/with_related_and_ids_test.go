package dbo4linkage

import (
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"reflect"
	"slices"
	"testing"
	"time"
)

func TestAddRelationshipAndID(t *testing.T) {
	type args struct {
		now               time.Time
		userID            string
		spaceID           coretypes.SpaceID
		withRelatedAndIDs *WithRelatedAndIDs
		command           RelationshipItemRolesCommand
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
				now:               now,
				userID:            "user1",
				spaceID:           "space1",
				withRelatedAndIDs: &WithRelatedAndIDs{},
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
			gotUpdates, err := AddRelationshipAndID(tt.args.now, tt.args.userID, tt.args.spaceID,
				&tt.args.withRelatedAndIDs.WithRelated, &tt.args.withRelatedAndIDs.WithRelatedIDs,
				tt.args.command)
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
		spaceID           coretypes.SpaceID
		withRelatedAndIDs *WithRelatedAndIDs
		ref               ItemRef
	}
	tests := []struct {
		name        string
		args        args
		wantUpdates []update.Update
	}{
		{
			name: "remove_non_existing_item",
			args: args{
				spaceID:           "space1",
				withRelatedAndIDs: &WithRelatedAndIDs{},
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
			if gotUpdates := RemoveRelatedAndID(tt.args.spaceID,
				&tt.args.withRelatedAndIDs.WithRelated, &tt.args.withRelatedAndIDs.WithRelatedIDs,
				tt.args.ref); !reflect.DeepEqual(gotUpdates, tt.wantUpdates) {
				t.Errorf("RemoveRelatedAndID() = %v, want %v", gotUpdates, tt.wantUpdates)
			}
		})
	}
}

func TestUpdateRelatedIDs(t *testing.T) {
	type args struct {
		spaceID           coretypes.SpaceID
		withRelatedAndIDs *WithRelatedAndIDs
	}
	tests := []struct {
		name           string
		args           args
		wantRelatedIDs []string
		wantUpdates    []update.Update
	}{
		{
			name: "empty",
			args: args{
				spaceID:           "space1",
				withRelatedAndIDs: &WithRelatedAndIDs{},
			},
			wantRelatedIDs: []string{"-"},
		},
		{
			name: "single_related_empty_ids",
			args: args{
				spaceID: "space1",
				withRelatedAndIDs: &WithRelatedAndIDs{
					WithRelated: WithRelated{
						Related: RelatedModules{
							"module1": {
								"collection1": {
									"item1": {},
								},
							},
						},
					},
				},
			},
			wantRelatedIDs: []string{
				"*",
				//"m=module1",
				//"m=module1&c=collection1",
				"m=module1&c=collection1&s=space1&i=item1",
			},
		},
		{
			name: "2_related_same_space_and_collection",
			args: args{
				spaceID: "space1",
				withRelatedAndIDs: &WithRelatedAndIDs{
					WithRelated: WithRelated{
						Related: RelatedModules{
							"module1": {
								"collection1": {
									"item1": {},
									"item2": {},
								},
							},
						},
					},
				},
			},
			wantRelatedIDs: []string{
				"*",
				//"m=module1",
				//"m=module1&c=collection1",
				"m=module1&c=collection1&s=space1&i=item1",
				"m=module1&c=collection1&s=space1&i=item2",
			},
		},
		{
			name: "2_related_same_collection_different_spaces",
			args: args{
				spaceID: "space1",
				withRelatedAndIDs: &WithRelatedAndIDs{
					WithRelated: WithRelated{
						Related: RelatedModules{
							"module1": {
								"collection1": {
									"item1":        {},
									"item2@space2": {},
								},
							},
						},
					},
				},
			},
			wantRelatedIDs: []string{
				"*",
				"s=space2", // add only for different spaces
				//"m=module1",
				//"m=module1&c=collection1",
				//"m=module1&c=collection1&s=space2", // add only for different spaces
				"m=module1&c=collection1&s=space1&i=item1",
				"m=module1&c=collection1&s=space2&i=item2",
			},
		},
		{
			name: "2_related_same_space_different_collections",
			args: args{
				spaceID: "space1",
				withRelatedAndIDs: &WithRelatedAndIDs{
					WithRelated: WithRelated{
						Related: RelatedModules{
							"module1": {
								"collection1": {
									"item1": {},
								},
								"collection2": {
									"item2": {},
								},
							},
						},
					},
				},
			},
			wantRelatedIDs: []string{
				"*",
				//"m=module1",
				//"m=module1&c=collection1",
				//"m=module1&c=collection2",
				"m=module1&c=collection1&s=space1&i=item1",
				"m=module1&c=collection2&s=space1&i=item2",
			},
		},
		{
			name: "2_related_different_spaces_and_different_collections",
			args: args{
				spaceID: "space1",
				withRelatedAndIDs: &WithRelatedAndIDs{
					WithRelated: WithRelated{
						Related: RelatedModules{
							"module1": {
								"collection1": {
									"item1": {},
								},
								"collection2": {
									"item2@space2": {},
								},
							},
						},
					},
				},
			},
			wantRelatedIDs: []string{
				"*",
				"s=space2", // add only for different spaces
				//"m=module1",
				//"m=module1&c=collection1",
				//"m=module1&c=collection2",
				//"m=module1&c=collection2&s=space2", // add only for different spaces
				"m=module1&c=collection1&s=space1&i=item1",
				"m=module1&c=collection2&s=space2&i=item2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Perform the tested code
			gotUpdates := UpdateRelatedIDs(tt.args.spaceID, &tt.args.withRelatedAndIDs.WithRelated, &tt.args.withRelatedAndIDs.WithRelatedIDs)

			// Assert the results
			if len(tt.args.withRelatedAndIDs.Related) == 0 {
				if tt.args.withRelatedAndIDs.RelatedIDs[0] != "-" {
					t.Errorf("first element of relatedIDs should be '-'")
				}
			} else if tt.args.withRelatedAndIDs.RelatedIDs[0] != "*" {
				t.Errorf("first element of relatedIDs should be '*'")
			}

			//sort.StringSlice(tt.wantRelatedIDs).Sort()
			slices.Sort(tt.wantRelatedIDs)
			if !slices.Equal(tt.args.withRelatedAndIDs.RelatedIDs, tt.wantRelatedIDs) {
				t.Errorf("UpdateRelatedIDs() = got\n\trelatedIDs:\n\t\t%+v\n\twant:\n\t\t%+v", tt.args.withRelatedAndIDs.RelatedIDs, tt.wantRelatedIDs)
			}
			if tt.wantUpdates != nil {
				if !reflect.DeepEqual(gotUpdates, tt.wantUpdates) {
					t.Errorf("UpdateRelatedIDs() = %v, want %v", gotUpdates, tt.wantUpdates)
				}
			}
		})
	}
}
