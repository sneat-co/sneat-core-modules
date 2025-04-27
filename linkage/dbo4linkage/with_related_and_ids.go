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

const NoRelatedID = "-"
const AnyRelatedID = "*"

// WithRelatedAndIDs defines relationships of the current contact record to other contacts.
type WithRelatedAndIDs struct {
	WithRelated
	WithRelatedIDs

	//	Example of related field as a JSON and relevant relatedIDs field:
	/*
	   DebtusSpaceContactEntry(id="child1") {
	   	relatedIDs: ["space1:parent1:contactus:contacts:parent"],
	   	related: {
	   		"space1": { // SpaceID ContactID
	   			"contactus": { // Module ContactID
	   				"contacts": { // Collection
	   					"parent1": { // Item ContactID
	   						relatedAs: {
	   							"parent": {} // RelationshipRole ContactID
	   						}
	   						relatesAs: {
	   							"child": {} // RelationshipRole ContactID
	   						},
	   					},
	   				}
	   			},
	   		},
	   	}
	   }
	*/
}

type WithRelatedIDs struct {
	// RelatedIDs holds identifiers of related records - needed for indexed search.
	RelatedIDs []string `json:"relatedIDs,omitempty" firestore:"relatedIDs,omitempty"`
}

func (v *WithRelatedIDs) Validate() error {
	for i, relatedID := range v.RelatedIDs {
		s := strings.TrimSpace(relatedID)
		if s == "" {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), "empty contact ContactID")
		}
		if s != relatedID {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), "has leading or trailing spaces")
		}
	}
	return nil
}

func ValidateRelatedAndRelatedIDs(withRelated WithRelated, relatedIDs []string) error {
	if err := withRelated.ValidateRelated(func(itemKey ItemRef) error {
		// needs space ID to verify this
		//if id := itemKey.ID(); !slices.Contains(relatedIDs, id) {
		//	return validation.NewErrBadRecordFieldValue("relatedIDs",
		//		fmt.Sprintf(`does not have relevant value in 'relatedIDs' field: relatedID="%s", relatedIDs=[%s]`, id, strings.Join(relatedIDs, ",")))
		//}
		if id := itemKey.ModuleCollectionID(); !slices.Contains(relatedIDs, id) {
			return validation.NewErrBadRecordFieldValue("relatedIDs",
				fmt.Sprintf(`does not have relevant value in 'relatedIDs' field: relatedID="%s", relatedIDs=[%s]`, id, strings.Join(relatedIDs, ",")))
		}
		if id := itemKey.ModuleID(); !slices.Contains(relatedIDs, id) {
			return validation.NewErrBadRecordFieldValue("relatedIDs",
				fmt.Sprintf(`does not have relevant value in 'relatedIDs' field: relatedID="%s", relatedIDs=[%s]`, id, strings.Join(relatedIDs, ",")))
		}
		return nil
	}); err != nil {
		return err
	}
	if len(withRelated.Related) == 0 && len(relatedIDs) == 0 {
		return validation.NewErrBadRecordFieldValue("relatedIDs", "record without related items should have relatedIDs=['-']")
	}
	if len(withRelated.Related) > 0 && len(relatedIDs) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("relatedIDs")
	}
	if relatedIDs[0] != AnyRelatedID && relatedIDs[0] != NoRelatedID {
		return validation.NewErrBadRecordFieldValue("relatedIDs[0]", fmt.Sprintf("first value of relatedIDs should be either '%s' or '%s'", AnyRelatedID, NoRelatedID))
	}
	for i, relatedID := range relatedIDs[1:] { // The first item is always either "*" or "-"
		if strings.TrimSpace(relatedID) == "" {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), "empty string value")
		}
		if strings.HasSuffix(relatedID, "."+AnyRelatedID) {
			// TODO: Validate search index values
			continue
		}
		delimitersCount := strings.Count(relatedID, "&")
		switch delimitersCount {
		case 0:
			if relatedID != "*" {
				parts := strings.Split(relatedID, "=")
				if len(parts) != 2 {
					return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), "should be in format '{key}={value}', got: "+relatedID)
				}
				switch parts[0] {
				case "m", "s": // Module or space
					if err := with.ValidateRecordID(parts[1]); err != nil {
						return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), err.Error())
					}
				default:
					return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), "single key should start with 'm=' or 's=', got: "+relatedID)
				}
			}
		case 2: // "m={module}&c={collection}&s={space}"
			params := strings.Split(relatedID, "&")
			if !strings.HasPrefix(params[0], "m=") {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), "1st part of a 3 part ID should start with 'm=', got: "+params[0])
			}
			if !strings.HasPrefix(params[1], "c=") {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), "2nd part of a 3 part ID should start with 'c=', got: "+params[0])
			}
			if !strings.HasPrefix(params[2], "s=") {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("relatedIDs[%d]", i), "3d part of a 3 part ID should start with 's=', got: "+params[0])
			}
		case 3: // "m={module}&c={collection}&s={space}&i=item"
			relatedRef, err := NewItemRefFromString(relatedID)
			if err != nil {
				return err
			}

			relatedCollections := withRelated.Related[string(relatedRef.Module)]
			if relatedCollections == nil {
				return validation.NewErrBadRecordFieldValue(
					fmt.Sprintf("relatedIDs[%d]", i),
					fmt.Sprintf("field 'related[%s]' does not have module value", relatedRef.Module))
			}
			relatedItems := relatedCollections[relatedRef.Collection]
			if relatedItems == nil {
				return validation.NewErrBadRecordFieldValue(
					fmt.Sprintf("relatedIDs[%d]", i),
					fmt.Sprintf("field 'related[%s][%s]' does not have collection value", relatedRef.Module, relatedRef.Collection))
			}

			if _, ok := relatedItems[relatedRef.ItemID]; !ok {
				itemID := relatedRef.ItemID[:strings.Index(relatedRef.ItemID, SpaceItemIDSeparator)]
				if _, ok = relatedItems[itemID]; !ok {
					return validation.NewErrBadRecordFieldValue(
						fmt.Sprintf("relatedIDs[%d]", i),
						fmt.Sprintf("field 'related[%s][%s]' does not have values for either '%s' or '%s'",
							relatedRef.Module, relatedRef.Collection, relatedRef.ItemID, itemID))
				}
			}
		}
	}
	return nil
}

