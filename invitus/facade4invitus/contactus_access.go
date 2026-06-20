package facade4invitus

import (
	"context"
	"strings"
	"time"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactusmodels/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
)

// ContactRequest defines a request for a single contact.
// It mirrors dto4contactus.ContactRequest so invitus does not depend on the contactus module.
type ContactRequest struct {
	dto4spaceus.SpaceRequest
	ContactID string `json:"contactID"`
}

// Validate returns error if request is invalid
func (v ContactRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.ContactID) == "" {
		return validation.NewErrRequestIsMissingRequiredField("contactID")
	}
	return nil
}

// MemberContact abstracts a single contactus contact record (a space member)
// without exposing contactus DAL/DBO types to invitus.
type MemberContact interface {
	// ID returns the contact ID.
	ID() string
	// Key returns the contact record key.
	Key() *dal.Key
	// Record returns the underlying contact record.
	Record() dal.Record
	// ContactBase exposes the contact's mutable base fields
	// (Names, Roles, UserID, etc.).
	ContactBase() *briefs4contactus.ContactBase
	// Emails returns the contact's email communication channels.
	Emails() map[string]*with.CommunicationChannelProps
	// Phones returns the contact's phone communication channels.
	Phones() map[string]*with.CommunicationChannelProps
}

// SpaceContactsSession wraps the contactus space worker params so invitus can
// run its flows inside a contactus transaction without depending on contactus
// DAL/DBO types.
type SpaceContactsSession interface {
	// Space returns the space entry.
	Space() dbo4spaceus.SpaceEntry

	// Contacts returns the space contact briefs map.
	Contacts() map[string]*briefs4contactus.ContactBrief

	// GetContactBriefByUserID returns the contact brief for the given user ID.
	GetContactBriefByUserID(userID string) (id string, brief *briefs4contactus.ContactBrief)

	// AddContact adds (or replaces) a contact brief in the space module DBO.
	AddContact(contactID string, brief *briefs4contactus.ContactBrief)

	// AddSpaceModuleUserID adds the user ID to the space module DBO and returns the resulting updates.
	AddSpaceModuleUserID(userID string) []update.Update

	// SpaceModuleUserIDs returns the user IDs stored on the space module DBO.
	SpaceModuleUserIDs() []string

	// SpaceModuleRecordExists reports whether the space module record exists.
	SpaceModuleRecordExists() bool

	// SpaceModuleKey returns the space module record key (for error reporting).
	SpaceModuleKey() *dal.Key

	// SpaceModuleRecordError returns the space module record error (for error reporting).
	SpaceModuleRecordError() error

	// AppendSpaceModuleUpdates appends updates to be applied to the space module record.
	AppendSpaceModuleUpdates(updates ...update.Update)

	// AppendSpaceUpdates appends updates to be applied to the space record.
	AppendSpaceUpdates(updates ...update.Update)

	// GetRecords batch-loads the space-module record plus any extra records.
	GetRecords(ctx context.Context, tx dal.ReadSession, extraRecords ...dal.Record) error

	// NewMemberContact creates a member contact entry (not loaded) for the given contact ID.
	NewMemberContact(contactID string) MemberContact

	// LoadMemberContact creates and loads a member contact entry via getter.Get.
	LoadMemberContact(ctx context.Context, getter dal.ReadSession, contactID string) (MemberContact, error)
}

// ContactSession is a SpaceContactsSession scoped to a single contact, exposing
// the invite-brief operations for that contact.
type ContactSession interface {
	SpaceContactsSession

	// Contact returns the contact this session is scoped to.
	Contact() MemberContact

	// GetContactInviteBriefByChannelAndInviterUserID returns the invite brief ID for the channel + inviter.
	GetContactInviteBriefByChannelAndInviterUserID(channel dbo4invitus.InviteChannel, inviterUserID string) (id string)

	// AppendContactDeleteInviteBrief deletes the invite brief and queues the contact update.
	AppendContactDeleteInviteBrief(inviteID string)

	// AppendContactAddInviteBrief adds the invite brief and queues the contact update.
	AppendContactAddInviteBrief(inviteID, createdByUserID string, channel dbo4invitus.InviteChannel, createdTime time.Time)
}

// ContactusAccess is the registered accessor invitus uses to run contactus
// transactions and load member briefs without importing contactus module packages.
type ContactusAccess interface {
	// RunSpaceContactsTx runs worker inside a read-write contactus space transaction.
	RunSpaceContactsTx(
		ctx facade.ContextWithUser,
		request dto4spaceus.SpaceRequest,
		worker func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, session SpaceContactsSession) error,
	) error

	// RunReadonlySpaceContactsTx runs worker inside a readonly contactus space transaction.
	RunReadonlySpaceContactsTx(
		ctx context.Context,
		userCtx facade.UserContext,
		request dto4spaceus.SpaceRequest,
		worker func(ctx context.Context, tx dal.ReadTransaction, session SpaceContactsSession) error,
	) error

	// RunContactTx runs worker inside a read-write contactus contact transaction.
	RunContactTx(
		ctx facade.ContextWithUser,
		request ContactRequest,
		worker func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, session ContactSession) error,
	) error

	// GetSpaceMemberContactBrief loads a single space member contact brief via a plain db.Get.
	GetSpaceMemberContactBrief(
		ctx context.Context,
		getter dal.ReadSession,
		spaceID coretypes.SpaceID,
		contactID string,
	) (*briefs4contactus.ContactBrief, error)
}

var contactusAccess ContactusAccess

// RegisterContactusAccess registers the contactus implementation used by invitus.
// Called once at startup from contactus.Extension().
func RegisterContactusAccess(a ContactusAccess) {
	contactusAccess = a
}
