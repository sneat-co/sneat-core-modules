package facade4invitus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-core-modules/userus/dal4userus"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"strings"
	"time"
)

type InviteClaimOperation string

const (
	InviteClaimOperationAccept  InviteClaimOperation = "accept"
	InviteClaimOperationDecline InviteClaimOperation = "decline"
)

var (
	ErrInvitePinDoesNotMatch = fmt.Errorf("%w: pin code does not match", facade.ErrBadRequest)
	ErrInviteAlreadyAccepted = fmt.Errorf("invite is already accepted")
	ErrInviteIsRevoked       = fmt.Errorf("invite is revoked")
	ErrInviteExpired         = fmt.Errorf("invite is expired")
)

// ClaimPersonalInviteRequest holds parameters for accepting a personal invite
type ClaimPersonalInviteRequest struct {
	InviteRequest
	Operation InviteClaimOperation `json:"operation"`

	NoPinRequired bool `json:"noPinRequired"`

	RemoteClient dbmodels.RemoteClientInfo `json:"remoteClient"`

	// TODO: Document why we need this and why it's called 'member'?
	//Member dbmodels.DtoWithID[*briefs4contactus.ContactBase] `json:"member"`

	//FullName string                      `json:"fullName"`
	//Email    string                      `json:"email"`
}

// Validate validates request
func (v *ClaimPersonalInviteRequest) Validate() error {
	if err := v.InviteRequest.Validate(); err != nil {
		return err
	}
	switch v.Operation {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("operation")
	case InviteClaimOperationAccept, InviteClaimOperationDecline:
		// OK
	default:
		return validation.NewErrBadRequestFieldValue("operation", "invalid value: "+string(v.Operation))
	}
	//if err := v.Member.Validate(); err != nil {
	//	return validation.NewErrBadRequestFieldValue("member", err.Error())
	//}
	return nil
}

type ClaimPersonalInviteResponse struct {
	Invite         InviteEntry
	Space          dbo4spaceus.SpaceEntry
	Contact        dal4contactus.ContactEntry
	ContactusSpace dal4contactus.ContactusSpaceEntry
}

// ClaimPersonalInvite accepts personal invite and joins user to a space.
// If needed, a user record should be created
func ClaimPersonalInvite(
	ctx facade.ContextWithUser, request ClaimPersonalInviteRequest,
) (
	response ClaimPersonalInviteResponse, err error,
) {
	if err = request.Validate(); err != nil {
		return
	}
	uid := ctx.User().GetUserID()

	err = dal4contactus.RunContactusSpaceWorker(ctx, request.SpaceRequest,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusSpaceWorkerParams) error {
			response.Space = params.Space
			response.ContactusSpace = params.SpaceModuleEntry

			if response.Invite, response.Contact, err = getPersonalInviteRecords(ctx, tx, params, request.InviteID); err != nil {
				return err
			}

			if response.Invite.Data.Pin != request.Pin {
				return ErrInvitePinDoesNotMatch
			}

			user := dbo4userus.NewUserEntry(uid)
			if err = dal4userus.GetUser(ctx, tx, user); err != nil {
				if !dal.IsNotFound(err) {
					return err
				}
			}

			now := time.Now()

			var oldInviteStatus = response.Invite.Data.Status
			var newInviteStatus dbo4invitus.InviteStatus
			switch request.Operation {
			case InviteClaimOperationAccept:
				if newInviteStatus = dbo4invitus.InviteStatusAccepted; oldInviteStatus != newInviteStatus {
					var spaceMember *briefs4contactus.ContactBase
					if spaceMember, err = updateContactusSpaceRecord(uid, response.Invite.Data.ToSpaceContactID, params, response.Contact); err != nil {
						return fmt.Errorf("failed to update space record: %w", err)
					}

					memberContext := dal4contactus.NewContactEntry(params.Space.ID, response.Contact.ID)

					if err = updateMemberRecord(ctx, tx, uid, memberContext, &response.Contact.Data.ContactBase, spaceMember); err != nil {
						return fmt.Errorf("failed to update space member record: %w", err)
					}

					if err = createOrUpdateUserRecord(ctx, tx, now, user, request, response.Contact, params, spaceMember, response.Invite); err != nil {
						return fmt.Errorf("failed to create or update user record: %w", err)
					}
				}
			case InviteClaimOperationDecline:
				newInviteStatus = dbo4invitus.InviteStatusDeclined
			}
			if newInviteStatus != oldInviteStatus {
				if err = updateInviteStatus(ctx, tx, uid, now, response.Invite, newInviteStatus); err != nil {
					return fmt.Errorf("failed to update invite record: %w", err)
				}
			}
			return err
		})
	return
}

func updateMemberRecord(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	uid string,
	member dal4contactus.ContactEntry,
	requestMember *briefs4contactus.ContactBase,
	spaceMember *briefs4contactus.ContactBase,
) (err error) {
	updates := []update.Update{update.ByFieldName("userID", uid)}
	updates = updatePersonDetails(&member.Data.ContactBase, requestMember, spaceMember, updates)
	if err = tx.Update(ctx, member.Key, updates); err != nil {
		return err
	}
	return err
}

