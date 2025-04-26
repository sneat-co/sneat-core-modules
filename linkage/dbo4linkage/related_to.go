package dbo4linkage

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
	"slices"
	"strings"
)

type ShortSpaceModuleDocRef struct {
	ID      string            `json:"id" firestore:"id"`
	SpaceID coretypes.SpaceID `json:"spaceID,omitempty" firestore:"spaceID,omitempty"`
}

func (v *ShortSpaceModuleDocRef) Validate() error {
	// Space can be empty for global collections like Happening
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("itemID")
	} else if err := validate.RecordID(v.ID); err != nil {
		return validation.NewErrBadRecordFieldValue("itemID", err.Error())
	}
	return nil
}

type SpaceModuleItemRef struct { // TODO: Move to sneat-go-core or document why not
	//Space      coretypes.SpaceID  `json:"space" firestore:"space"`
	Module     coretypes.ModuleID `json:"module" firestore:"module"`
	Collection string             `json:"collection" firestore:"collection"`
	ItemID     string             `json:"itemID" firestore:"itemID"`
}

func NewSpaceModuleItemRef(spaceID coretypes.SpaceID, module coretypes.ModuleID, collection, itemID string) SpaceModuleItemRef {
	if spaceID == "" {
		panic("spaceID is required")
	}
	if module == "" {
		panic("module is required")
	}
	if collection == "" {
		panic("collection is required")
	}
	if itemID == "" {
		panic("itemID is required")
	}
	return SpaceModuleItemRef{
		//Space:      spaceID,
		Module:     module,
		Collection: collection,
		ItemID:     itemID,
	}
}

func NewSpaceModuleItemRefFromString(id string) (itemRef SpaceModuleItemRef, err error) {
	ids := strings.Split(id, "&")
	if len(ids) != 4 {
		panic(fmt.Sprintf("invalid ContactID: '%s'", id))
	}
	for i, s := range ids {
		if s[1] != '=' {
			err = fmt.Errorf("expected to have '=' as 2nd charcter of value #%d, got '%s'", i, string(s[1]))
			return
		}
		switch s[0] {
		case 'm':
			itemRef.Module = coretypes.ModuleID(s[2:])
		case 'c':
			itemRef.Collection = s[2:]
		case 'i':
			itemRef.ItemID = s[2:]
		default:
			err = fmt.Errorf("unexpected key for value #%d - '%s'", i, id)
			return
		}
	}
	return
}

func (v SpaceModuleItemRef) ID() string {
	// The order is important for RelatedIDs field
	return fmt.Sprintf("m=%s&c=%s&i=%s", v.Module, v.Collection, v.ItemID)
}

func (v SpaceModuleItemRef) String() string {
	return fmt.Sprintf("{Module=%s. Collection=%s, ItemID=%s}", v.Module, v.Collection, v.ItemID)
}

func (v SpaceModuleItemRef) ModuleCollectionPath() string {
	return fmt.Sprintf("%s.%s", v.Module, v.Collection)
}

func (v SpaceModuleItemRef) Validate() error {
	// Space can be empty for global collections like Happening
	if v.Module == "" {
		return validation.NewErrRecordIsMissingRequiredField("moduleID")
	}
	if v.Collection == "" {
		return validation.NewErrRecordIsMissingRequiredField("collection")
	}
	if v.ItemID == "" {
		return validation.NewErrRecordIsMissingRequiredField("itemID")
	} else if err := validate.RecordID(v.ItemID); err != nil {
		return validation.NewErrBadRecordFieldValue("itemID", err.Error())
	}
	return nil
}

type RolesCommand struct {
	RolesOfItem []RelationshipRoleID `json:"rolesOfItem,omitempty" firestore:"rolesOfItem,omitempty"`
	RolesToItem []RelationshipRoleID `json:"rolesToItem,omitempty" firestore:"rolesToItem,omitempty"`
}

type RelationshipItemRolesCommand struct {
	ItemRef SpaceModuleItemRef `json:"itemRef"`
	Add     *RolesCommand      `json:"add,omitempty"`
	Remove  *RolesCommand      `json:"remove,omitempty"`
}

func (v RelationshipItemRolesCommand) Validate() error {
	//if err := v.SpaceModuleItemRef.Validate(); err != nil {
	//	return err
	//}
	if err := v.Add.Validate(); err != nil {
		return validation.NewErrBadRequestFieldValue("add", err.Error())
	}
	if err := v.Remove.Validate(); err != nil {
		return validation.NewErrBadRequestFieldValue("remove", err.Error())
	}
	return nil
}
func (v *RolesCommand) Validate() error {
	if v == nil {
		return nil
	}
	validateRelationIDs := func(field string, relations []string) error {
		for i, s := range relations {
			if strings.TrimSpace(s) != s {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("%s[%d]", field, i),
					"must not have leading or trailing spaces")
			}
			if slices.Contains(relations[:i], s) {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("%s[%d]", field, i),
					"duplicate relationship role value: "+s)
			}
		}
		return nil
	}
	if v.RolesToItem == nil && v.RolesOfItem == nil {
		return validation.NewErrRecordIsMissingRequiredField("rolesOfItem|rolesToItem")
	}
	if v.RolesToItem != nil {
		if err := validateRelationIDs("rolesOfItem", v.RolesOfItem); err != nil {
			return err
		}
	}
	if v.RolesToItem != nil {
		if err := validateRelationIDs("rolesToItem", v.RolesToItem); err != nil {
			return err
		}
	}
	return nil
}
