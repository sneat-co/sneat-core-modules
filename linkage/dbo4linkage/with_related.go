package dbo4linkage

import (
	"fmt"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"slices"
	"strings"
	"time"
)

type RelationshipRoleID = string

type RelationshipRole struct {
	with.CreatedField
}

func (v RelationshipRole) Validate() error {
	return nil
	//return v.CreatedField.Validate()
}

type RelationshipRoles = map[RelationshipRoleID]*RelationshipRole

type RelatedItemKey struct {
	SpaceID coretypes.SpaceID `json:"spaceID" firestore:"spaceID"`
	ItemID  string            `json:"itemID" firestore:"itemID"`
}

func (v RelatedItemKey) String() string {
	return fmt.Sprintf("%s@%s", v.ItemID, v.SpaceID)
}

func (v RelatedItemKey) Validate() error {
	if v.SpaceID == "" {
		return validation.NewErrRecordIsMissingRequiredField("spaceID")
	}
	if v.ItemID == "" {
		return validation.NewErrRecordIsMissingRequiredField("itemID")
	}
	return nil
}

func GetRelatedItemByRef(relatedModules RelatedModules, itemRef SpaceModuleItemRef, createIfMissing bool) *RelatedItem {
	relatedCollections := relatedModules[string(itemRef.Module)]
	if !createIfMissing && len(relatedCollections) == 0 {
		return nil
	}
	relatedItems := relatedCollections[itemRef.Collection]
	if !createIfMissing && len(relatedItems) == 0 {
		return nil
	}
	relatedItem, exists := relatedItems[itemRef.ItemID]
	if exists {
		return relatedItem
	}
	if createIfMissing {
		relatedItem = NewRelatedItem()
		if relatedItems == nil {
			relatedItems = make(RelatedItems, 1)
		}
		relatedItems[itemRef.ItemID] = relatedItem
		if relatedCollections == nil {
			relatedCollections = make(RelatedCollections, 1)
		}
		relatedCollections[itemRef.Collection] = relatedItems
		if relatedModules == nil {
			relatedModules = make(RelatedModules, 1)
		}
		relatedModules[string(itemRef.Module)] = relatedCollections
		return relatedItem
	}
	return nil
}

type RelatedItem struct {
	//Keys []RelatedItemKey `json:"keys" firestore:"keys"` // TODO: document why we need multiple keys, provide a use case

	Note string `json:"note,omitempty" firestore:"note,omitempty"`

	// RolesOfItem - if related item is a child of the current record, then rolesOfItem = {"child": ...}
	RolesOfItem RelationshipRoles `json:"rolesOfItem,omitempty" firestore:"rolesOfItem,omitempty"`

	// RolesToItem - if related item is a child of the current contact, then rolesToItem = {"parent": ...}
	RolesToItem RelationshipRoles `json:"rolesToItem,omitempty" firestore:"rolesToItem,omitempty"`
}

func (v *RelatedItem) String() string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprintf("RelatedItem{RolesOfItem=%+v, RolesToItem=%+v}", v.RolesOfItem, v.RolesToItem)
}

func NewRelatedItem() *RelatedItem {
	return new(RelatedItem)
}

func (v *RelatedItem) Validate() error {
	if err := v.validateRelationships(v.RolesOfItem); err != nil {
		return validation.NewErrBadRecordFieldValue("rolesOfItem", err.Error())
	}
	if err := v.validateRelationships(v.RolesToItem); err != nil {
		return validation.NewErrBadRecordFieldValue("rolesToItem", err.Error())
	}
	return nil
}

func (*RelatedItem) validateRelationships(related RelationshipRoles) error {
	for relationshipID, relationshipDetails := range related {
		if strings.TrimSpace(relationshipID) == "" {
			return validation.NewValidationError("relationship key is empty string")
		}
		if err := relationshipDetails.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(relationshipID, err.Error())
		}
	}
	return nil
}

type RelatedItems = map[string]*RelatedItem
type RelatedCollections = map[string]RelatedItems
type RelatedModules = map[string]RelatedCollections

const relatedField = "related"

var _ Relatable = (*WithRelatedAndIDs)(nil)

func (v *WithRelatedAndIDs) GetRelated() *WithRelatedAndIDs {
	return v
}

type WithRelated struct {
	// Related defines relationships of the current contact to other contacts.
	// Key is space ContactID.
	Related RelatedModules `json:"related,omitempty" firestore:"related,omitempty"`
}

