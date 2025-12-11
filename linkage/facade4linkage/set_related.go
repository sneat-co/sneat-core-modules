package facade4linkage

import (
	"fmt"
	"time"

	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

//type RelatableAdapter[D dbo4linkage.Relatable] interface {
//	VerifyItem(ctx context.Context, tx dal.ReadTransaction, recordRef dbo4linkage.ItemRef) (err error)
//	//GetRecord(ctx context.Context, tx dal.ReadTransaction, recordRef dbo4linkage.ItemRef) (record.DataWithID[string, D], error)
//}
//type relatableAdapter[D dbo4linkage.Relatable] struct {
//	verifyItem func(ctx context.Context, tx dal.ReadTransaction, recordRef dbo4linkage.ItemRef) (err error)
//}
//
//func (v relatableAdapter[D]) VerifyItem(ctx context.Context, tx dal.ReadTransaction, recordRef dbo4linkage.ItemRef) (err error) {
//	return v.verifyItem(ctx, tx, recordRef)
//}
//
//func NewRelatableAdapter[D dbo4linkage.Relatable](
//	verifyItem func(ctx context.Context, tx dal.ReadTransaction, recordRef dbo4linkage.ItemRef) (err error),
//) RelatableAdapter[D] {
//	return relatableAdapter[D]{
//		verifyItem: verifyItem,
//	}
//}

//func (relatableAdapter[D]) GetRecord(ctx context.Context, tx dal.ReadTransaction, recordRef dbo4linkage.ItemRef) (record.DataWithID[string, D], error) {
//	return nil, nil
//}

type SetRelatedResult struct {
	ItemUpdates []update.Update
}

// SetRelated updates related records to define relationships
func SetRelated(
	now time.Time,
	userID string,
	spaceID coretypes.SpaceID,
	object dbo4linkage.Relatable,
	objectRef dbo4linkage.ItemRef,
	command dbo4linkage.RelationshipItemRolesCommand,
) (
	result SetRelatedResult, err error,
) {

	{
		const invalidArgPrefix = "facade4linkage.SetRelated got invalid argument"
		if err = objectRef.Validate(); err != nil {
			err = fmt.Errorf("%s `objectRef dbo4linkage.ItemRef`: %w", invalidArgPrefix, err)
			return
		}
		if err = command.Validate(); err != nil {
			return
		}
	}

	var relUpdates []update.Update

	objectWithRelated := object.GetRelated()
	if objectWithRelated.Related == nil {
		objectWithRelated.Related = make(dbo4linkage.RelatedModules, 1)
	}

	// TODO: Duplicate of GetOppositeRole()! - needs to be in 1 place and document why 1 place chosen over the other
	/*addReciprocal := func(roles1, roles2 []dbo4linkage.RelationshipRoleID) []dbo4linkage.RelationshipRoleID {
		for _, r := range roles1 {
			switch r {
			case dbo4linkage.RelationshipRoleSibling, dbo4linkage.RelationshipRoleSpouse:
				// mutual relationships
				if !slices.Contains(roles2, r) {
					roles2 = append(roles2, r)
				}
			case dbo4linkage.RelationshipRoleParent:
				if !slices.Contains(roles2, dbo4linkage.RelationshipRoleChild) {
					roles2 = append(roles2, dbo4linkage.RelationshipRoleChild)
				}
			case dbo4linkage.RelationshipRoleChild:
				if !slices.Contains(roles2, dbo4linkage.RelationshipRoleParent) {
					roles2 = append(roles2, dbo4linkage.RelationshipRoleParent)
				}
			case dbo4linkage.RelationshipRoleManager:
				if !slices.Contains(roles2, dbo4linkage.RelationshipRoleParent) {
					roles2 = append(roles2, dbo4linkage.RelationshipRoleParent)
				}
			}
		}
		return roles2
	}*/

	//command.Add.RolesToItem = addReciprocal(command.Add.RolesOfItem, command.Add.RolesToItem)
	//command.Add.RolesOfItem = addReciprocal(command.Add.RolesToItem, command.Add.RolesOfItem)

	//makeRelationships := func(ids []string, now time.Time) (relationships dbo4linkage.RelationshipRoles) {
	//	relationships = make(dbo4linkage.RelationshipRoles, len(ids))
	//	for _, r := range ids {
	//		relationships[r] = &dbo4linkage.RelationshipRole{
	//			CreatedField: with.CreatedField{
	//				Created: with.Created{
	//					At: now.Format(time.DateOnly),
	//					By: userID,
	//				},
	//			},
	//		}
	//	}
	//	return
	//}
	//rolesOfItem := makeRelationships(command.Add.RolesOfItem, now)
	//rolesToItem := makeRelationships(command.Add.RolesToItem, now)
	//
	//objectWithRelated.AddRelationshipsAndIDs(
	//	itemRef,
	//	rolesOfItem,
	//	rolesToItem,
	//)
	if relUpdates, err = objectWithRelated.AddRelationshipAndID(now, userID, spaceID, command); err != nil {
		return
	}
	result.ItemUpdates = append(result.ItemUpdates, relUpdates...)

	//for _, itemUpdate := range itemUpdates {
	//	if strings.HasSuffix(itemUpdate.Field, "relatedIDs") {
	//		continue // Ignore relatedIDs for spaceModuleUpdates
	//	}
	//	spaceModuleUpdates = append(spaceModuleUpdates, update.Update{
	//		Field: fmt.Sprintf("%s.%s.%s", objectRef.Collection, objectRef.ItemID, itemUpdate.Field),
	//		Value: itemUpdate.Value,
	//	})
	//}

	return
}