func updateContactusSpaceRecord(
	uid, memberID string,
	params *dal4contactus.ContactusSpaceWorkerParams,
	requestMember dal4contactus.ContactEntry,
) (spaceMember *briefs4contactus.ContactBase, err error) {
	if uid == "" {
		panic("required parameter `uid` is empty string")
	}

	inviteToMemberID := memberID[strings.Index(memberID, ":")+1:]
	for contactID, m := range params.SpaceModuleEntry.Data.Contacts {
		if contactID == inviteToMemberID {
			m.UserID = uid
			params.SpaceModuleEntry.Data.AddUserID(uid)
			params.SpaceModuleEntry.Data.AddContact(contactID, m)
			//request.ContactID.Roles = m.Roles
			//m = request.ContactID
			m.UserID = uid
			spaceMember = &briefs4contactus.ContactBase{
				ContactBrief: *m,
			}
			//space.Members[i] = m
			updatePersonDetails(spaceMember, &requestMember.Data.ContactBase, spaceMember, nil)
			params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, params.SpaceModuleEntry.Data.AddUserID(uid)...)
			if m.AddRole(const4contactus.SpaceMemberRoleMember) {
				params.SpaceModuleUpdates = append(params.SpaceModuleUpdates,
					update.ByFieldPath([]string{"contacts", contactID, "roles"}, m.Roles))
			}
			break
		}
	}
	if spaceMember == nil {
		return spaceMember, fmt.Errorf("space member is not found by inviteToMemberID=%s", inviteToMemberID)
	}

	if params.Space.Data.HasUserID(uid) {
		goto UserIdAdded
	}
	params.SpaceUpdates = append(params.SpaceUpdates, update.ByFieldName("userIDs", params.Space.Data.UserIDs))
UserIdAdded:
	return spaceMember, err
}

func createOrUpdateUserRecord(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	now time.Time,
	user dbo4userus.UserEntry,
	request ClaimPersonalInviteRequest,
	member dal4contactus.ContactEntry,
	params *dal4contactus.ContactusSpaceWorkerParams,
	spaceMember *briefs4contactus.ContactBase,
	invite InviteEntry,
) (err error) {
	if spaceMember == nil {
		panic("spaceMember == nil")
	}
	existingUser := user.Record.Exists()
	if existingUser {
		spaceInfo := user.Data.GetUserSpaceInfoByID(request.SpaceID)
		if spaceInfo != nil {
			return nil
		}
	}

	userSpaceInfo := dbo4userus.UserSpaceBrief{
		SpaceBrief: params.Space.Data.SpaceBrief,
		Roles:      invite.Data.Roles, // TODO: Validate roles?
	}
	if err = userSpaceInfo.Validate(); err != nil {
		return fmt.Errorf("invalid user space info: %w", err)
	}
	user.Data.Spaces[string(request.SpaceID)] = &userSpaceInfo
	user.Data.SpaceIDs = append(user.Data.SpaceIDs, string(request.SpaceID))
	if existingUser {
		userUpdates := []update.Update{
			update.ByFieldName("spaces", user.Data.Spaces),
		}
		userUpdates = updatePersonDetails(&user.Data.ContactBase, &member.Data.ContactBase, spaceMember, userUpdates)
		if err = user.Data.Validate(); err != nil {
			return fmt.Errorf("user record prepared for update is not valid: %w", err)
		}
		if err = tx.Update(ctx, user.Key, userUpdates); err != nil {
			return fmt.Errorf("failed to update user record: %w", err)
		}
	} else { // New user record
		user.Data.CreatedAt = now
		user.Data.Created.Client = request.RemoteClient
		user.Data.Type = briefs4contactus.ContactTypePerson
		user.Data.Names = member.Data.Names
		if user.Data.Names.IsEmpty() {
			user.Data.Names = spaceMember.Names
		}
		updatePersonDetails(&user.Data.ContactBase, &member.Data.ContactBase, spaceMember, nil)
		if user.Data.Gender == "" {
			user.Data.Gender = "unknown"
		}
		if user.Data.CountryID == "" {
			user.Data.CountryID = with.UnknownCountryID
		}
		if len(member.Data.Emails) > 0 {
			user.Data.Emails = member.Data.Emails
		}
		if len(member.Data.Phones) > 0 {
			user.Data.Phones = member.Data.Phones
		}
		if err = user.Data.Validate(); err != nil {
			return fmt.Errorf("user record prepared for insert is not valid: %w", err)
		}
		if err = tx.Insert(ctx, user.Record); err != nil {
			return fmt.Errorf("failed to insert user record: %w", err)
		}
	}
	return err
}

func updatePersonDetails(personContact *briefs4contactus.ContactBase, member *briefs4contactus.ContactBase, spaceMember *briefs4contactus.ContactBase, updates []update.Update) []update.Update {
	if member.Names != nil {
		if personContact.Names == nil {
			personContact.Names = new(person.NameFields)
		}
		if personContact.Names.FirstName == "" {
			name := member.Names.FirstName
			if name == "" {
				name = spaceMember.Names.FirstName
			}
			if name != "" {
				personContact.Names.FirstName = name
				if updates != nil {
					updates = append(updates, update.ByFieldName("name.first", name))
				}
			}
		}
		if personContact.Names.LastName == "" {
			name := member.Names.LastName
			if name == "" {
				name = spaceMember.Names.LastName
			}
			if name != "" {
				personContact.Names.LastName = name
				if updates != nil {
					updates = append(updates, update.ByFieldName("name.last", name))
				}
			}
		}
		if personContact.Names.FullName == "" {
			name := member.Names.FullName
			if name == "" {
				name = spaceMember.Names.FullName
			}
			if name != "" {
				personContact.Names.FullName = name
				if updates != nil {
					updates = append(updates, update.ByFieldName("name.full", name))
				}
			}
		}
	}
	if personContact.Gender == "" || personContact.Gender == "unknown" {
		gender := member.Gender
		if gender == "" || gender == "unknown" {
			gender = spaceMember.Gender
		}
		if gender == "" {
			gender = "unknown"
		}
		if personContact.Gender == "" || gender != "unknown" {
			personContact.Gender = member.Gender
			if updates != nil {
				updates = append(updates, update.ByFieldName("gender", gender))
			}
		}
	}
	return updates
}
