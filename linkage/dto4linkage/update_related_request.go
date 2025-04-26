package dto4linkage

import (
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

type UpdateRelatedRequest struct {
	dto4spaceus.SpaceItemRequest
	coretypes.ModuleCollectionRef
	UpdateRelatedFieldRequest
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
