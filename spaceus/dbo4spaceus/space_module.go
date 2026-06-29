package dbo4spaceus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/core/coremodels"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

const SpaceModulesCollection = coremodels.ExtCollection

// specscore: decisions/0002-reserved-extension-space-ids
// Spaceless key-builder: when the space is empty, build /ext/{ext-id}/... and do
// NOT call NewSpaceKey (which panics on empty). See sneat-specs Decision 0002:
// https://github.com/sneat-co/sneat-specs/blob/main/spec/decisions/0002-reserved-extension-space-ids.md
func NewSpaceModuleKey(spaceID coretypes.SpaceID, moduleID coretypes.ExtID) *dal.Key {
	if spaceID == "" {
		// Spaceless system namespace: global/system extension records live at
		// /ext/{ext-id}/... with no /spaces/{space-id} prefix (sneat-specs
		// Decision 0002). NewSpaceKey is intentionally NOT called here.
		return dal.NewKeyWithID(SpaceModulesCollection, moduleID)
	}
	spaceKey := NewSpaceKey(spaceID)
	return dal.NewKeyWithParentAndID(spaceKey, SpaceModulesCollection, moduleID)
}

func NewSpaceModuleEntry[D any](spaceID coretypes.SpaceID, moduleID coretypes.ExtID, data D) record.DataWithID[coretypes.ExtID, D] {
	key := NewSpaceModuleKey(spaceID, moduleID)
	return record.NewDataWithID(moduleID, key, data)
}
