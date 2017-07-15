package twitter_manager

import (
	_ "fmt"
	"github.com/ChimeraCoder/anaconda"
)

type TwitterManager struct {
	*anaconda.TwitterApi
}

type ResPostTweet struct {
	Tweet *anaconda.Tweet
	Err   error
}

func New(access_token string, access_token_secret string) *TwitterManager {
	tw_api := anaconda.NewTwitterApi(access_token, access_token_secret)
	tw_api.SetLogger(anaconda.BasicLogger)
	return &TwitterManager{
		TwitterApi: tw_api,
	}
}
