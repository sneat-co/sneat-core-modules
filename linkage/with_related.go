package linkage

import (
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
	"strings"
	"time"
)

type RelationshipID = string

type Relationship struct {
	dbmodels.WithCreatedField
}

func (v Relationship) Validate() error {
	return v.WithCreatedField.Validate()
}

type Relationships = map[RelationshipID]*Relationship

type RelatedItem struct {
	// Brief any // TODO: do we need a brief here?

	// RelatedAs - if related contact is a child of the current contact, then relatedAs = {"child": ...}
	RelatedAs Relationships `json:"relatedAs,omitempty" firestore:"relatedAs,omitempty"`

	// RelatesAs - if related contact is a child of the current contact, then relatesAs = {"parent": ...}
	RelatesAs Relationships `json:"relatesAs,omitempty" firestore:"relatesAs,omitempty"`
}

func (v *RelatedItem) String() string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprintf("RelatedItem{RelatedAs=%+v, RelatesAs=%+v}", v.RelatedAs, v.RelatesAs)
}

func NewRelatedItem() *RelatedItem {
	return &RelatedItem{
		RelatedAs: make(Relationships, 1),
		RelatesAs: make(Relationships, 1),
	}
}

func (v *RelatedItem) Validate() error {
	if err := v.validateRelationships(v.RelatedAs); err != nil {
		return validation.NewErrBadRecordFieldValue("relatedAs", err.Error())
	}
	if err := v.validateRelationships(v.RelatesAs); err != nil {
		return validation.NewErrBadRecordFieldValue("relatesAs", err.Error())
	}
	return nil
}

func (*RelatedItem) validateRelationships(related Relationships) error {
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

type RelatedByItemID = map[string]*RelatedItem
type RelatedByCollectionID = map[string]RelatedByItemID
type RelatedByModuleID = map[string]RelatedByCollectionID
type RelatedByTeamID = map[string]RelatedByModuleID

const relatedField = "related"

// WithRelated defines relationship of the current contact record to other contacts.
type WithRelated struct {
	/* Example of related field as a JSON:

	Contact(id="child1") {
		relatedIDs: ["team1:parent1:contactus:contacts:parent"],
		related: {
			"team1": { // Team ID
				"contactus": { // Module ID
					"contacts": { // Collection
						"parent1": { // Item ID
							relatedAs: {
								"parent": {} // Relationship ID
							}
							relatesAs: {
								"child": {} // Relationship ID
							},
						},
					}
				},
			},
		}
	}
	*/

	// Related defines relationship of the current contact to other contacts. Key is contact ID.
	Related RelatedByTeamID `json:"related,omitempty" firestore:"related,omitempty"`

	// RelatedIDs is a list of IDs of records that are related to the current record - this is needed for indexed search.
	RelatedIDs []string `json:"relatedIDs,omitempty" firestore:"relatedIDs,omitempty"`
}

// Validate returns error if not valid
func (v *WithRelated) Validate() error {
	for teamID, relatedByModuleID := range v.Related {
		if teamID == "" {
			return validation.NewErrBadRecordFieldValue(relatedField, "has empty team ID")
		}
		for moduleID, relatedByCollectionID := range relatedByModuleID {
			if moduleID == "" {
				return validation.NewErrBadRecordFieldValue(
					relatedField+"."+teamID,
					"has empty module ID")
			}
			for collectionID, relatedByItemID := range relatedByCollectionID {
				if collectionID == "" {
					return validation.NewErrBadRecordFieldValue(
						fmt.Sprintf("%s.%s.%s", relatedField, teamID, moduleID),
						"has empty collection ID",
					)
				}
				for itemID, relatedItem := range relatedByItemID {
					if itemID == "" {
						return validation.NewErrBadRecordFieldValue(
							fmt.Sprintf("%s.%s.%s.%s", relatedField, teamID, moduleID, collectionID),
							"has empty item ID")
					}

					key := fmt.Sprintf("%s.%s.%s.%s", teamID, moduleID, collectionID, itemID)
					field := relatedField + "." + key

					if relatedItem == nil {
						return validation.NewErrRecordIsMissingRequiredField(field)
					}

					if err := relatedItem.Validate(); err != nil {
						return validation.NewErrBadRecordFieldValue(field, err.Error())
					}

					if !slice.Contains(v.RelatedIDs, key) {
						return validation.NewErrBadRecordFieldValue("relatedIDs",
							"does not have relevant value in 'relatedIDs' field: "+key)
					}
				}
			}
		}
	}
	for i, relatedID := range v.RelatedIDs {
		if strings.TrimSpace(relatedID) == "" {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), "empty contact ID")
		}
		relatedRef := NewTeamModuleDocRef(relatedID)

		relatedByModuleID := v.Related[relatedRef.TeamID]
		if relatedByModuleID == nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), fmt.Sprintf("field 'related'  does not have value for team ID=%s", relatedRef.TeamID))
		}
		relatedByCollectionID := relatedByModuleID[relatedRef.ModuleID]
		if relatedByCollectionID == nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), fmt.Sprintf("field 'related[%s]' does not have value for module ID=%s", relatedRef.TeamID, relatedRef.ModuleID))
		}
		relatedByItemID := relatedByCollectionID[relatedRef.Collection]
		if relatedByItemID == nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), fmt.Sprintf("field 'related[%s][%s]' does not have value for collection ID=%s", relatedRef.TeamID, relatedRef.ModuleID, relatedRef.Collection))
		}
		_, ok := relatedByItemID[relatedRef.ItemID]
		if !ok {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), fmt.Sprintf("field 'related[%s][%s][%s]' does not have value for item ID=%s", relatedRef.TeamID, relatedRef.ModuleID, relatedRef.Collection, relatedRef.ItemID))
		}
	}
	return nil
}

