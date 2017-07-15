package entity

import (
	"time"
)

type User struct {
	Id                   uint64    `json:"id" xorm:"pk autoincr"`
	UserName             string    `json:"user_name" xorm:"unique(username) notnull default ''"`
	DisplayName          string    `json:"display_name" xorm:"notnull default ''"`
	LoginType            string    `json:"login_type" xorm:"notnull default ''"`
	Thumbnail            string    `json:"thumbnail" xorm:"notnull default ''"`
	SupportsRemainingNum uint64    `json:"supports_remaining_num" xorm:"notnull default 0"`
	SupportsNumRecoverAt time.Time `json:"supports_num_recover_at"`

	//email
	Email           string `json:"email" xorm:"unique(email)"`
	IsEmailVerified bool   `json:"is_email_verified" xorm:"notnull default 0"`

	//twitter
	TwitterUserId            string `json:"twitter_user_id" xorm:"notnull default ''"`
	TwitterAccessToken       string `json:"twitter_access_token" xorm:"notnull default ''"`
	TwitterAccessTokenSecret string `json:"twitter_access_token_secret" xorm:"notnull default ''"`

	//cancdidate
	IsCandidate bool `json:"is_cadidate" xorm:"notnull default 0"`

	//time
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	UpdatedAt time.Time `json:"updated_at" xorm:"updated"`
}

func (u User) GetInitNullableCols() []string {
	return []string{"email"}
}

func (u User) GetInitOmitCols() []string {
	//golangの初期値とdefaultがずれるときは、omitさせる
	return []string{"supports_remaining_num"}
}
