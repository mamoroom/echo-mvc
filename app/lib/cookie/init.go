package cookie

import (
	"github.com/mamoroom/echo-mvc/app/config"
)

var conf = config.Conf
var cookieOf map[string]config.CookieConf

func init() {
	j := conf.Jwt.Web.Cookie
	c := conf.Csrf.Cookie
	i := conf.Bypass.Incentive.Cookie
	cookieOf = map[string]config.CookieConf{
		j.Name: j,
		c.Name: c,
		i.Name: i,
	}
}
