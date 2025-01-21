package unsorted4auth

import (
	"context"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/dal-go/dalgo/dal"
	models4auth2 "github.com/sneat-co/sneat-core-modules/auth/models4auth"
	//"github.com/sneat-co/sneat-core-modules/bots/anybot"
	dbo4userus2 "github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"time"
)

type CreateUserData struct {
	//FbUserID     string
	//GoogleUserID string
	//VkUserID     int64
	FirstName  string
	LastName   string
	ScreenName string
	Nickname   string
}

// CreateUserEntity
// Deprecated
func CreateUserEntity(createUserData CreateUserData) (user dbo4userus2.UserEntry) {
	return
	//return &models4debtus.DebutsAppUserDataOBSOLETE{
	//	//FbUserID: createUserData.FbUserID,
	//	//VkUserID: createUserData.VkUserID,
	//	//GoogleUniqueUserID: createUserData.GoogleUserID,
	//	//ContactDetails: dto4contactus.ContactDetails{
	//	//	NameFields: person.NameFields{
	//	//		FirstName:  createUserData.FirstName,
	//	//		LastName:   createUserData.LastName,
	//	//		ScreenName: createUserData.ScreenName,
	//	//		NickName:   createUserData.Nickname,
	//	//	},
	//	//},
	//}
}

type UserDal interface {
	GetUserByStrID(ctx context.Context, userID string) (dbo4userus2.UserEntry, error)
	GetUserByVkUserID(ctx context.Context, vkUserID int64) (dbo4userus2.UserEntry, error)
	CreateAnonymousUser(ctx context.Context) (dbo4userus2.UserEntry, error)
	CreateUser(ctx context.Context, userEntity *dbo4userus2.UserDbo) (dbo4userus2.UserEntry, error)
	DelaySetUserPreferredLocale(ctx context.Context, delay time.Duration, userID string, localeCode5 string) error
}

type PasswordResetDal interface {
	GetPasswordResetByID(ctx context.Context, tx dal.ReadSession, id int) (models4auth2.PasswordReset, error)
	CreatePasswordResetByID(ctx context.Context, tx dal.ReadwriteTransaction, entity *models4auth2.PasswordResetData) (models4auth2.PasswordReset, error)
	SavePasswordResetByID(ctx context.Context, tx dal.ReadwriteTransaction, record models4auth2.PasswordReset) (err error)
}

type UserGoogleDal interface {
	GetUserGoogleByID(ctx context.Context, googleUserID string) (userGoogle models4auth2.UserAccountEntry, err error)
	DeleteUserGoogle(ctx context.Context, googleUserID string) (err error)
}

type UserVkDal interface {
	GetUserVkByID(ctx context.Context, vkUserID int64) (userGoogle models4auth2.UserVk, err error)
	SaveUserVk(ctx context.Context, userVk models4auth2.UserVk) (err error)
}

type UserEmailDal interface {
	GetUserEmailByID(ctx context.Context, tx dal.ReadSession, email string) (userEmail models4auth2.UserEmailEntry, err error)
	SaveUserEmail(ctx context.Context, tx dal.ReadwriteTransaction, userEmail models4auth2.UserEmailEntry) (err error)
}

type UserGooglePlusDal interface {
	GetUserGooglePlusByID(ctx context.Context, id string) (userGooglePlus models4auth2.UserGooglePlus, err error)
	//SaveUserGooglePlusByID(ctx context.Context, userGooglePlus models4auth.UserGooglePlus) (err error)
}

type UserFacebookDal interface {
	GetFbUserByFbID(ctx context.Context, fbAppOrPageID, fbUserOrPageScopeID string) (fbUser models4auth2.UserFacebook, err error)
	SaveFbUser(ctx context.Context, tx dal.ReadwriteTransaction, fbUser models4auth2.UserFacebook) (err error)
	DeleteFbUser(ctx context.Context, fbAppOrPageID, fbUserOrPageScopeID string) (err error)
	//CreateFbUserRecord(ctx context.Context, fbUserID string, appUserID int64) (fbUser models.UserFacebook, err error)
}

type LoginPinDal interface {
	GetLoginPinByID(ctx context.Context, tx dal.ReadSession, loginID int) (loginPin models4auth2.LoginPin, err error)
	SaveLoginPin(ctx context.Context, tx dal.ReadwriteTransaction, loginPin models4auth2.LoginPin) (err error)
	CreateLoginPin(ctx context.Context, tx dal.ReadwriteTransaction, channel, gaClientID string, createdUserID string) (loginPin models4auth2.LoginPin, err error)
}

type LoginCodeDal interface {
	NewLoginCode(ctx context.Context, userID string) (code int, err error)
	ClaimLoginCode(ctx context.Context, code int) (userID string, err error)
}

//type TgChatDal interface {
//	GetTgChatByID(ctx context.Context, tgBotID string, tgChatID int64) (tgChat anybot.SneatAppTgChatEntry, err error)
//	DoSomething( // TODO: WTF name?
//		ctx context.Context,
//		userTask *sync.WaitGroup,
//		currency string,
//		tgChatID int64,
//		authInfo token4auth.AuthInfo,
//		user dbo4userus2.UserEntry,
//		sendToTelegram func(tgChat botsfwtgmodels.TgChatData) error,
//	) (err error)
//}

type TgUserDal interface {
	FindByUserName(ctx context.Context, tx dal.ReadSession, userName string) (tgUsers []botsfwtgmodels.TgPlatformUser, err error)
}

var User UserDal

var UserFacebook UserFacebookDal

var UserGooglePlus UserGooglePlusDal

var PasswordReset PasswordResetDal

var UserEmail UserEmailDal

//var TgChat TgChatDal

var TgUser TgUserDal
