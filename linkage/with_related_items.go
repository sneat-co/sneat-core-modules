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

type RelatedItemsByID = map[string]*RelatedItem
type RelatedItemsByCollection = map[string]RelatedItemsByID
type RelatedItemsByModule = map[string]RelatedItemsByCollection
type RelatedItemsByTeam = map[string]RelatedItemsByModule

const relatedItemsField = "relatedItems"

// WithRelatedItems defines relationship of the current contact record to other contacts.
type WithRelatedItems struct {
	/* Example of relatedItems field as a JSON:

	Contact(id="child1") {
		relatedItemIDs: ["team1:parent1:contactus:contacts:parent"],
		relatedItems: {
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

	// RelatedItems defines relationship of the current contact to other contacts. Key is contact ID.
	RelatedItems RelatedItemsByTeam `json:"relatedItems,omitempty" firestore:"relatedItems,omitempty"`

	// RelatedItemIDs is a list of contact IDs that are related to the current contact - needed for indexed search.
	RelatedItemIDs []string `json:"relatedItemIDs,omitempty" firestore:"relatedItemIDs,omitempty"`
}

// Validate returns error if not valid
func (v *WithRelatedItems) Validate() error {
	for teamID, relatedItemsByModule := range v.RelatedItems {
		if teamID == "" {
			return validation.NewErrBadRecordFieldValue(relatedItemsField, "has empty team ID")
		}
		for moduleID, relatedItemsByCollectionID := range relatedItemsByModule {
			if moduleID == "" {
				return validation.NewErrBadRecordFieldValue(
					relatedItemsField+"."+teamID,
					"has empty module ID")
			}
			for collectionID, relatedItemsByID := range relatedItemsByCollectionID {
				if collectionID == "" {
					return validation.NewErrBadRecordFieldValue(
						fmt.Sprintf("%s.%s.%s", relatedItemsField, teamID, moduleID),
						"has empty collection ID",
					)
				}
				for itemID, relatedItem := range relatedItemsByID {
					if itemID == "" {
						return validation.NewErrBadRecordFieldValue(
							fmt.Sprintf("%s.%s.%s.%s", relatedItemsField, teamID, moduleID, collectionID),
							"has empty item ID")
					}

					key := fmt.Sprintf("%s.%s.%s.%s", teamID, moduleID, collectionID, itemID)
					field := relatedItemsField + "." + key

					if relatedItem == nil {
						return validation.NewErrRecordIsMissingRequiredField(field)
					}

					if err := relatedItem.Validate(); err != nil {
						return validation.NewErrBadRecordFieldValue(field, err.Error())
					}

					if !slice.Contains(v.RelatedItemIDs, key) {
						return validation.NewErrBadRecordFieldValue("relatedItemIDs",
							"does not have relevant value in 'relatedItemIDs' field: "+key)
					}
				}
			}
		}
	}
	for i, contactID := range v.RelatedItemIDs {
		if strings.TrimSpace(contactID) == "" {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedItemIDs[%d]", i), "empty contact ID")
		}
		if _, ok := v.RelatedItems[contactID]; !ok {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedItemIDs[%d]", i), "does not have relevant value in 'relatedItems' field")
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

func (v *WithRelatedItems) SetRelationshipToItem(
	userID string,
	recordRef TeamModuleDocRef,
	link Link,
	now time.Time,
) (updates []dal.Update, err error) {
	if err = link.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate link: %w", err)
	}

	var alreadyHasRelatedAs, alreadyHasRelatesAs bool

	if relatedItemsByModule := v.RelatedItems[link.TeamID]; relatedItemsByModule != nil {
		if relatedItemsByCollection := relatedItemsByModule[link.ModuleID]; relatedItemsByCollection != nil {
			if relatedItemsByID := relatedItemsByCollection[const4contactus.ContactsCollection]; relatedItemsByID != nil {
				if relatedItem := relatedItemsByID[link.ItemID]; relatedItem != nil {
					addIfNeeded := func(f string, itemRelationships Relationships, linkRelationshipIDs []RelationshipID) {
						field := func() string {
							return fmt.Sprintf("%s.%s.%s", relatedItemsField, link.ID(), f)
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
					addIfNeeded("relatedAs", relatedItem.RelatedAs, link.RelatedAs)
					addIfNeeded("relatesAs", relatedItem.RelatesAs, link.RelatesAs)
				}
			}
		}
	}

	if alreadyHasRelatedAs && alreadyHasRelatesAs {
		if len(v.RelatedItems) == 0 {
			v.RelatedItems = nil
		}
		return updates, nil
	}

	var relationshipUpdate []dal.Update
	if relationshipUpdate, err = v.AddRelationship(userID, recordRef, link, now); err != nil {
		return updates, err
	}
	updates = append(updates, relationshipUpdate...)
	updates = append(updates, dal.Update{Field: "relatedItemIDs", Value: v.RelatedItemIDs})

	return updates, err
}

func (v *WithRelatedItems) AddRelationship(
	userID string,
	recordRef TeamModuleDocRef,
	link Link,
	now time.Time,
) (updates []dal.Update, err error) {
	if err := recordRef.Validate(); err != nil {
		return nil, err
	}
	if v.RelatedItems == nil {
		v.RelatedItems = make(RelatedItemsByTeam, 1)
	}

	for _, linkRelatedAs := range link.RelatesAs {
		if relatesAs := GetRelatesAsFromRelated(linkRelatedAs); relatesAs != "" && !slice.Contains(link.RelatesAs, relatesAs) {
			link.RelatesAs = append(link.RelatesAs, "child")
		}
	}

	relatedItemsByModule := v.RelatedItems[link.TeamID]
	if relatedItemsByModule == nil {
		relatedItemsByModule = make(RelatedItemsByModule, 1)
		v.RelatedItems[link.TeamID] = relatedItemsByModule
	}

	relatedItemsByCollection := relatedItemsByModule[link.ModuleID]
	if relatedItemsByCollection == nil {
		relatedItemsByCollection = make(RelatedItemsByCollection, 1)
		relatedItemsByModule[const4contactus.ModuleID] = relatedItemsByCollection
	}

	relatedItemsByID := relatedItemsByCollection[const4contactus.ContactsCollection]
	if relatedItemsByID == nil {
		relatedItemsByID = make(RelatedItemsByID, 1)
		relatedItemsByCollection[const4contactus.ContactsCollection] = relatedItemsByID
	}

	relatedItem := relatedItemsByID[link.ItemID]
	if relatedItem == nil {
		relatedItem = NewRelatedItem()
		relatedItemsByID[link.ItemID] = relatedItem
	}

	relatedItemID := link.TeamModuleDocRef.ID()
	if !slice.Contains(v.RelatedItemIDs, relatedItemID) {
		v.RelatedItemIDs = append(v.RelatedItemIDs, relatedItemID)
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
					Field: fmt.Sprintf("relatedItems.%s.%s.%s", link.ID(), field, relationshipID),
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
func (v *WithRelatedItems) RemoveRelationshipToContact(contactID string) (updates []dal.Update) {
	if _, ok := v.RelatedItems[contactID]; ok {
		delete(v.RelatedItems, contactID)
		updates = append(updates, dal.Update{
			Field: relatedItemsField + "." + contactID,
			Value: dal.DeleteField,
		})
		v.RelatedItemIDs = slice.RemoveInPlace(contactID, v.RelatedItemIDs)
	}
	if slice.Contains(v.RelatedItemIDs, contactID) {
		v.RelatedItemIDs = slice.RemoveInPlace(contactID, v.RelatedItemIDs)
		updates = append(updates, dal.Update{
			Field: "relatedItemIDs",
			Value: v.RelatedItemIDs,
		})
	}
	return updates
}
