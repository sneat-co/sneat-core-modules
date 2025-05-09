package dto4userus

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/security"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/strongoapp/strongoauth"
	"github.com/strongo/validation"
	"net/mail"
	"strings"
)

var _ facade.Request = (*InitUserRecordRequest)(nil)

type InitSpaceInfo struct {
	Type  coretypes.SpaceType `json:"type"`
	Title string              `json:"title"`
}

func (v InitSpaceInfo) Validate() error {
	if strings.TrimSpace(string(v.Type)) == "" {
		return validation.NewErrRequestIsMissingRequiredField("type")
	}
	if strings.TrimSpace(v.Title) == "" {
		return validation.NewErrRequestIsMissingRequiredField("title")
	}
	return nil
}

// InitUserRecordRequest request
type InitUserRecordRequest struct {
	AuthProvider    string             `json:"authProvider,omitempty"`
	Email           string             `json:"email,omitempty"`
	EmailIsVerified bool               `json:"emailIsVerified,omitempty"`
	IanaTimezone    string             `json:"ianaTimezone,omitempty"`
	Names           *person.NameFields `json:"names"`

	// RemoteClient contains information about the remote client making the request.
	RemoteClient dbmodels.RemoteClientInfo `json:"remoteClient"`
}

// Validate validates request
func (v *InitUserRecordRequest) Validate() error {
	if v.AuthProvider != "" {
		if err := strongoauth.ValidateAuthProviderID(v.AuthProvider); err != nil {
			return validation.NewErrBadRequestFieldValue("authProvider", err.Error())
		}
	}
	if v.Names != nil {
		if err := v.Names.Validate(); err != nil {
			return fmt.Errorf("%w: %v", facade.ErrBadRequest, err)
		}
	}
	if v.Email != "" {
		if _, err := mail.ParseAddress(v.Email); err != nil {
			return validation.NewErrBadRequestFieldValue("email", err.Error())
		}
	}
	//if v.SpaceID != nil {
	//	if err := v.SpaceID.Validate(); err != nil {
	//		return validation.NewErrBadRecordFieldValue("space", err.Error())
	//	}
	//}
	return nil
}

// CreateUserRequest DTO
type CreateUserRequest struct {
	Creator string `json:"creator"`
	Title   string `json:"title,omitempty"`
	Email   string `json:"email"`
}

// Validate validates request
func (v *CreateUserRequest) Validate() error {
	if err := validate.OptionalEmail(v.Email, "email"); err != nil {
		return err
	}
	if v.Creator == "" {
		return validation.NewErrRecordIsMissingRequiredField("creator")
	} else if !security.IsKnownHost(v.Creator) {
		return validation.NewErrBadRequestFieldValue("creator", "unknown creator: "+v.Creator)
	}
	return nil
}

// CreateUserRequestWithRemoteClientInfo a request DTO
type CreateUserRequestWithRemoteClientInfo struct {
	CreateUserRequest
	RemoteClient dbmodels.RemoteClientInfo
}
