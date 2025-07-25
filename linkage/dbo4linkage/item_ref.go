package dbo4linkage

import (
	"errors"
	"fmt"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
	"net/url"
	"strings"
)

type ItemRef struct { // TODO: Move to sneat-go-core or document why not
	ExtID      coretypes.ExtID `json:"module" firestore:"module"` // TODO: change to `json:"extID" firestore:"extID"`?
	Collection string          `json:"collection" firestore:"collection"`
	ItemID     string          `json:"itemID" firestore:"itemID"`
	//SpaceID    coretypes.SpaceID  `json:"spaceID,omitempty" firestore:"spaceID,omitempty"`
}

func NewItemRefSameSpace(extID coretypes.ExtID, collection, itemID string) ItemRef {
	if strings.Contains(itemID, "@") {
		panic("itemID must not contain a spaceID separated by '@'")
	}
	return newItemRef(extID, collection, itemID)
}

func newItemRef(extID coretypes.ExtID, collection, itemID string) ItemRef {
	if extID == "" {
		panic("extID is required")
	}
	if collection == "" {
		panic("collection is required")
	}
	if itemID == "" {
		panic("itemID is required")
	}
	return ItemRef{
		//SpaceID:    spaceID,
		ExtID:      extID,
		Collection: collection,
		ItemID:     itemID,
	}
}

func NewItemRefFromQueryString(values url.Values) (itemRef ItemRef, err error) {
	if itemRef.ExtID = coretypes.ExtID(values.Get("m")); strings.TrimSpace(string(itemRef.ExtID)) == "" {
		return itemRef, errors.New("extension ID 'm' parameter is required")
	}
	if itemRef.Collection = values.Get("c"); strings.TrimSpace(string(itemRef.Collection)) == "" {
		return itemRef, errors.New("collectionID 'c' parameter is required")
	}
	if itemRef.ItemID = values.Get("i"); strings.TrimSpace(itemRef.ItemID) == "" {
		return itemRef, errors.New("itemID 'i' parameter is required")
	}
	if spaceID := values.Get("s"); spaceID != "" {
		itemRef.ItemID = itemRef.ItemID + SpaceItemIDSeparator + spaceID
	}
	return
}

func (v ItemRef) ID() string {
	// The order is important for the RelatedIDs field
	return fmt.Sprintf("m=%s&c=%s&i=%s", v.ExtID, v.Collection, v.ItemID)
}

func (v ItemRef) String() string {
	return fmt.Sprintf("{ExtID=%s,Collection=%s,ItemID=%s}", v.ExtID, v.Collection, v.ItemID)
}

func (v ItemRef) ExtensionID() string {
	return "m=" + string(v.ExtID)
}

func (v ItemRef) ExtensionCollectionPath() string {
	// The order is important for the RelatedIDs field
	return fmt.Sprintf("%s.%s", v.ExtID, v.Collection)
}

func (v ItemRef) ExtensionCollectionID() string {
	return fmt.Sprintf("m=%s&c=%s", v.ExtID, v.Collection)
}

func (v ItemRef) Validate() error {
	// SpaceID can be empty for global collections like Happening
	if v.ExtID == "" {
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
