package config

import (
	"github.com/jinzhu/configor"

	_ "fmt"
	"path"
	"runtime"
)

type config struct {
	Server struct {
		Domain string `json:"domain"`
		Port   string `json:"port"`
	}
	Echo struct {
		IsDebug bool `json:"is_debug"`
		Log     EchoLog
	}
	Db struct {
		User   DbConf
		Master DbConf
	}
	Oauth struct {
		Google  GoogleConf
		Twitter TwitterConf
		Login   OauthPathConf
		App     OauthPathConf
		Chat    OauthPathConf
		Cookie  CookieConf
	}
	Register struct {
		Cookie CookieConf
	}
	Jwt struct {
		Web struct {
			Cookie          CookieConf
			ExpiresDuration int    `json:"expires_duration"`
			SigningKey      string `json:"signing_key"`
		}
		App struct {
			HeaderName      string `json:"header_name"`
			ContextKey      string `json:"context_key"`
			SigningKey      string `json:"signing_key"`
			ExpiresDuration int    `json:"expires_duration"`
		}
	}
	Csrf struct {
		Cookie      CookieConf `json:"cookie"`
		TokenLength uint8      `json:"token_length"`
		ContextKey  string     `json:"context_key"`
	}
}

type EchoLog struct {
	SetLevel string `json:"set_level"`
}

type DbConf struct {
	R          DbConnectConf
	W          DbConnectConf
	DriverName string `json:"driver_name"`
	DbName     string `json:"db_name"`
}

type DbConnectConf struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	IsSqlLog bool   `json:"is_sql_log"`
	LogMode  string `json:"log_mode"`
}

type TwitterConf struct {
	Account OauthAccountConf
}

type GoogleConf struct {
	AutoDiscoveryUrl string `json:"auto_discovery_url"`
	Prompt           string `json:"prompt"`
	Account          OauthAccountConf
}

type OauthAccountConf struct {
	ClientKey          string `json:"client_key"`
	ClientSecret       string `json:"client_secret"`
	ProviderNamePrefix string `json:"provider_name_prefix"`
}

type OauthPathConf struct {
	Domain             string   `json:"domain"`
	CallbackUri        string   `json:"callback_uri"`
	ProviderNameSuffix string   `json:"provider_name_suffix"`
	Scopes             []string `json:"scopes"`
	Session            OauthSessionConf
}

type OauthSessionConf struct {
	CookieName string `json:"cookie_name"`
	SecretKey  string `json:"secret_key"`
}

type CookieConf struct {
	Path     string `json:"path"`
	Domain   string `json:"domain"`
	Name     string `json:"name"`
	MaxAge   int    `json:"max_age"`
	Secure   bool   `json:"secure"`
	HttpOnly bool   `json:"http_only"`
}

var Conf config
var ROOT_PATH string

func init() {
	ROOT_PATH = get_current_dir() + "/.."
	configor.Load(&Conf, ROOT_PATH+"/config/app.json")
}

// __DIR__
func get_current_dir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
