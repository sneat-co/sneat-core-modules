package dbo4linkage

import (
	"fmt"
	"github.com/strongo/validation"
	"slices"
	"strings"
)

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
