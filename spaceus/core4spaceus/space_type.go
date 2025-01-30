package core4spaceus

import (
	"fmt"
	"strings"
)

type SpaceType string

const (
	// SpaceTypePrivate is a "private" space type
	SpaceTypePrivate SpaceType = "private"

	// SpaceTypeFamily is a "family" space type
	SpaceTypeFamily SpaceType = "family"

	// SpaceTypeCompany is a "company" space type
	SpaceTypeCompany SpaceType = "company"

	// SpaceTypeSpace is a "space" space type
	SpaceTypeSpace SpaceType = "space"

	// SpaceTypeClub is a "club" space type
	SpaceTypeClub SpaceType = "club"
)

type SpaceRef string

func (v SpaceRef) SpaceType() SpaceType {
	if i := strings.Index(string(v), SpaceRefSeparator); i >= 0 {
		return SpaceType(v[:i])
	}
	if IsValidSpaceType(SpaceType(v)) {
		return SpaceType(v)
	}
	return ""
}

func (v SpaceRef) SpaceID() string {
	if i := strings.Index(string(v), SpaceRefSeparator); i >= 0 {
		return string(v[i+1:])
	}
	if !IsValidSpaceType(SpaceType(v)) {
		return string(v)
	}
	return ""
}

func (v SpaceRef) UrlPath() string {
	return fmt.Sprintf("%s/%s", v.SpaceType(), v.SpaceID())
}

const SpaceRefSeparator = "!"

func NewSpaceRef(spaceType SpaceType, spaceID string) SpaceRef {
	if !IsValidSpaceType(spaceType) {
		panic(fmt.Errorf("invalid space type: %v", spaceType))
	}
	if spaceID == "" {
		panic("spaceID is an empty string")
	}
	return SpaceRef(string(spaceType) + SpaceRefSeparator + spaceID)
}

func NewWeakSpaceRef(spaceType SpaceType) SpaceRef {
	return SpaceRef(spaceType)
}

// IsValidSpaceType checks if space has a valid/known type
func IsValidSpaceType(v SpaceType) bool {
	switch v {
	case SpaceTypeFamily, SpaceTypePrivate, SpaceTypeCompany, SpaceTypeSpace, SpaceTypeClub:
		return true
	default:
		return false
	}
}
