package dbo4linkage

import (
	"encoding/json"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/strongo/strongoapp/with"
)

func TestWithRelatedAndIDs_SetRelationshipToItem(t *testing.T) {
	type fields struct {
		Related    RelatedModules
		relatedIDs []string
	}
	type args struct {
		userID  string
		spaceID coretypes.SpaceID
		command RelationshipItemRolesCommand
		now     time.Time
	}
	now := time.Now()
	createdField := with.CreatedField{
		Created: with.Created{
			At: now.Format(time.RFC3339),
			By: "u1",
		},
	}
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
				userID:  "u1",
				spaceID: "space1",
				command: RelationshipItemRolesCommand{
					ItemRef: ItemRef{
						ExtID:      const4contactus.ExtensionID,
						Collection: const4contactus.ContactsCollection,
						ItemID:     "contact1",
					},
					Add: &RolesCommand{
						RolesOfItem: []RelationshipRoleID{"parent"},
					},
				},
				now: now,
			},
			wantUpdates: []update.Update{
				update.ByFieldPath([]string{"related", "contactus", "contacts", "contact1"}, // space1.c2.relatedAs.child
					&RelatedItem{
						RolesOfItem: RelationshipRoles{
							"parent": &RelationshipRole{
								CreatedField: createdField,
							},
						},
						RolesToItem: RelationshipRoles{
							"child": &RelationshipRole{
								CreatedField: createdField,
							},
						},
					}),
				//{Field: "related.space1.contactus.contacts.c2.relatesAs.child", Value: &RelationshipRole{WithCreatedField: dbmodels.WithCreatedField{Created: dbmodels.Created{By: "u1", On: now.Format(time.DateTime)}}}},
				update.ByFieldName("relatedIDs", []string{
					"*",
					"m=contactus&c=contacts&s=space1&i=contact1",
				}),
			},
		},
		{
			name:   "set_related_as_child_for_empty",
			fields: fields{},
			args: args{
				userID:  "u1",
				spaceID: "space1",
				command: RelationshipItemRolesCommand{
					ItemRef: ItemRef{
						ExtID:      const4contactus.ExtensionID,
						Collection: const4contactus.ContactsCollection,
						ItemID:     "contact1",
					},
					Add: &RolesCommand{
						RolesOfItem: []RelationshipRoleID{"child"},
					},
				},
				now: now,
			},
			wantUpdates: []update.Update{
				update.ByFieldPath([]string{"related", "contactus", "contacts", "contact1"}, // space1.c2.relatedAs.child
					&RelatedItem{
						RolesOfItem: RelationshipRoles{
							"child": &RelationshipRole{
								CreatedField: createdField,
							},
						},
						RolesToItem: RelationshipRoles{
							"parent": &RelationshipRole{
								CreatedField: createdField,
							},
						},
					}),
				update.ByFieldName("relatedIDs", []string{
					"*",
					"m=contactus&c=contacts&s=space1&i=contact1",
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
			gotUpdates, gotErr := v.AddRelationshipAndID(
				now,
				tt.args.userID,
				tt.args.spaceID,
				tt.args.command,
			)
			if gotErr != nil {
				t.Fatal(gotErr)
			}
			if len(gotUpdates) != len(tt.wantUpdates) {
				t.Errorf("SetRelationshipToItem()\n\t got:\n\t\t%+v,\n\twant:\n\t\t%+v", gotUpdates, tt.wantUpdates)
			}
			for i, gotUpdate := range gotUpdates {
				wantUpdate := tt.wantUpdates[i]
				if updateFieldName := gotUpdate.FieldName(); updateFieldName != "" {
					if updateFieldName != wantUpdate.FieldName() {
						t.Errorf("SetRelationshipToItem()[%d]\nactual.FieldName:\n\t%+v,\nwant.FieldName:\n\t%+v",
							i, updateFieldName, wantUpdate.FieldName())
					}
				} else if !slices.Equal(gotUpdate.FieldPath(), wantUpdate.FieldPath()) {
					t.Errorf("SetRelationshipToItem()[%d]\nactual.FieldPath:\n\t%+v,\nwant.FieldPath:\n\t%+v",
						i, gotUpdate.FieldPath(), wantUpdate.FieldPath())
				}

				if !reflect.DeepEqual(gotUpdate.FieldPath(), wantUpdate.FieldPath()) {
					t.Errorf("SetRelationshipToItem()[%d]\nactual.Field:\n\t%+v,\nwant.Field:\n\t%+v",
						i, gotUpdate.FieldPath(), wantUpdate.FieldPath())
				}
				if gotUpdate.FieldName() == "related.contactus.contacts" {
					gotRelated := gotUpdate.Value().(map[string]*RelatedItem)
					wantRelated := wantUpdate.Value().(map[string]*RelatedItem)
					if len(gotRelated) != len(wantRelated) {
						t.Fatalf("expected to have %d related items, but got %d", len(wantRelated), len(gotRelated))
					}
					for relatedItemID, relatedItem := range gotRelated {
						wantRelatedItem := wantRelated[relatedItemID]
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
						// We do not care about the order of the relatedIDs
						gotRelatedIDs := gotUpdate.Value().([]string)
						wantRelatedIDs := wantUpdate.Value().([]string)
						slices.Sort(gotRelatedIDs)
						slices.Sort(wantRelatedIDs)
						if !slices.Equal(gotRelatedIDs, wantRelatedIDs) {
							t.Errorf("SetRelationshipToItem()[%d] Field=%s\nactual.Value:\n\t%+v\nwant.Value:\n\t%+v",
								i, gotUpdate.FieldName(), gotRelatedIDs, wantRelatedIDs)
						}
					} else if fieldPath := gotUpdate.FieldPath(); len(fieldPath) > 0 && fieldPath[0] == "related" {
						if gotItem, ok := gotUpdate.Value().(*RelatedItem); !ok {
							t.Errorf("SetRelationshipToItem()[%d]\nactual type:\n\t%T\nwant type:\n\t%T", i, gotUpdate.Value(), wantUpdate.Value())
							return
						} else if wantItem, ok2 := wantUpdate.Value().(*RelatedItem); !ok2 {
							t.Errorf("SetRelationshipToItem()[%d]\nactual type:\n\t%T\nwant type:\n\t%T", i, gotUpdate.Value(), wantUpdate.Value())
							return
						} else if !reflect.DeepEqual(gotItem.RolesOfItem, wantItem.RolesOfItem) {
							gotJson, _ := json.Marshal(gotItem.RolesOfItem)
							wantJson, _ := json.Marshal(wantItem.RolesOfItem)
							t.Errorf("SetRelationshipToItem()[%d]\nactual RolesOfItem:\n\t%s,\nwant RolesOfItem:\n\t%s", i, string(gotJson), wantJson)
						} else if !reflect.DeepEqual(gotItem.RolesToItem, wantItem.RolesToItem) {
							t.Errorf("SetRelationshipToItem()[%d]\nactual RolesToItem:\n\t%+v,\nwant RolesToItem:\n\t%+v", i, gotItem.RolesToItem, wantItem.RolesToItem)
						}
					}
				}
			}
		})
	}
}

