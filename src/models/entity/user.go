package entity

import (
	"time"
)

type User struct {
	Id                       uint64    `json:"id" xorm:"pk autoincr"`
	Name                     string    `json:"name" xorm:"notnull default ''"`
	Lang                     string    `json:"lang" xorm:"notnull"`
	Points                   uint64    `json:"points" xorm:"notnull default 0"`
	Coins                    uint64    `json:"coins" xorm:"notnull default 0"`
	AuthProvider             string    `json:"auth_provider" xorm:"unique(auth) "`
	AuthAccountId            string    `json:"auth_account_id" xorm:" unique(auth)"`
	AuthAccessToken          string    `json:"auth_access_token" xorm:"notnull default ''"`
	AuthRefreshToken         string    `json:"auth_refresh_token" xorm:"notnull default ''"`
	AuthExpiresAt            time.Time `json:"expires_at" xorm:"utc"`
	TwitterUserId            string    `json:"twitter_user_id" xorm:"unique(twitter)"`
	TwitterAccessToken       string    `json:"twitter_access_token" xorm:"notnull default ''"`
	TwitterAccessTokenSecret string    `json:"twitter_access_token_secret" xorm:"notnull default ''"`
	TwitterAvatarUrl         string    `json:"twitter_avater_url" xorm:"notnull default ''"`
	IsTutorialDone           bool      `json:"is_tutorial_done" xorm:"notnull default 0"`
	IsFirstGachaDone         bool      `json:"is_first_gacha_done" xorm:"notnull default 0"`
	IsNotificationsOn        bool      `json:"is_notifications_on" xorm:"notnull default 0"`
	SeqLoginCnt              uint8     `json:"seq_login_cnt" xorm:"notnull default 0"`
	CreatedAt                time.Time `json:"created_at" xorm:"created utc"`
	UpdatedAt                time.Time `json:"updated_at" xorm:"updated utc"`
}

func (u User) GetNullableCols() []string {
	return []string{"auth_provider", "auth_account_id", "twitter_user_id"}
}

func (u *User) IsAuthedUser() bool {
	return u.AuthProvider != ""
}

func (u *User) HasTwitterUserId() bool {
	return u.TwitterUserId != ""
}

func (u *User) GetLoginType() string {
	if !u.IsAuthedUser() {
		return "guest"
	}
	return u.AuthProvider
}