// Validate returns error if not valid
func (v *WithRelatedAndIDs) Validate() error {
	if err := v.WithRelatedIDs.Validate(); err != nil {
		return err
	}
	return ValidateRelatedAndRelatedIDs(v.WithRelated, v.RelatedIDs)
}

/*func (v *WithRelatedAndIDs) AddRelationshipsAndIDs(
	itemRef ItemRef,
	rolesOfItem RelationshipRoles,
	rolesToItem RelationshipRoles, // TODO: needs implementation
) (
	updates []update.Update, err error,
) {
	var command RelationshipItemRolesCommand
	if len(rolesOfItem) > 0 || len(rolesToItem) > 0 {
		if command.Add == nil {
			command.Add = new(RolesCommand)
		}
	}
	for roleOfItem := range rolesOfItem {
		command.Add.RolesOfItem = append(command.Add.RolesOfItem, roleOfItem)
	}
	for roleToItem := range rolesToItem {
		command.Add.RolesToItem = append(command.Add.RolesToItem, roleToItem)
	}
	return v.AddRelationshipAndID(itemRef, command)
	//return nil, errors.New("not implemented yet - AddRelationshipsAndIDs")
}*/

func UpdateRelatedIDs(spaceID coretypes.SpaceID, withRelated *WithRelated, withRelatedIDs *WithRelatedIDs) (updates []update.Update) {
	searchIndex := []string{AnyRelatedID}
	var spaceIDs []string
	var moduleCollectionSpaces []string
	for moduleID, relatedCollections := range withRelated.Related {
		searchIndex = append(searchIndex, "m="+string(moduleID))
		for collectionID, relatedItems := range relatedCollections {
			searchIndex = append(searchIndex, fmt.Sprintf("m=%s&c=%s", moduleID, collectionID))
			for itemID := range relatedItems {
				var itemSpaceID string
				if i := strings.Index(itemID, SpaceItemIDSeparator); i > 0 {
					itemSpaceID = itemID[i+1:]
					itemID = itemID[:i]
				} else {
					itemSpaceID = string(spaceID)
				}
				searchIndex = append(searchIndex, fmt.Sprintf("m=%s&c=%s&s=%s&i=%s", moduleID, collectionID, itemSpaceID, itemID))
				if !slices.Contains(spaceIDs, itemSpaceID) {
					searchIndex = append(searchIndex, fmt.Sprintf("s=%s", itemSpaceID))
					spaceIDs = append(spaceIDs, itemSpaceID)
				}
				if v := fmt.Sprintf("m=%s&c=%s&s=%s", moduleID, collectionID, itemSpaceID); !slices.Contains(moduleCollectionSpaces, v) {
					searchIndex = append(searchIndex, v)
					moduleCollectionSpaces = append(moduleCollectionSpaces, v)
				}
			}
		}
	}
	if len(searchIndex) > 2 {
		withRelatedIDs.RelatedIDs = searchIndex
		slices.Sort(searchIndex)
		updates = append(updates, update.ByFieldName("relatedIDs", withRelatedIDs.RelatedIDs))
	} else if len(withRelatedIDs.RelatedIDs) != 1 || withRelatedIDs.RelatedIDs[0] != NoRelatedID {
		withRelatedIDs.RelatedIDs = []string{NoRelatedID}
		updates = append(updates, update.ByFieldName("relatedIDs", withRelatedIDs.RelatedIDs))
	}
	return
}