func TestWithRelated_RemoveRelatedItem(t *testing.T) {
	type fields struct {
		Related RelatedModules
	}
	type args struct {
		ref ItemRef
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantRelated RelatedModules
		wantUpdates []update.Update
	}{
		{
			name: "remove_the_only_related_item",
			fields: fields{
				Related: RelatedModules{
					"module1": {
						"collection1": {
							"item1": {},
						},
					},
				},
			},
			args: args{
				ref: ItemRef{
					ExtID:      "module1",
					Collection: "collection1",
					ItemID:     "item1",
				},
			},
			wantUpdates: []update.Update{update.ByFieldPath([]string{"related"}, update.DeleteField)},
			wantRelated: nil,
		},
		{
			name: "remove_the_1_of_2_related_item_in_same_collection",
			fields: fields{
				Related: RelatedModules{
					"module1": {
						"collection1": {
							"item1": {},
							"item2": {},
						},
					},
				},
			},
			args: args{
				ref: ItemRef{
					ExtID:      "module1",
					Collection: "collection1",
					ItemID:     "item1",
				},
			},
			wantUpdates: []update.Update{update.ByFieldPath([]string{"related", "module1", "collection1", "item1"}, update.DeleteField)},
			wantRelated: RelatedModules{
				"module1": {
					"collection1": {
						"item2": {},
					},
				},
			},
		},
		{
			name: "remove_1_of_2_related_item_in_different_collection",
			fields: fields{
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
			args: args{
				ref: ItemRef{
					ExtID:      "module1",
					Collection: "collection1",
					ItemID:     "item1",
				},
			},
			wantUpdates: []update.Update{update.ByFieldPath([]string{"related", "module1", "collection1"}, update.DeleteField)},
			wantRelated: RelatedModules{
				"module1": {
					"collection2": {
						"item2": {},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &WithRelated{
				Related: tt.fields.Related,
			}
			if gotUpdates := v.RemoveRelatedItem(tt.args.ref); !reflect.DeepEqual(gotUpdates, tt.wantUpdates) {
				t.Errorf("RemoveRelatedItem() = \n\tgot updates=%v, \n\twant updates=%v", gotUpdates, tt.wantUpdates)
			}
			if !reflect.DeepEqual(v.Related, tt.wantRelated) {
				t.Errorf("RemoveRelatedItem() => \n\t got related=%v, \n\twant related=%v", v.Related, tt.wantRelated)
			}
		})
	}
}

func TestWithRelated_removeRolesFromRelatedItem(t *testing.T) {
	type fields struct {
		Related RelatedModules
	}
	type args struct {
		itemRef ItemRef
		remove  RolesCommand
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantUpdates []update.Update
		wantRelated RelatedModules
	}{
		{
			name: "remove_roles_from_nil_related",
			fields: fields{
				Related: nil,
			},
			args: args{
				itemRef: ItemRef{
					ExtID:      "module1",
					Collection: "collection1",
					ItemID:     "item1",
				},
				remove: RolesCommand{
					RolesOfItem: []RelationshipRoleID{"child"},
					RolesToItem: []RelationshipRoleID{"parent"},
				},
			},
			wantUpdates: nil,
			wantRelated: nil,
		},
		{
			name: "remove_roles_from_single_related_item_with_single_reciprocal_role",
			fields: fields{
				Related: RelatedModules{
					"module1": {
						"collection1": {
							"item1": {
								RolesOfItem: RelationshipRoles{"child": &RelationshipRole{}},
								RolesToItem: RelationshipRoles{"parent": &RelationshipRole{}},
							},
						},
					},
				},
			},
			args: args{
				itemRef: ItemRef{
					ExtID:      "module1",
					Collection: "collection1",
					ItemID:     "item1",
				},
				remove: RolesCommand{
					RolesOfItem: []RelationshipRoleID{"child"},
					RolesToItem: []RelationshipRoleID{"parent"},
				},
			},
			wantUpdates: []update.Update{update.ByFieldPath(update.FieldPath{relatedField}, update.DeleteField)},
			wantRelated: nil,
		},
		{
			name: "remove_1_reciprocal_role_from_single_related_item_with_2_reciprocal_roles",
			fields: fields{
				Related: RelatedModules{
					"module1": {
						"collection1": {
							"item1": {
								RolesOfItem: RelationshipRoles{
									"child":  &RelationshipRole{}, // reciprocal for parent
									"spouse": &RelationshipRole{}, // reciprocal for spouse
								},
								RolesToItem: RelationshipRoles{
									"parent": &RelationshipRole{}, // reciprocal for child
									"spouse": &RelationshipRole{}, // reciprocal for spouse
								},
							},
						},
					},
				},
			},
			args: args{
				itemRef: ItemRef{
					ExtID:      "module1",
					Collection: "collection1",
					ItemID:     "item1",
				},
				remove: RolesCommand{
					RolesOfItem: []RelationshipRoleID{"child"},
					RolesToItem: []RelationshipRoleID{"parent"},
				},
			},
			wantUpdates: []update.Update{
				update.ByFieldPath(
					[]string{relatedField, "module1", "collection1", "item1", "rolesOfItem", "child"},
					update.DeleteField,
				),
				update.ByFieldPath(
					[]string{relatedField, "module1", "collection1", "item1", "rolesToItem", "parent"},
					update.DeleteField,
				),
			},
			wantRelated: RelatedModules{
				"module1": {
					"collection1": {
						"item1": {
							RolesOfItem: RelationshipRoles{
								"spouse": &RelationshipRole{},
							},
							RolesToItem: RelationshipRoles{
								"spouse": &RelationshipRole{},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &WithRelated{
				Related: tt.fields.Related,
			}
			gotUpdates := v.removeRolesFromRelatedItem(tt.args.itemRef, tt.args.remove)
			if !reflect.DeepEqual(gotUpdates, tt.wantUpdates) {
				t.Errorf("removeRolesFromRelatedItem():\n\t got updates=%+v,\n\twant updates=%+v", gotUpdates, tt.wantUpdates)
			}
			if !reflect.DeepEqual(v.Related, tt.wantRelated) {
				t.Errorf("removeRolesFromRelatedItem():\n\t got related=%+v,\n\twant related=%+v", v.Related, tt.wantRelated)
			}
		})
	}
}
