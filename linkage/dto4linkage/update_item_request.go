package dto4linkage

import (
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
)

type UpdateItemRequest struct {
	dto4spaceus.SpaceRequest
	dbo4linkage.ItemRef `json:"itemRef"`
	UpdateRelatedFieldRequest
}

func (v *UpdateItemRequest) Validate() error {
	if err := v.ItemRef.Validate(); err != nil {
		return err
	}
	if err := v.UpdateRelatedFieldRequest.Validate(); err != nil {
		return err
	}
	return nil
}
