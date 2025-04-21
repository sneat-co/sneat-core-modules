package dbo4linkage

import (
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"reflect"
	"slices"
	"testing"
	"time"
)

func TestWithRelatedAndIDs_SetRelationshipToItem(t *testing.T) {
	type fields struct {
		Related    RelatedByModuleID
		relatedIDs []string
	}
	type args struct {
		userID  string
		command RelationshipItemRolesCommand
		now     time.Time
	}
	now := time.Now()
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantUpdates []update.Update
	}{
		{
			name:   "set_related_as_parent_for_empty",
			fields: fields{},
			args: args{
				userID: "u1",
				command: RelationshipItemRolesCommand{
					ItemRef: SpaceModuleItemRef{
						Space:      "space1",
						Module:     const4contactus.ModuleID,
						Collection: const4contactus.ContactsCollection,
						ItemID:     "c2",
					},
					Add: &RolesCommand{
						RolesOfItem: []RelationshipRoleID{"parent"},
					},
				},
				now: now,
			},
			wantUpdates: []update.Update{
				update.ByFieldName("related.contactus.contacts", // space1.c2.relatedAs.child
					[]*RelatedItem{
						{
							Keys: []RelatedItemKey{
								{SpaceID: "space1", ItemID: "c2"},
							},
							RolesOfItem: RelationshipRoles{
								"parent": &RelationshipRole{},
							},
							RolesToItem: RelationshipRoles{
								"child": &RelationshipRole{},
							},
						},
					}),
				//{Field: "related.space1.contactus.contacts.c2.relatesAs.child", Value: &RelationshipRole{WithCreatedField: dbmodels.WithCreatedField{Created: dbmodels.Created{By: "u1", On: now.Format(time.DateTime)}}}},
				update.ByFieldName("relatedIDs", []string{
					"*",
					"s=space1",
					"m=contactus",
					"m=contactus&c=contacts",
					"s=space1&m=contactus&c=contacts",
					"s=space1&m=contactus&c=contacts&i=c2",
				}),
			},
		},
		{
			name:   "set_related_as_child_for_empty",
			fields: fields{},
			args: args{
				userID: "u1",
				command: RelationshipItemRolesCommand{
					ItemRef: SpaceModuleItemRef{
						Space:      "space1",
						Module:     const4contactus.ModuleID,
						Collection: const4contactus.ContactsCollection,
						ItemID:     "c2",
					},
					Add: &RolesCommand{
						RolesOfItem: []RelationshipRoleID{"child"},
					},
				},
				now: now,
			},
			wantUpdates: []update.Update{
				update.ByFieldName("related.contactus.contacts", // space1.c2.relatedAs.child
					[]*RelatedItem{
						{
							Keys: []RelatedItemKey{
								{SpaceID: "space1", ItemID: "c2"},
							},
							RolesOfItem: RelationshipRoles{
								"child": &RelationshipRole{},
							},
							RolesToItem: RelationshipRoles{
								"parent": &RelationshipRole{},
							},
						},
					}),
				update.ByFieldName("relatedIDs", []string{
					"*",
					"s=space1",
					"m=contactus",
					"m=contactus&c=contacts",
					"s=space1&m=contactus&c=contacts",
					"s=space1&m=contactus&c=contacts&i=c2",
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &WithRelatedAndIDs{
				WithRelated: WithRelated{
					Related: tt.fields.Related,
				},
				WithRelatedIDs: WithRelatedIDs{
					RelatedIDs: tt.fields.relatedIDs,
				},
			}
			userID := "u123"
			gotUpdates, gotErr := v.AddRelationshipAndID(
				now,
				userID,
				tt.args.command,
			)
			if gotErr != nil {
				t.Fatal(gotErr)
			}
			if len(gotUpdates) != len(tt.wantUpdates) {
				t.Errorf("SetRelationshipToItem()\nactual:\n%+v,\nwant:\n%+v", gotUpdates, tt.wantUpdates)
			}
			for i, gotUpdate := range gotUpdates {
				wantUpdate := tt.wantUpdates[i]
				if gotUpdate.FieldName() != wantUpdate.FieldName() {
					t.Errorf("SetRelationshipToItem()[%d]\nactual.Field:\n\t%+v,\nwant.Field:\n\t%+v", i, gotUpdate.FieldName(), wantUpdate.FieldName())
				}
				if !reflect.DeepEqual(gotUpdate.FieldPath(), wantUpdate.FieldPath()) {
					t.Errorf("SetRelationshipToItem()[%d]\nactual.Field:\n\t%+v,\nwant.Field:\n\t%+v", i, gotUpdate.FieldPath(), wantUpdate.FieldPath())
				}
				if gotUpdate.FieldName() == "related.contactus.contacts" {
					gotRelated := gotUpdate.Value().([]*RelatedItem)
					wantRelated := wantUpdate.Value().([]*RelatedItem)
					if len(gotRelated) != len(wantRelated) {
						t.Fatalf("expected to have %d related items, but got %d", len(wantRelated), len(gotRelated))
					}
					for i, relatedItem := range gotRelated {
						wantRelatedItem := wantRelated[i]
						for gotRoleOfItem := range relatedItem.RolesOfItem {
							if wantRelatedItem.RolesOfItem[gotRoleOfItem] == nil {
								t.Error("got unexpected roleOfItem=" + gotRoleOfItem)
							}
						}
						for gotRoleToItem := range relatedItem.RolesToItem {
							if wantRelatedItem.RolesToItem[gotRoleToItem] == nil {
								t.Error("got unexpected roleToItem=" + gotRoleToItem)
							}
						}
					}
				} else if !reflect.DeepEqual(gotUpdate.Value(), wantUpdate.Value()) {
					if gotUpdate.FieldName() == "relatedIDs" {
						// We do not care about order of the relatedIDs
						gotRelatedIDs := gotUpdate.Value().([]string)
						wantRelatedIDs := wantUpdate.Value().([]string)
						slices.Sort(gotRelatedIDs)
						slices.Sort(wantRelatedIDs)
						if !slices.Equal(gotRelatedIDs, wantRelatedIDs) {
							t.Errorf("SetRelationshipToItem()[%d] Field=%s\nactual.Value:\n\t%+v\nwant.Value:\n\t%+v", i, gotUpdate.FieldName(), gotRelatedIDs, wantRelatedIDs)
						}
					} else {
						t.Errorf("SetRelationshipToItem()[%d] Field=%s\nactual.Value:\n\t%+v\nwant.Value:\n\t%+v", i, gotUpdate.FieldName(), gotUpdate.Value(), wantUpdate.Value())

						if gotUpdate.FieldName() == "related" {
							gotItems, ok := gotUpdate.Value().([]*RelatedItem)
							if !ok {
								t.Errorf("SetRelationshipToItem()[%d]\nactual type:\n\t%T\nwant type:\n\t%T", i, gotUpdate.Value(), wantUpdate.Value())
								return
							}
							wantItems := wantUpdate.Value().([]*RelatedItem)
							if len(gotItems) != len(wantItems) {
								t.Errorf("SetRelationshipToItem()[%d]\nactual.Value:\n\t%+v,\nwant.Value:\n\t%+v", i, gotItems, wantItems)
								return
							}
							for j, gotItem := range gotItems {
								wantItem := wantItems[j]
								if !reflect.DeepEqual(gotItem.Keys, wantItem.Keys) {
									t.Errorf("SetRelationshipToItem()[%d]\nactual.Value[%d].Keys:\n\t%+v,\nwant.Value[%d].Keys:\n\t%+v", i, j, gotItem.Keys, j, wantItem.Keys)
								}
								if !reflect.DeepEqual(gotItem.RolesOfItem, wantItem.RolesOfItem) {
									t.Errorf("SetRelationshipToItem()[%d]\nactual.Value[%d].RolesOfItem:\n\t%+v,\nwant.Value[%d].RolesOfItem:\n\t%+v", i, j, gotItem.RolesOfItem, j, wantItem.RolesOfItem)
								}
								if !reflect.DeepEqual(gotItem.RolesToItem, wantItem.RolesToItem) {
									t.Errorf("SetRelationshipToItem()[%d]\nactual.Value[%d].RolesToItem:\n\t%+v,\nwant.Value[%d].RolesToItem:\n\t%+v", i, j, gotItem.RolesToItem, j, wantItem.RolesToItem)
								}
							}
						}
					}
				}
			}
		})
	}
}