func (v *WithRelatedAndIDs) AddRelationshipAndID(
	now time.Time,
	userID string,
	spaceID coretypes.SpaceID,
	command RelationshipItemRolesCommand,
) (
	updates []update.Update, err error,
) {
	updates, err = v.AddRelationship(now, userID, command)
	updates = append(updates, UpdateRelatedIDs(spaceID, &v.WithRelated, &v.WithRelatedIDs)...)
	return
}

func AddRelationshipAndID(
	now time.Time,
	userID string,
	spaceID coretypes.SpaceID,
	withRelated *WithRelated,
	withRelatedIDs *WithRelatedIDs,
	command RelationshipItemRolesCommand,
) (updates []update.Update, err error) {
	updates, err = withRelated.AddRelationship(now, userID, command)
	updates = append(updates, UpdateRelatedIDs(spaceID, withRelated, withRelatedIDs)...)
	return
}

func RemoveRelatedAndID(spaceID coretypes.SpaceID, withRelated *WithRelated, withRelatedIDs *WithRelatedIDs, ref ItemRef) (updates []update.Update) {
	updates = withRelated.RemoveRelatedItem(ref)
	updates = append(updates, UpdateRelatedIDs(spaceID, withRelated, withRelatedIDs)...)
	return updates
}

const (
	RelationshipRoleSpouse    = "spouse"
	RelationshipRoleParent    = "parent"
	RelationshipRoleChild     = "child"
	RelationshipRoleCousin    = "cousin"
	RelationshipRoleSibling   = "sibling"
	RelationshipRolePartner   = "partner"
	RelationshipRoleSpacemate = "space-mate"
)

const (
	RelationshipRoleManager      = "manager"
	RelationshipRoleDirectReport = "direct_report"
)

// Should provide a way for modules to register opposite roles?
var oppositeRoles = map[RelationshipRoleID]RelationshipRoleID{
	RelationshipRoleParent:       RelationshipRoleChild,
	RelationshipRoleChild:        RelationshipRoleParent,
	RelationshipRoleManager:      RelationshipRoleDirectReport,
	RelationshipRoleDirectReport: RelationshipRoleManager,
}

// Should provide a way for modules to register reciprocal roles?
var reciprocalRoles = []string{
	RelationshipRoleSpouse,
	RelationshipRoleSibling,
	RelationshipRoleCousin,
	RelationshipRolePartner,
	RelationshipRoleSpacemate,
}

func IsReciprocalRole(role RelationshipRoleID) bool {
	return slices.Contains(reciprocalRoles, role)
}

// GetOppositeRole returns relationship ContactID for the opposite direction
func GetOppositeRole(relationshipRoleID RelationshipRoleID) RelationshipRoleID {
	if relationshipRoleID == "" {
		return ""
	}
	if IsReciprocalRole(relationshipRoleID) {
		return relationshipRoleID
	}
	return oppositeRoles[relationshipRoleID]
}