func GetRelatesAsFromRelated(relatedAs RelationshipID) RelationshipID {
	switch relatedAs {
	case "parent":
		return "child"
	case "spouse":
		return "spouse"
	}
	return ""
}

func (v *WithRelated) SetRelationshipToItem(
	userID string,
	recordRef TeamModuleDocRef,
	link Link,
	now time.Time,
) (updates []dal.Update, err error) {
	if err = link.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate link: %w", err)
	}

	var alreadyHasRelatedAs, alreadyHasRelatesAs bool

	if relatedByModuleID := v.Related[link.TeamID]; relatedByModuleID != nil {
		if relatedByCollectionID := relatedByModuleID[link.ModuleID]; relatedByCollectionID != nil {
			if relatedByItemID := relatedByCollectionID[const4contactus.ContactsCollection]; relatedByItemID != nil {
				if related := relatedByItemID[link.ItemID]; related != nil {
					addIfNeeded := func(f string, itemRelationships Relationships, linkRelationshipIDs []RelationshipID) {
						field := func() string {
							return fmt.Sprintf("%s.%s.%s", relatedField, link.ID(), f)
						}
						for _, linkRelationshipID := range linkRelationshipIDs {
							for itemRelationshipID := range itemRelationships {
								if itemRelationshipID == linkRelationshipID {
									alreadyHasRelatedAs = true
								} else {
									updates = append(updates, dal.Update{Field: field(), Value: dal.DeleteField})
								}
							}
						}
					}
					addIfNeeded("relatedAs", related.RelatedAs, link.RelatedAs)
					addIfNeeded("relatesAs", related.RelatesAs, link.RelatesAs)
				}
			}
		}
	}

	if alreadyHasRelatedAs && alreadyHasRelatesAs {
		if len(v.Related) == 0 {
			v.Related = nil
		}
		return updates, nil
	}

	var relationshipUpdate []dal.Update
	if relationshipUpdate, err = v.AddRelationship(userID, recordRef, link, now); err != nil {
		return updates, err
	}
	updates = append(updates, relationshipUpdate...)
	updates = append(updates, dal.Update{Field: "relatedIDs", Value: v.RelatedIDs})

	return updates, err
}