func (v *WithRelated) Validate() error {
	return v.ValidateRelated(nil)
}

// RemoveRelatedItem removes all relationships to a given item
// TODO(help-wanted): needs 100% code coverage by tests
func (v *WithRelated) RemoveRelatedItem(ref SpaceModuleItemRef) (updates []update.Update) {
	relatedCollections := v.Related[string(ref.Module)]
	if relatedCollections == nil {
		return
	}
	relatedItems := relatedCollections[ref.Collection]
	if len(relatedItems) == 0 {
		return
	}
	if _, ok := relatedItems[ref.ItemID]; ok {
		delete(relatedItems, ref.ItemID)
		if len(relatedItems) == 0 {
			delete(relatedCollections, ref.Collection)
			if len(v.Related) == 0 {
				delete(v.Related, string(ref.Module))
				updates = append(updates,
					update.ByFieldName(relatedField, update.DeleteField))
			} else {
				updates = append(updates,
					update.ByFieldPath([]string{relatedField, string(ref.Module)}, update.DeleteField))
			}
		} else {
			updates = append(updates,
				update.ByFieldPath([]string{relatedField, string(ref.Module), ref.Collection, ref.ItemID}, update.DeleteField))
		}
	}
	return updates
}

func (v *WithRelated) ValidateRelated(validateID func(itemKey SpaceModuleItemRef) error) error {
	for moduleID, relatedCollections := range v.Related {
		if moduleID == "" {
			return validation.NewErrBadRecordFieldValue(relatedField, "has an empty module key")
		}
		for collectionID, relatedItems := range relatedCollections {
			if collectionID == "" {
				return validation.NewErrBadRecordFieldValue(
					fmt.Sprintf("%s.%s", relatedField, moduleID),
					"has an empty collection key",
				)
			}
			for itemID, relatedItem := range relatedItems {
				switch itemID {
				case "":
					return validation.NewErrBadRecordFieldValue(
						fmt.Sprintf("%s.%s.%s", relatedField, moduleID, collectionID),
						"has an empty item key")
				case "itemID":
					return validation.NewErrBadRecordFieldValue(
						fmt.Sprintf("%s.%s.%s", relatedField, moduleID, collectionID),
						"item key should not be 'itemID'")
				}
				if err := relatedItem.Validate(); err != nil {
					return validation.NewErrBadRecordFieldValue(
						fmt.Sprintf("%s.%s.%s.%s", relatedField, moduleID, collectionID, itemID),
						err.Error())
				}
				if validateID != nil {
					if err := validateID(SpaceModuleItemRef{
						Module:     coretypes.ModuleID(moduleID),
						Collection: collectionID,
						ItemID:     itemID,
					}); err != nil {
						return fmt.Errorf("validateID(%s) returned error: %w", itemID, err)
					}
				}
			}
		}
	}
	return nil
}

