package util

import (
	"github.com/mamoroom/echo-mvc/app/config"

	"regexp"
	"strings"
)

var conf = config.Conf

//oauth
var loginAuthProviderRegexp *regexp.Regexp
var chatAuthProviderRegexp *regexp.Regexp

//s3 upload
var imageContentTypeRegexp *regexp.Regexp

//game code duplicate entry
var duplicateGameCodeRegexp *regexp.Regexp

//crawler
var crawlerUaRegexp *regexp.Regexp

//carrier domain
var carrierDomainRegexp *regexp.Regexp

func init() {
	_init_login_auth_provider_regexp()
	_init_chat_auth_provider_regexp()
	_init_image_content_type_regexp()
	_init_duplicate_game_code_regexp()
	_init_crawler_ua_regexp()
	_init_carrier_domain_regexp()
}

func _init_login_auth_provider_regexp() {
	regexp_prefix := "^("
	regexp_suffix := ")$"
	regexp_body := []string{GetOauthProviderName(conf.Oauth.Twitter.Account.ProviderNamePrefix, conf.Oauth.Login.ProviderNameSuffix)}
	loginAuthProviderRegexp = regexp.MustCompile(regexp_prefix + strings.Join(regexp_body, "|") + regexp_suffix)
}

func _init_chat_auth_provider_regexp() {
	regexp_prefix := "^("
	regexp_suffix := ")$"
	regexp_body := []string{GetOauthProviderName(conf.Oauth.Twitter.Account.ProviderNamePrefix, conf.Oauth.Post.ProviderNameSuffix)}
	chatAuthProviderRegexp = regexp.MustCompile(regexp_prefix + strings.Join(regexp_body, "|") + regexp_suffix)
}

func _init_image_content_type_regexp() {
	regexp_prefix := "^("
	regexp_suffix := ")$"
	regexp_body := []string{"image/png", "image/jpg", "image/jpeg", "image/gif"}
	imageContentTypeRegexp = regexp.MustCompile(regexp_prefix + strings.Join(regexp_body, "|") + regexp_suffix)
}

func _init_duplicate_game_code_regexp() {
	duplicateGameCodeRegexp = regexp.MustCompile(`Duplicate entry.*for key.*(UQE_user_dscode|UQE_user_appcode)`)
}

func _init_crawler_ua_regexp() {
	regexp_body := []string{"facebookexternalhit/1.1", "twitterbot", "googlebot"}
	crawlerUaRegexp = regexp.MustCompile(strings.Join(regexp_body, "|"))
}

func _init_carrier_domain_regexp() {
	//regexp_body := []string{"@ezweb.ne.jp", "@docomo.ne.jp", "@softbank.ne.jp", "@textmail"}
	regexp_body := []string{"@ezweb.ne.jp", "@docomo.ne.jp", "@softbank.ne.jp"}
	carrierDomainRegexp = regexp.MustCompile(strings.Join(regexp_body, "|"))
}
