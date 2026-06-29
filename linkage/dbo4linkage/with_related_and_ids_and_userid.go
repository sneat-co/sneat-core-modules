package dbo4linkage

import (
	"github.com/sneat-co/sneat-go-core/acl/dbo4acl"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

type WithRelatedAndIDsAndUserID struct {
	dbmodels.WithUserID
	*WithRelatedAndIDs

	// ACL is the record's per-record access-control list. It is populated when
	// loading a spaceless system-namespace record (/ext/{ext-id}/...), whose
	// writes are authorized per-record rather than by space membership (Decision
	// 0002). Space-bound records leave it nil and are gated by space membership.
	// specscore: features/reserved-extension-space-ids/R4
	ACL *dbo4acl.ACL `json:"acl,omitempty" firestore:"acl,omitempty"`
}

func (v *WithRelatedAndIDsAndUserID) Validate() error {
	if err := v.WithUserID.Validate(); err != nil {
		return err
	}
	if err := v.WithRelatedAndIDs.Validate(); err != nil {
		return err
	}
	if v.ACL != nil {
		if err := v.ACL.Validate(); err != nil {
			return err
		}
	}
	return nil
}
