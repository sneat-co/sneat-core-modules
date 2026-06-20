package contactus

import (
	"fmt"
	"time"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/contactusmodels/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/contactusmodels/const4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

// spaceusContactusContributor implements facade4spaceus.ContactusSpaceContributor.
// It builds the contactus records (contactus space + creator member) that must be
// persisted when a new space is created, keeping spaceus decoupled from contactus DAL types.
type spaceusContactusContributor struct{}

func (spaceusContactusContributor) BuildSpaceCreationRecords(
	spaceID coretypes.SpaceID,
	userContactID string,
	creatorBrief briefs4contactus.ContactBrief,
	createdAt time.Time,
	byUserID string,
) (records []dal.Record, err error) {
	contactusSpace := dal4contactus.NewContactusSpaceEntry(spaceID)
	contactusSpace.Data.AddContact(userContactID, &creatorBrief)
	if err = contactusSpace.Data.Validate(); err != nil {
		return nil, fmt.Errorf("newly created contactus space record is not valid: %w", err)
	}
	contactusSpace.Record.MarkAsChanged()
	contactusSpace.Record.SetError(nil)

	member, err := newMemberContactEntryFromContactBrief(spaceID, userContactID, creatorBrief, createdAt, byUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to create member's record: %w", err)
	}
	return []dal.Record{contactusSpace.Record, member.Record}, nil
}

// newMemberContactEntryFromContactBrief creates a member record from member's brief.
// Moved here from spaceus/facade4spaceus/member_helpers.go as part of the contactus cycle-break.
func newMemberContactEntryFromContactBrief(
	spaceID coretypes.SpaceID,
	contactID string,
	memberBrief briefs4contactus.ContactBrief,
	now time.Time,
	byUserID string,
) (
	member dal4contactus.ContactEntry,
	err error,
) {
	if err = memberBrief.Validate(); err != nil {
		return member, fmt.Errorf("supplied member brief is not valid: %w", err)
	}
	member = dal4contactus.NewContactEntry(spaceID, contactID)
	member.Data.ContactBrief = memberBrief
	member.Data.Status = dbmodels.StatusActive
	_ = member.Data.AddRole(const4contactus.SpaceMemberRoleMember)
	member.Data.CreatedAt = now
	member.Data.CreatedBy = byUserID
	dbo4linkage.UpdateRelatedIDs(spaceID, &member.Data.WithRelated, &member.Data.WithRelatedIDs)
	member.Data.IncreaseVersion(now, byUserID)
	if err = member.Data.Validate(); err != nil {
		member.Record.SetError(err)
		return member, fmt.Errorf("failed to validate member data: %w", err)
	}
	member.Record.SetError(nil)
	return member, nil
}
