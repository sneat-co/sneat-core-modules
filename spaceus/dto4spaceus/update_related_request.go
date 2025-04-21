package dto4spaceus

import (
	"github.com/sneat-co/sneat-core-modules/linkage/dto4linkage"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

type UpdateRelatedRequest struct {
	SpaceItemRequest
	coretypes.ModuleCollectionRef
	dto4linkage.UpdateRelatedFieldRequest
}

func (v *UpdateRelatedRequest) Validate() error {
	if err := v.SpaceItemRequest.Validate(); err != nil {
		return err
	}
	if err := v.ModuleCollectionRef.Validate(); err != nil {
		return err
	}
	if err := v.UpdateRelatedFieldRequest.Validate(); err != nil {
		return err
	}
	return nil
}
