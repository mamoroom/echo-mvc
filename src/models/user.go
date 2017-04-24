package models

import (
	"github.com/mamoroom/echo-mvc/src/models/db/db_conns"
	"github.com/mamoroom/echo-mvc/src/models/db/db_conns/db_user_r"
	"github.com/mamoroom/echo-mvc/src/models/db/db_conns/db_user_w"
	"github.com/mamoroom/echo-mvc/src/models/entity"

	"fmt"
	"time"
)

func NewUserW() *UserModel {
	return NewUser(db_user_w.New())
}
func NewUserR() *UserModel {
	return NewUser(db_user_r.New())
}

func NewUser(db_conn db_conns.DbConn) *UserModel {
	b := NewBase(db_conn)
	return &UserModel{
		BaseModel: *b,
	}
}

type UserModel struct {
	BaseModel
	user *entity.User
}

func (u *UserModel) SetUserEntity(entity_user *entity.User) {
	u.user = entity_user
}

func (u *UserModel) GetUserEntity() *entity.User {
	return u.user
}

func (u *UserModel) IsUserEntityEmpty() bool {
	return *u.user == entity.User{}
}

func (u *UserModel) IsUserEntityNil() bool {
	return u.user == nil
}

func (u *UserModel) FindUserById(id uint64) error {
	var user entity.User
	_, err := u.Dbh.Handle().Id(id).Get(&user)
	if err != nil {
		return err
	}
	u.SetUserEntity(&user)
	return nil
}

func (u *UserModel) FindUserByAuth(auth_provider string, auth_account_id string) error {
	var user entity.User
	_, err := u.Dbh.Handle().Where("auth_provider = ?", auth_provider).And("auth_account_id = ?", auth_account_id).Get(&user)
	if err != nil {
		return err
	}
	u.SetUserEntity(&user)
	return nil
}

func (u *UserModel) DeleteAuth() (int64, error) {
	if u.IsUserEntityNil() || u.IsUserEntityEmpty() {
		panic("Must set user entity at first")
	}
	u.Dbh.SetNullable(u.user.GetLogoutNullableCols())
	return u.UpdateAuth("", "", "", "", time.Time{})
}

func (u *UserModel) DeleteChatAuth() (int64, error) {
	if u.IsUserEntityNil() || u.IsUserEntityEmpty() {
		panic("Must set user entity at first")
	}
	return u.UpdateTwitterAuth("", "", "", "")
}

func (u *UserModel) UpdateAuth(provider string, account_id string, access_token string, refresh_token string, expires_at time.Time) (int64, error) {
	if u.IsUserEntityNil() || u.IsUserEntityEmpty() {
		panic("Must set user entity at first")
	}

	_user := *u.user
	// set
	_user.AuthProvider = provider
	_user.AuthAccountId = account_id
	_user.AuthAccessToken = access_token
	_user.AuthRefreshToken = refresh_token
	_user.AuthExpiresAt = expires_at

	rows, err := u.Dbh.Handle().Id(u.GetUserEntity().Id).Cols("auth_provider", "auth_account_id", "auth_access_token", "auth_refresh_token", "auth_expires_at").Update(&_user)
	if err != nil {
		return 0, err
	}
	if rows > 0 {
		u.SetUserEntity(&_user)
	}
	return rows, nil
}

func (u *UserModel) UpdateTwitterAuth(twitter_user_id string, twitter_access_token string, twitter_access_token_secret string, twitter_avatar_url string) (int64, error) {
	if u.IsUserEntityNil() || u.IsUserEntityEmpty() {
		panic("Must set user entity at first")
	}

	_user := *u.user
	// set
	_user.TwitterUserId = twitter_user_id
	_user.TwitterAccessToken = twitter_access_token
	_user.TwitterAccessTokenSecret = twitter_access_token_secret
	_user.TwitterAvatarUrl = twitter_avatar_url

	rows, err := u.Dbh.Handle().Id(u.GetUserEntity().Id).Cols("twitter_user_id", "twitter_access_token", "twitter_access_token_secret", "twitter_avatar_url").Update(&_user)
	if err != nil {
		return 0, err
	}
	if rows > 0 {
		u.SetUserEntity(&_user)
	}
	return rows, nil
}

func (u *UserModel) UpdateTutorialDone(name string) (int64, error) {
	if u.IsUserEntityNil() || u.IsUserEntityEmpty() {
		panic("Must set user entity at first")
	}

	_user := *u.user
	// set
	_user.Name = name
	_user.IsTutorialDone = true

	rows, err := u.Dbh.Handle().Id(u.GetUserEntity().Id).Cols("name", "is_tutorial_done").Update(&_user)
	if err != nil {
		return 0, err
	}
	if rows > 0 {
		u.SetUserEntity(&_user)
	}
	return rows, nil
}

func (u *UserModel) InsertByAuth(lang string, provider string, account_id string, access_token string, refresh_token string, expires_at time.Time) (int64, error) {
	//[todo]: validation
	fmt.Println(expires_at)
	_user := entity.User{
		Name:             "",
		Lang:             lang,
		AuthProvider:     provider,
		AuthAccountId:    account_id,
		AuthAccessToken:  access_token,
		AuthRefreshToken: refresh_token,
		AuthExpiresAt:    expires_at,
	}
	u.Dbh.SetNullable(_user.GetInitNullableCols())
	rows, err := u.Dbh.Handle().Insert(&_user)
	if err != nil {
		return 0, err
	}
	if rows > 0 {
		u.SetUserEntity(&_user)
	}
	return rows, nil
}

func (u *UserModel) Insert(lang string) (int64, error) {
	//[todo]: validation
	_user := entity.User{
		Name: "",
		Lang: lang,
	}
	u.Dbh.SetNullable(_user.GetInitNullableCols())
	rows, err := u.Dbh.Handle().Insert(&_user)
	if err != nil {
		return 0, err
	}
	if rows > 0 {
		u.SetUserEntity(&_user)
	}
	return rows, nil
}

/*func (u *UserModel) SetUserAuthById(id uint64, auth_provider string, auth_account_id string, access_token string) (*User, error) {
}*/
