package models4contactus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"reflect"
	"testing"
	"time"
)

func TestWithRelatedContacts_SetSingleRelationshipToContact(t *testing.T) {
	type fields struct {
		RelatedContacts   RelatedContacts
		RelatedContactIDs []string
	}
	type args struct {
		userID           string
		currentContactID string
		relatedContactID string
		relatedAs        ContactRelationshipID
		now              time.Time
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
				userID:           "u1",
				currentContactID: "u1c",
				relatedContactID: "c2",
				relatedAs:        "parent",
				now:              now,
			},
			wantUpdates: []dal.Update{
				{Field: "relatedContacts.c2.relatedAs.parent", Value: &ContactRelationship{WithCreatedField: dbmodels.WithCreatedField{Created: dbmodels.Created{By: "u1", On: now.Format(time.DateOnly)}}}},
				{Field: "relatedContacts.c2.relatesAs.child", Value: &ContactRelationship{WithCreatedField: dbmodels.WithCreatedField{Created: dbmodels.Created{By: "u1", On: now.Format(time.DateOnly)}}}},
				{Field: "relatedContactIDs", Value: []string{"c2"}},
			},
		},
		{
			name:   "set_related_as_child_for_empty",
			fields: fields{},
			args: args{
				userID:           "u1",
				currentContactID: "u1c",
				relatedContactID: "c2",
				relatedAs:        "child",
				now:              now,
			},
			wantUpdates: []dal.Update{
				{Field: "relatedContacts.c2.relatedAs.child", Value: &ContactRelationship{WithCreatedField: dbmodels.WithCreatedField{Created: dbmodels.Created{By: "u1", On: now.Format(time.DateOnly)}}}},
				{Field: "relatedContactIDs", Value: []string{"c2"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &WithRelatedContacts{
				RelatedContacts:   tt.fields.RelatedContacts,
				RelatedContactIDs: tt.fields.RelatedContactIDs,
			}
			gotUpdates, gotErr := v.SetSingleRelationshipToContact(tt.args.userID, tt.args.currentContactID, tt.args.relatedContactID, tt.args.relatedAs, tt.args.now)
			if gotErr != nil {
				t.Fatal(gotErr)
			}
			if !reflect.DeepEqual(gotUpdates, tt.wantUpdates) {
				t.Errorf("SetSingleRelationshipToContact() = \n%+v,\nwant:\n%+v", gotUpdates, tt.wantUpdates)
			}
		})
	}
}
