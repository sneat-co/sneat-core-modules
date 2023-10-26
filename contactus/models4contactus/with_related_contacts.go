package models4contactus

import (
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
	"strings"
	"time"
)

type ContactRelationshipID = string

type ContactRelationship struct {
	dbmodels.WithCreated
}

func (v ContactRelationship) Validate() error {
	if err := v.WithCreated.Validate(); err != nil {
		return err
	}
	return nil
}

type Relationships = map[ContactRelationshipID]*ContactRelationship

type RelatedContact struct {

	// RelatedAs - if related contact is a child of the current contact, then relatedAs = {"child": ...}
	RelatedAs Relationships

	// RelatesAs - if related contact is a child of the current contact, then relatesAs = {"parent": ...}
	RelatesAs Relationships
}

func NewRelatedContact() *RelatedContact {
	return &RelatedContact{
		RelatedAs: make(Relationships, 1),
		RelatesAs: make(Relationships, 1),
	}
}

func (v *RelatedContact) Validate() error {
	if err := v.validateRelationships(v.RelatedAs); err != nil {
		return validation.NewErrBadRecordFieldValue("relatedAs", err.Error())
	}
	if err := v.validateRelationships(v.RelatesAs); err != nil {
		return validation.NewErrBadRecordFieldValue("relatesAs", err.Error())
	}
	return nil
}

