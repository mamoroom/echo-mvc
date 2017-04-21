package oauth

import (
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/markbates/goth/providers/twitter"

	"github.com/mamoroom/echo-mvc/src/config"
	"github.com/mamoroom/echo-mvc/src/lib/util"

	_ "errors"
	_ "fmt"
)

var conf = config.Conf

func init() {
	goth.UseProviders(
		getGplusProvider(conf.Oauth.Google, conf.Oauth.Login),
		getGplusProvider(conf.Oauth.Google, conf.Oauth.App),
		getTwitterProvider(conf.Oauth.Twitter, conf.Oauth.Chat),
	)
	//fmt.Printf("Providers: %v\n", goth.GetProviders())
}

func getTwitterProvider(conf_twitter config.TwitterConf, conf_path config.OauthPathConf) *twitter.Provider {
	provider_name := util.GetOauthProviderName(conf_twitter.Account.ProviderNamePrefix, conf_path.ProviderNameSuffix)
	tw := twitter.New(conf_twitter.Account.ClientKey, conf_twitter.Account.ClientSecret, getCallbackUrl(conf_path, provider_name))
	tw.SetName(provider_name)
	//fmt.Printf("Callbacks: %v\n", tw.CallbackURL)
	return tw
}

func getGplusProvider(conf_google config.GoogleConf, conf_path config.OauthPathConf) *gplus.Provider {
	provider_name := util.GetOauthProviderName(conf_google.Account.ProviderNamePrefix, conf_path.ProviderNameSuffix)
	gp := gplus.New(conf_google.Account.ClientKey, conf_google.Account.ClientSecret, getCallbackUrl(conf_path, provider_name), conf_path.Scopes...)
	gp.SetPrompt(conf_google.Prompt)
	gp.SetName(provider_name)
	//fmt.Printf("Callbacks: %v\n", gp.CallbackURL)
	return gp
}

func getOpenIdProvider(conf_google config.GoogleConf, conf_path config.OauthPathConf) *openidConnect.Provider {
	provider_name := util.GetOauthProviderName(conf_google.Account.ProviderNamePrefix, conf_path.ProviderNameSuffix)
	oc, err := openidConnect.New(conf_google.Account.ClientKey, conf_google.Account.ClientSecret, getCallbackUrl(conf_path, provider_name), conf_google.AutoDiscoveryUrl, conf_path.Scopes...)
	if err != nil {
		panic(err)
	}
	oc.SetName(provider_name)
	//fmt.Printf("Callbacks: %v\n", oc.CallbackURL)
	return oc
}

func getCallbackUrl(conf_path config.OauthPathConf, provider_name string) string {
	return conf_path.Domain + conf_path.CallbackUri + provider_name + "/callback"
}
