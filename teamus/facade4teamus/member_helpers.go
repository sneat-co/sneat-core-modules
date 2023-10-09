package facade4teamus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"time"
)

// CreateMemberRecordFromBrief creates a member record from member's brief
func CreateMemberRecordFromBrief(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	teamID string,
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
	member = dal4contactus.NewContactEntry(teamID, contactID)
	//member.Brief = &memberBrief
	//member.Data.TeamID = teamID
	member.Data.ContactBrief = memberBrief
	member.Data.Status = dbmodels.StatusActive
	_ = member.Data.AddRole(const4contactus.TeamMemberRoleTeamMember)
	member.Data.CreatedAt = now
	member.Data.CreatedBy = byUserID
	member.Data.IncreaseVersion(now, byUserID)
	if err = member.Data.Validate(); err != nil {
		return member, fmt.Errorf("failed to validate member data: %w", err)
	}
	if err := tx.Insert(ctx, member.Record); err != nil {
		return member, fmt.Errorf("failed to inser member record into DB: %w", err)
	}
	return member, nil
}