func (*RelatedContact) validateRelationships(related Relationships) error {
	for relationshipID, relationshipDetails := range related {
		if strings.TrimSpace(relationshipID) == "" {
			return errors.New("key is empty string")
		}
		if err := relationshipDetails.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type RelatedContacts = map[string]*RelatedContact

// WithRelatedContacts defines relationship of the current contact record to other contacts.
type WithRelatedContacts struct {
	/* Example of relatedContats field as a JSON:

	Contact(id="child1") {
		relatedContacts: {
			"parent1": {
				relatesAs: {
					"child": {}
				},
				relatedAs: {
					"parent": {}
				}
			},
		}
	}
	*/

	// RelatedContacts defines relationship of the current contact to other contacts. Key is contact ID.
	RelatedContacts RelatedContacts `json:"relatedContacts,omitempty" firestore:"relatedContacts,omitempty"`

	// RelatedContactIDs is a list of contact IDs that are related to the current contact - needed for indexed search.
	RelatedContactIDs []string `json:"relatedContactIDs,omitempty" firestore:"relatedContactIDs,omitempty"`
}

func (v *WithRelatedContacts) SetSingleRelationshipToContact(
	userID, currentContactID, relatedContactID string,
	relatedAs ContactRelationshipID, // if contact is a child of the user, then relatedAs = "child"
	now time.Time,
) (updates []dal.Update, err error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("argument 'userID' is empty string")
	}
	if strings.TrimSpace(currentContactID) == "" {
		return nil, errors.New("argument 'currentContactID' is empty string")
	}
	if strings.TrimSpace(relatedContactID) == "" {
		return nil, errors.New("argument 'relatedContactID' is empty string")
	}
	if strings.TrimSpace(relatedAs) == "" {
		return nil, errors.New("argument 'relatedAs' is empty string")
	}

	var relatesAs ContactRelationshipID // if related contact is a child of the current contact, then relatesAs = "parent"

	switch relatedAs {
	case "parent":
		relatesAs = "child"
	case "spouse":
		relatesAs = "spouse"
	}

	relatedContact := v.RelatedContacts[relatedContactID]
	if relatedContact == nil {
		relatedContact = NewRelatedContact()
		if v.RelatedContacts == nil {
			v.RelatedContacts = make(RelatedContacts, 1)
		}
		v.RelatedContacts[relatedContactID] = relatedContact
	}

	var alreadyHasRelatedAs, alreadyHasRelatesAs bool

	for relationshipID := range relatedContact.RelatedAs {
		if relationshipID == relatedAs {
			alreadyHasRelatedAs = true
		} else {
			updates = append(updates, dal.Update{Field: fmt.Sprintf("relatedContacts.%s.relatedAs.%s", relatedContactID, relationshipID), Value: dal.DeleteField})
		}
	}

	if relatesAs != "" {
		for relationshipID := range relatedContact.RelatesAs {
			if relationshipID == relatedAs {
				alreadyHasRelatesAs = true
			} else {
				updates = append(updates, dal.Update{Field: fmt.Sprintf("relatedContacts.%s.relatesAs.%s", relatedContactID, relationshipID), Value: dal.DeleteField})
			}
		}
	}

	if alreadyHasRelatedAs && alreadyHasRelatesAs {
		return updates, nil
	}

	updates = append(updates, v.AddRelationshipToContact(userID, currentContactID, relatedContactID, relatedAs, relatesAs, now)...)

	return updates, nil
}

func (v *WithRelatedContacts) AddRelationshipToContact(userID, userContactID string, relatedContactID string, relatedAs, relatesAs ContactRelationshipID, now time.Time) (updates []dal.Update) {
	if v.RelatedContacts == nil {
		v.RelatedContacts = make(RelatedContacts, 1)
	}
	relatedContact := v.RelatedContacts[relatedContactID]
	if relatedContact == nil {
		relatedContact = NewRelatedContact()
		v.RelatedContacts[relatedContactID] = relatedContact
	}

	addRelationship := func(field string, relationshipID ContactRelationshipID, relationships Relationships) Relationships {
		if relationships == nil {
			relationships = make(Relationships, 1)
		}
		if relationship, ok := relationships[relatedAs]; !ok {
			relationship = &ContactRelationship{
				WithCreated: dbmodels.WithCreated{
					CreatedBy: userID,
					CreatedAt: now,
				},
			}
			relationships[relatedAs] = relationship
			updates = append(updates, dal.Update{
				Field: field + "." + relationshipID,
				Value: relationship,
			})
		}
		return relationships
	}

	relatedContact.RelatedAs = addRelationship(fmt.Sprintf("relatedContacts.%s.relatedAs", relatedContactID), relatedAs, relatedContact.RelatedAs)
	if relatesAs != "" {
		relatedContact.RelatesAs = addRelationship(fmt.Sprintf("relatedContacts.%s.relatesAs", relatedContactID), relatesAs, relatedContact.RelatesAs)
	}

	if !slice.Contains(v.RelatedContactIDs, relatedContactID) {
		v.RelatedContactIDs = append(v.RelatedContactIDs, relatedContactID)
		updates = append(updates, dal.Update{
			Field: "relatedContactIDs",
			Value: v.RelatedContactIDs,
		})
	}
	return updates
}

// RemoveRelationshipToContact removes all relationships to a given contact
func (v *WithRelatedContacts) RemoveRelationshipToContact(contactID string) (updates []dal.Update) {
	if _, ok := v.RelatedContacts[contactID]; ok {
		delete(v.RelatedContacts, contactID)
		updates = append(updates, dal.Update{
			Field: "relatedContacts." + contactID,
			Value: dal.DeleteField,
		})
		v.RelatedContactIDs = slice.RemoveInPlace(contactID, v.RelatedContactIDs)
	}
	if slice.Contains(v.RelatedContactIDs, contactID) {
		v.RelatedContactIDs = slice.RemoveInPlace(contactID, v.RelatedContactIDs)
		updates = append(updates, dal.Update{
			Field: "relatedContactIDs",
			Value: v.RelatedContactIDs,
		})
	}
	return updates
}

// Validate returns error if not valid
func (v *WithRelatedContacts) Validate() error {
	for contactID, relatedContact := range v.RelatedContacts {
		if contactID == "" {
			return validation.NewErrBadRecordFieldValue("relatedContacts", "has empty key for contact ID")
		}
		if !slice.Contains(v.RelatedContactIDs, contactID) {
			return validation.NewErrBadRecordFieldValue("relatedContacts."+contactID, "does not have relevant value in 'relatedContactIDs' field")
		}
		var field = func() string {
			return "relatedContacts." + contactID
		}
		if relatedContact == nil {
			return validation.NewErrRecordIsMissingRequiredField(field())
		}
		if err := relatedContact.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(field(), err.Error())
		}
	}
	for i, contactID := range v.RelatedContactIDs {
		if strings.TrimSpace(contactID) == "" {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedContactIDs[%d]", i), "empty contact ID")
		}
		if _, ok := v.RelatedContacts[contactID]; !ok {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedContactIDs[%d]", i), "does not have relevant value in 'relatedContacts' field")
		}
	}
	return nil
}
