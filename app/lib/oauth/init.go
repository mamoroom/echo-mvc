package oauth

import (
	"github.com/prtytokyo/goth"
	"github.com/prtytokyo/goth/providers/gplus"
	"github.com/prtytokyo/goth/providers/l5id"
	"github.com/prtytokyo/goth/providers/openidConnect"
	"github.com/prtytokyo/goth/providers/twitter"

	"github.com/mamoroom/echo-mvc/app/config"
	"github.com/mamoroom/echo-mvc/app/lib/util"

	_ "errors"
	_ "fmt"
)

var conf = config.Conf

func init() {
	goth.UseProviders(
		getGplusProvider(conf.Oauth.Google, conf.Oauth.Login),
		getL5idProvider(conf.Oauth.L5id, conf.Oauth.Login),

		getGplusProvider(conf.Oauth.Google, conf.Oauth.App),
		getTwitterProvider(conf.Oauth.Twitter, conf.Oauth.Chat),
	)
}

func getTwitterProvider(conf_twitter config.TwitterConf, conf_path config.OauthPathConf) *twitter.Provider {
	provider_name := util.GetOauthProviderName(conf_twitter.Account.ProviderNamePrefix, conf_path.ProviderNameSuffix)
	tw := twitter.New(conf_twitter.Account.ClientKey, conf_twitter.Account.ClientSecret, getCallbackUrl(conf_path, conf.Server.Domain, provider_name))
	tw.SetName(provider_name)
	return tw
}

func getGplusProvider(conf_google config.GoogleConf, conf_path config.OauthPathConf) *gplus.Provider {
	provider_name := util.GetOauthProviderName(conf_google.Account.ProviderNamePrefix, conf_path.ProviderNameSuffix)
	gp := gplus.New(conf_google.Account.ClientKey, conf_google.Account.ClientSecret, getCallbackUrl(conf_path, conf.Server.Domain, provider_name), conf_path.Scopes...)
	gp.SetPrompt(conf_google.Prompt)
	gp.SetName(provider_name)
	return gp
}

func getL5idProvider(conf_l5id config.L5idConf, conf_path config.OauthPathConf) *l5id.Provider {
	provider_name := util.GetOauthProviderName(conf_l5id.Account.ProviderNamePrefix, conf_path.ProviderNameSuffix)
	l5id := l5id.New(conf_l5id.Account.ClientKey, conf_l5id.Account.ClientSecret, getCallbackUrl(conf_path, conf.Server.Domain, provider_name))
	l5id.SetName(provider_name)
	return l5id
}

func getOpenIdProvider(conf_google config.GoogleConf, conf_path config.OauthPathConf) *openidConnect.Provider {
	provider_name := util.GetOauthProviderName(conf_google.Account.ProviderNamePrefix, conf_path.ProviderNameSuffix)
	oc, err := openidConnect.New(conf_google.Account.ClientKey, conf_google.Account.ClientSecret, getCallbackUrl(conf_path, conf.Server.Domain, provider_name), conf_google.AutoDiscoveryUrl, conf_path.Scopes...)
	if err != nil {
		panic(err)
	}
	oc.SetName(provider_name)
	return oc
}

func getCallbackUrl(conf_path config.OauthPathConf, domain string, provider_name string) string {
	return domain + conf_path.CallbackUri + provider_name + "/callback"
}
