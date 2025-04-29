package dbo4linkage

import (
	"fmt"
	"github.com/dal-go/dalgo/update"
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

func GetRelatedItemByRef(relatedModules RelatedModules, itemRef ItemRef, createIfMissing bool) *RelatedItem {
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
func (v *WithRelated) RemoveRelatedItem(ref ItemRef) (updates []update.Update) {
	relatedCollections := v.Related[string(ref.Module)]
	if relatedCollections == nil {
		return
	}
	relatedItems := relatedCollections[ref.Collection]
	deletePath := []string{
		relatedField,
		string(ref.Module),
		ref.Collection,
		ref.ItemID,
	}
	if len(relatedItems) == 0 {
		return
	}
	if _, ok := relatedItems[ref.ItemID]; ok {
		delete(relatedItems, ref.ItemID)
		if len(relatedItems) == 0 {
			delete(relatedCollections, ref.Collection)
			deletePath = deletePath[:len(deletePath)-1]
			if len(relatedCollections) == 0 {
				delete(v.Related, string(ref.Module))
				deletePath = deletePath[:len(deletePath)-1]
				if len(v.Related) == 0 {
					v.Related = nil
					deletePath = []string{relatedField}
				}
			}
		}
		return []update.Update{
			update.ByFieldPath(deletePath, update.DeleteField),
		}
	}
	return
}

func (v *WithRelated) ValidateRelated(validateID func(itemKey ItemRef) error) error {
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
					if err := validateID(ItemRef{
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

func (v *WithRelated) removeRolesFromRelatedItem(itemRef ItemRef, remove RolesCommand) (updates []update.Update) {
	if len(remove.RolesOfItem) == 0 || len(remove.RolesToItem) == 0 {
		return
	}
	relatedCollections := v.Related[string(itemRef.Module)]
	if relatedCollections == nil {
		return
	}
	relatedItems := relatedCollections[itemRef.Collection]
	if relatedItems == nil {
		return
	}
	relatedItem := relatedItems[itemRef.ItemID]
	if relatedItem == nil {
		return
	}

	for _, role := range remove.RolesOfItem {
		if oppositeRole := GetOppositeRole(role); oppositeRole != "" {
			if slices.Contains(remove.RolesToItem, oppositeRole) {
				remove.RolesToItem = append(remove.RolesToItem, oppositeRole)
			}
		}
	}
	for _, role := range remove.RolesToItem {
		if oppositeRole := GetOppositeRole(role); oppositeRole != "" {
			if slices.Contains(remove.RolesOfItem, oppositeRole) {
				remove.RolesOfItem = append(remove.RolesOfItem, oppositeRole)
			}
		}
	}

	removeRoles := func(field string, roles RelationshipRoles, rolesToRemove []RelationshipRoleID) (updates []update.Update) {
		var roleUpdates []update.Update
		for _, role := range rolesToRemove {
			if roles[role] != nil {
				delete(roles, role)
				roleUpdates = append(roleUpdates, update.ByFieldPath([]string{
					relatedField,
					string(itemRef.Module),
					itemRef.Collection,
					itemRef.ItemID,
					field,
					role,
				}, update.DeleteField))
			}
		}
		if len(roles) > 0 {
			updates = roleUpdates
		}
		return
	}
	if len(relatedItem.RolesOfItem) > 0 || len(relatedItem.RolesToItem) > 0 {
		itemUpdates := removeRoles("rolesOfItem", relatedItem.RolesOfItem, remove.RolesOfItem)
		itemUpdates = append(itemUpdates, removeRoles("rolesToItem", relatedItem.RolesToItem, remove.RolesToItem)...)
		if len(relatedItem.RolesOfItem) == 0 && len(relatedItem.RolesToItem) == 0 {
			updates = append(updates, v.RemoveRelatedItem(itemRef)...)
		} else {
			updates = append(updates, itemUpdates...)
		}
	}
	return
}

func (v *WithRelated) addRolesToRelatedItem(itemRef ItemRef, add RolesCommand, userID string, now time.Time) (updates []update.Update) {
	if v.Related == nil {
		v.Related = make(RelatedModules, 1)
	}
	addOppositeRoles := func(roles []RelationshipRoleID, oppositeRoles []RelationshipRoleID) []RelationshipRoleID {
		for _, roleOfItem := range roles {
			if oppositeRole := GetOppositeRole(roleOfItem); oppositeRole != "" && !slices.Contains(add.RolesToItem, oppositeRole) {
				oppositeRoles = append(oppositeRoles, oppositeRole)
			}
		}
		return oppositeRoles
	}
	add.RolesToItem = addOppositeRoles(add.RolesOfItem, add.RolesToItem)
	add.RolesOfItem = addOppositeRoles(add.RolesToItem, add.RolesOfItem)

	relatedCollections := v.Related[string(itemRef.Module)]
	if relatedCollections == nil {
		relatedCollections = make(RelatedCollections, 1)
		v.Related[string(itemRef.Module)] = relatedCollections
	}

	relatedItems := relatedCollections[itemRef.Collection]
	if relatedItems == nil {
		relatedItems = make(RelatedItems, 1)
		relatedCollections[itemRef.Collection] = relatedItems
	}

	var relatedItemChanged bool
	relatedItem := relatedItems[itemRef.ItemID]
	if relatedItem == nil {
		relatedItem = NewRelatedItem()
		relatedItems[itemRef.ItemID] = relatedItem
		relatedItemChanged = true
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
				relatedItemChanged = true
			}
		}
		return relationships
	}
	relatedItem.RolesOfItem = addRelationship("rolesOfItem", add.RolesOfItem, relatedItem.RolesOfItem)
	relatedItem.RolesToItem = addRelationship("rolesToItem", add.RolesToItem, relatedItem.RolesToItem)

	if relatedItemChanged {
		updates = append(updates, update.ByFieldPath(
			[]string{relatedField, string(itemRef.Module), itemRef.Collection, itemRef.ItemID},
			relatedItem,
		))
	}
	return
}

func (v *WithRelated) ProcessRelatedCommand(
	now time.Time,
	userID string,
	command RelationshipItemRolesCommand,
) (
	updates []update.Update, err error,
) {
	if err = command.Validate(); err != nil {
		return nil, err
	}

	if command.Remove != nil {
		updates = append(updates, v.removeRolesFromRelatedItem(command.ItemRef, *command.Remove)...)
	}

	if command.Add != nil {
		updates = append(updates, v.addRolesToRelatedItem(command.ItemRef, *command.Add, userID, now)...)
	}

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
