package contactus

import (
	"context"
	"time"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/contactusmodels/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-core-modules/invitus/facade4invitus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/strongoapp/with"
)

// invitusContactusAccess implements facade4invitus.ContactusAccess.
// It runs contactus workers and adapts contactus DAL/DBO types to the
// contactus-free interfaces consumed by invitus, keeping invitus decoupled
// from the contactus module while preserving transaction atomicity (invitus
// owns the whole transaction via a single worker callback).
type invitusContactusAccess struct{}

// memberContact adapts dal4contactus.ContactEntry to facade4invitus.MemberContact.
type memberContact struct {
	entry dal4contactus.ContactEntry
}

func (m memberContact) ID() string                                 { return m.entry.ID }
func (m memberContact) Key() *dal.Key                              { return m.entry.Key }
func (m memberContact) Record() dal.Record                         { return m.entry.Record }
func (m memberContact) ContactBase() *briefs4contactus.ContactBase { return &m.entry.Data.ContactBase }
func (m memberContact) Emails() map[string]*with.CommunicationChannelProps {
	return m.entry.Data.Emails
}
func (m memberContact) Phones() map[string]*with.CommunicationChannelProps {
	return m.entry.Data.Phones
}

// spaceContactsSession adapts *dal4contactus.ContactusSpaceWorkerParams to facade4invitus.SpaceContactsSession.
type spaceContactsSession struct {
	params *dal4contactus.ContactusSpaceWorkerParams
}

func (s *spaceContactsSession) Space() dbo4spaceus.SpaceEntry {
	return s.params.Space
}

func (s *spaceContactsSession) Contacts() map[string]*briefs4contactus.ContactBrief {
	return s.params.SpaceModuleEntry.Data.Contacts
}

func (s *spaceContactsSession) GetContactBriefByUserID(userID string) (id string, brief *briefs4contactus.ContactBrief) {
	return s.params.SpaceModuleEntry.Data.GetContactBriefByUserID(userID)
}

func (s *spaceContactsSession) AddContact(contactID string, brief *briefs4contactus.ContactBrief) {
	s.params.SpaceModuleEntry.Data.AddContact(contactID, brief)
}

func (s *spaceContactsSession) AddSpaceModuleUserID(userID string) []update.Update {
	return s.params.SpaceModuleEntry.Data.AddUserID(userID)
}

func (s *spaceContactsSession) SpaceModuleUserIDs() []string {
	return s.params.SpaceModuleEntry.Data.UserIDs
}

func (s *spaceContactsSession) SpaceModuleRecordExists() bool {
	return s.params.SpaceModuleEntry.Record.Exists()
}

func (s *spaceContactsSession) SpaceModuleKey() *dal.Key {
	return s.params.SpaceModuleEntry.Key
}

func (s *spaceContactsSession) SpaceModuleRecordError() error {
	return s.params.SpaceModuleEntry.Record.Error()
}

func (s *spaceContactsSession) AppendSpaceModuleUpdates(updates ...update.Update) {
	s.params.SpaceModuleUpdates = append(s.params.SpaceModuleUpdates, updates...)
}

func (s *spaceContactsSession) AppendSpaceUpdates(updates ...update.Update) {
	s.params.SpaceUpdates = append(s.params.SpaceUpdates, updates...)
}

func (s *spaceContactsSession) GetRecords(ctx context.Context, tx dal.ReadSession, extraRecords ...dal.Record) error {
	return s.params.GetRecords(ctx, tx, extraRecords...)
}

func (s *spaceContactsSession) NewMemberContact(contactID string) facade4invitus.MemberContact {
	return memberContact{entry: dal4contactus.NewContactEntry(s.params.Space.ID, contactID)}
}

func (s *spaceContactsSession) LoadMemberContact(ctx context.Context, getter dal.ReadSession, contactID string) (facade4invitus.MemberContact, error) {
	entry := dal4contactus.NewContactEntry(s.params.Space.ID, contactID)
	if err := getter.Get(ctx, entry.Record); err != nil {
		return memberContact{entry: entry}, err
	}
	return memberContact{entry: entry}, nil
}

// contactSession adapts *dal4contactus.ContactWorkerParams to facade4invitus.ContactSession.
type contactSession struct {
	spaceContactsSession
	params *dal4contactus.ContactWorkerParams
}

func (s *contactSession) Contact() facade4invitus.MemberContact {
	return memberContact{entry: s.params.Contact}
}

func (s *contactSession) GetContactInviteBriefByChannelAndInviterUserID(channel dbo4invitus.InviteChannel, inviterUserID string) (id string) {
	id, _ = s.params.Contact.Data.GetInviteBriefByChannelAndInviterUserID(channel, inviterUserID)
	return id
}

func (s *contactSession) AppendContactDeleteInviteBrief(inviteID string) {
	s.params.ContactUpdates = append(s.params.ContactUpdates, s.params.Contact.Data.DeleteInviteBrief(inviteID))
}

func (s *contactSession) AppendContactAddInviteBrief(inviteID, createdByUserID string, channel dbo4invitus.InviteChannel, createdTime time.Time) {
	s.params.ContactUpdates = append(s.params.ContactUpdates,
		s.params.Contact.Data.AddInviteBrief(inviteID, createdByUserID, channel, createdTime))
}

func (invitusContactusAccess) RunSpaceContactsTx(
	ctx facade.ContextWithUser,
	request dto4spaceus.SpaceRequest,
	worker func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, session facade4invitus.SpaceContactsSession) error,
) error {
	return dal4contactus.RunContactusSpaceWorker(ctx, request,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusSpaceWorkerParams) error {
			return worker(ctx, tx, &spaceContactsSession{params: params})
		})
}

func (invitusContactusAccess) RunReadonlySpaceContactsTx(
	ctx context.Context,
	userCtx facade.UserContext,
	request dto4spaceus.SpaceRequest,
	worker func(ctx context.Context, tx dal.ReadTransaction, session facade4invitus.SpaceContactsSession) error,
) error {
	return dal4contactus.RunReadonlyContactusSpaceWorker(ctx, userCtx, request,
		func(ctx context.Context, tx dal.ReadTransaction, params *dal4contactus.ContactusSpaceWorkerParams) error {
			return worker(ctx, tx, &spaceContactsSession{params: params})
		})
}

func (invitusContactusAccess) RunContactTx(
	ctx facade.ContextWithUser,
	request facade4invitus.ContactRequest,
	worker func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, session facade4invitus.ContactSession) error,
) error {
	contactRequest := dto4contactusContactRequest(request)
	return dal4contactus.RunContactWorker(ctx, contactRequest,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4contactus.ContactWorkerParams) error {
			session := &contactSession{
				spaceContactsSession: spaceContactsSession{params: params.ContactusSpaceWorkerParams},
				params:               params,
			}
			return worker(ctx, tx, session)
		})
}

func (invitusContactusAccess) GetSpaceMemberContactBrief(
	ctx context.Context,
	getter dal.ReadSession,
	spaceID coretypes.SpaceID,
	contactID string,
) (*briefs4contactus.ContactBrief, error) {
	member := dal4contactus.NewContactEntry(spaceID, contactID)
	if err := getter.Get(ctx, member.Record); err != nil {
		return nil, err
	}
	return &member.Data.ContactBrief, nil
}

func dto4contactusContactRequest(request facade4invitus.ContactRequest) dto4contactus.ContactRequest {
	return dto4contactus.ContactRequest{
		SpaceRequest: request.SpaceRequest,
		ContactID:    request.ContactID,
	}
}
