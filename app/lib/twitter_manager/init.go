package twitter_manager

import (
	"github.com/ChimeraCoder/anaconda"

	"github.com/mamoroom/echo-mvc/app/config"
)

var conf = config.Conf

func init() {
	anaconda.SetConsumerKey(conf.Oauth.Twitter.Account.ClientKey)
	anaconda.SetConsumerSecret(conf.Oauth.Twitter.Account.ClientSecret)
}