func (v *WithRelated) AddRelationship(
	now time.Time,
	userID string,
	command RelationshipItemRolesCommand,
) (
	updates []update.Update, err error,
) {
	if err := command.Validate(); err != nil {
		return nil, err
	}
	if v.Related == nil {
		v.Related = make(RelatedModules, 1)
	}

	if command.Add != nil {
		addOppositeRoles := func(roles []RelationshipRoleID, oppositeRoles []RelationshipRoleID) []RelationshipRoleID {
			for _, roleOfItem := range roles {
				if oppositeRole := GetOppositeRole(roleOfItem); oppositeRole != "" && !slices.Contains(command.Add.RolesToItem, oppositeRole) {
					oppositeRoles = append(oppositeRoles, oppositeRole)
				}
			}
			return oppositeRoles
		}
		command.Add.RolesToItem = addOppositeRoles(command.Add.RolesOfItem, command.Add.RolesToItem)
		command.Add.RolesOfItem = addOppositeRoles(command.Add.RolesToItem, command.Add.RolesOfItem)
	}

	relatedByCollectionID := v.Related[string(command.ItemRef.Module)]
	if relatedByCollectionID == nil {
		relatedByCollectionID = make(RelatedCollections, 1)
		v.Related[string(command.ItemRef.Module)] = relatedByCollectionID
	}

	relatedItems := relatedByCollectionID[const4contactus.ContactsCollection]
	if relatedItems == nil {
		relatedItems = make(RelatedItems, 1)
		relatedByCollectionID[const4contactus.ContactsCollection] = relatedItems
	}

	relatedItem := relatedItems[command.ItemRef.ItemID]
	if relatedItem == nil {
		relatedItem = NewRelatedItem()
		relatedItems[command.ItemRef.ItemID] = relatedItem
		relatedByCollectionID[const4contactus.ContactsCollection] = relatedItems
	}

	addRelationship := func(field string, relationshipIDs []RelationshipRoleID, relationships RelationshipRoles) RelationshipRoles {
		if len(relationshipIDs) == 0 {
			return relationships
		}
		if relationships == nil {
			relationships = make(RelationshipRoles, len(relationshipIDs))
		}
		for _, relationshipID := range relationshipIDs {
			if relationship := relationships[relationshipID]; relationship == nil {
				relationship = &RelationshipRole{
					CreatedField: with.CreatedField{
						Created: with.Created{
							By: userID,
							At: now.Format(time.RFC3339),
						},
					},
				}
				relationships[relationshipID] = relationship
			}
		}
		return relationships
	}

	if command.Add != nil {
		relatedItem.RolesOfItem = addRelationship("rolesOfItem", command.Add.RolesOfItem, relatedItem.RolesOfItem)
		relatedItem.RolesToItem = addRelationship("rolesToItem", command.Add.RolesToItem, relatedItem.RolesToItem)
	}

	updates = append(updates, update.ByFieldName(
		fmt.Sprintf("related.%s", command.ItemRef.ModuleCollectionPath()),
		relatedItems))

	return updates, nil
}

//func (v *WithRelated) SetRelationshipToItem(
//	userID string,
//	command RelationshipItemRolesCommand,
//	now time.Time,
//) (updates []update.Update, err error) {
//	if err = command.Validate(); err != nil {
//		return nil, fmt.Errorf("failed to validate command: %w", err)
//	}
//
//	//var alreadyHasRelatedAs bool
//
//	changed := false
//
//	if v.Related == nil {
//		v.Related = make(RelatedModules, 1)
//	}
//	relatedByCollectionID := v.Related[command.Module]
//	if relatedByCollectionID == nil {
//		relatedByCollectionID = make(RelatedCollections, 1)
//		v.Related[command.Module] = relatedByCollectionID
//	}
//	relatedItems := relatedByCollectionID[const4contactus.ContactsCollection]
//	//if relatedItems == nil {
//	//	relatedItems = make([]*RelatedItem, 0, 1)
//	//	relatedByCollectionID[const4contactus.ContactsCollection] = relatedItems
//	//}
//	relatedItemKey := RelatedItemKey{SpaceID: command.SpaceID, ItemID: command.ItemID}
//	relatedItem := GetRelatedItemByKey(relatedItems, relatedItemKey)
//	if relatedItem == nil {
//		relatedItem = NewRelatedItem(relatedItemKey)
//		relatedItems = append(relatedItems, relatedItem)
//		relatedByCollectionID[const4contactus.ContactsCollection] = relatedItems
//		changed = true
//	}
//
//	//addIfNeeded := func(f string, itemRelationships RelationshipRoles, linkRelationshipIDs []RelationshipRoleID) {
//	//	field := func() string {
//	//		return fmt.Sprintf("%s.%s.%s", relatedField, command.ContactID(), f)
//	//	}
//	//	for _, linkRelationshipID := range linkRelationshipIDs {
//	//		itemRelationship := itemRelationships[linkRelationshipID]
//	//		if itemRelationship == nil {
//	//			itemRelationships[linkRelationshipID] = &RelationshipRole{
//	//				CreatedField: with.CreatedField{
//	//					Created: with.Created{
//	//						By: userID,
//	//						At: now.Format(time.DateOnly),
//	//					},
//	//				},
//	//			}
//	//			alreadyHasRelatedAs = true
//	//		}
//	//	}
//	//}
//	//addIfNeeded("rolesOfItem", relatedItem.RolesOfItem, command.RolesOfItem)
//	//addIfNeeded("rolesToItem", relatedItem.RolesToItem, command.RolesToItem)
//
//	var relationshipUpdate []update.Update
//	if relationshipUpdate, err = v.AddRelationshipAndID(userID, command, now); err != nil {
//		return updates, err
//	}
//	updates = append(updates, relationshipUpdate...)
//
//	return updates, err
//}