func (v *WithRelated) AddRelationship(
	userID string,
	recordRef TeamModuleDocRef,
	link Link,
	now time.Time,
) (updates []dal.Update, err error) {
	if err := recordRef.Validate(); err != nil {
		return nil, err
	}
	if v.Related == nil {
		v.Related = make(RelatedByTeamID, 1)
	}

	for _, linkRelatedAs := range link.RelatesAs {
		if relatesAs := GetRelatesAsFromRelated(linkRelatedAs); relatesAs != "" && !slice.Contains(link.RelatesAs, relatesAs) {
			link.RelatesAs = append(link.RelatesAs, "child")
		}
	}

	relatedByModuleID := v.Related[link.TeamID]
	if relatedByModuleID == nil {
		relatedByModuleID = make(RelatedByModuleID, 1)
		v.Related[link.TeamID] = relatedByModuleID
	}

	relatedByCollectionID := relatedByModuleID[link.ModuleID]
	if relatedByCollectionID == nil {
		relatedByCollectionID = make(RelatedByCollectionID, 1)
		relatedByModuleID[const4contactus.ModuleID] = relatedByCollectionID
	}

	relatedByItemID := relatedByCollectionID[const4contactus.ContactsCollection]
	if relatedByItemID == nil {
		relatedByItemID = make(RelatedByItemID, 1)
		relatedByCollectionID[const4contactus.ContactsCollection] = relatedByItemID
	}

	relatedItem := relatedByItemID[link.ItemID]
	if relatedItem == nil {
		relatedItem = NewRelatedItem()
		relatedByItemID[link.ItemID] = relatedItem
	}

	relatedItemID := link.TeamModuleDocRef.ID()
	if !slice.Contains(v.RelatedIDs, relatedItemID) {
		v.RelatedIDs = append(v.RelatedIDs, relatedItemID)
	}

	addRelationship := func(field string, relationshipIDs []RelationshipID, relationships Relationships) Relationships {
		if relationships == nil {
			relationships = make(Relationships, len(relationshipIDs))
		}
		for _, relationshipID := range relationshipIDs {
			if relationship := relationships[relationshipID]; relationship == nil {
				relationship = &Relationship{
					WithCreatedField: dbmodels.WithCreatedField{
						Created: dbmodels.Created{
							By: userID,
							On: now.Format(time.DateOnly),
						},
					},
				}
				relationships[relationshipID] = relationship
				updates = append(updates, dal.Update{
					Field: fmt.Sprintf("related.%s.%s.%s", link.ID(), field, relationshipID),
					Value: relationship,
				})
			}
		}
		return relationships
	}

	relatedItem.RelatedAs = addRelationship("relatedAs", link.RelatedAs, relatedItem.RelatedAs)
	relatedItem.RelatedAs = addRelationship("relatesAs", link.RelatesAs, relatedItem.RelatesAs)

	return updates, nil
}

// RemoveRelationshipToContact removes all relationships to a given contact
func (v *WithRelated) RemoveRelationshipToContact(contactID string) (updates []dal.Update) {
	if _, ok := v.Related[contactID]; ok {
		delete(v.Related, contactID)
		updates = append(updates, dal.Update{
			Field: relatedField + "." + contactID,
			Value: dal.DeleteField,
		})
		v.RelatedIDs = slice.RemoveInPlace(contactID, v.RelatedIDs)
	}
	if slice.Contains(v.RelatedIDs, contactID) {
		v.RelatedIDs = slice.RemoveInPlace(contactID, v.RelatedIDs)
		updates = append(updates, dal.Update{
			Field: "relatedIDs",
			Value: v.RelatedIDs,
		})
	}
	return updates
}
