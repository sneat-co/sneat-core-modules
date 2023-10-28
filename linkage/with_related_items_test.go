package linkage

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"reflect"
	"testing"
	"time"
)

func TestWithRelatedItems_SetRelationshipToItem(t *testing.T) {
	type fields struct {
		RelatedItems   RelatedItemsByTeam
		RelatedItemIDs []string
	}
	type args struct {
		userID    string
		recordRef TeamModuleDocRef
		link      Link
		now       time.Time
	}
	now := time.Now()
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantUpdates []dal.Update
	}{
		{
			name:   "set_related_as_parent_for_empty",
			fields: fields{},
			args: args{
				userID: "u1",
				recordRef: TeamModuleDocRef{
					TeamID:     "team1",
					ModuleID:   const4contactus.ModuleID,
					Collection: const4contactus.ContactsCollection,
					ItemID:     "u1c",
				},
				link: Link{
					TeamModuleDocRef: TeamModuleDocRef{
						TeamID:     "team1",
						ModuleID:   const4contactus.ModuleID,
						Collection: const4contactus.ContactsCollection,
						ItemID:     "c2",
					},
					RelatedAs: []RelationshipID{"parent"},
				},
				now: now,
			},
			wantUpdates: []dal.Update{
				{Field: "relatedItems.team1.contactus.contacts.c2.relatedAs.parent", Value: &Relationship{WithCreatedField: dbmodels.WithCreatedField{Created: dbmodels.Created{By: "u1", On: now.Format(time.DateOnly)}}}},
				//{Field: "relatedItems.team1.contactus.contacts.c2.relatesAs.child", Value: &Relationship{WithCreatedField: dbmodels.WithCreatedField{Created: dbmodels.Created{By: "u1", On: now.Format(time.DateOnly)}}}},
				{Field: "relatedItemIDs", Value: []string{"team1.contactus.contacts.c2"}},
			},
		},
		{
			name:   "set_related_as_child_for_empty",
			fields: fields{},
			args: args{
				userID: "u1",
				recordRef: TeamModuleDocRef{
					TeamID:     "team1",
					ModuleID:   const4contactus.ModuleID,
					Collection: const4contactus.ContactsCollection,
					ItemID:     "u1c",
				},
				link: Link{
					TeamModuleDocRef: TeamModuleDocRef{
						TeamID:     "team1",
						ModuleID:   const4contactus.ModuleID,
						Collection: const4contactus.ContactsCollection,
						ItemID:     "c2",
					},
					RelatedAs: []RelationshipID{"child"},
				},
				now: now,
			},
			wantUpdates: []dal.Update{
				{Field: "relatedItems.team1.contactus.contacts.c2.relatedAs.child",
					Value: &Relationship{
						WithCreatedField: dbmodels.WithCreatedField{
							Created: dbmodels.Created{By: "u1",
								On: now.Format(time.DateOnly),
							},
						},
					},
				},
				{Field: "relatedItemIDs", Value: []string{"team1.contactus.contacts.c2"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &WithRelatedItems{
				RelatedItems:   tt.fields.RelatedItems,
				RelatedItemIDs: tt.fields.RelatedItemIDs,
			}
			gotUpdates, gotErr := v.SetRelationshipToItem(
				tt.args.userID,
				tt.args.recordRef,
				tt.args.link,
				tt.args.now,
			)
			if gotErr != nil {
				t.Fatal(gotErr)
			}
			if !reflect.DeepEqual(gotUpdates, tt.wantUpdates) {
				t.Errorf("SetRelationshipToItem() = \n%+v,\nwant:\n%+v", gotUpdates, tt.wantUpdates)
			}
		})
	}
}
