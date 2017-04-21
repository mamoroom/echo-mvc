package cookie

import (
	"github.com/mamoroom/echo-mvc/src/config"
)

var conf = config.Conf
var cookieOf map[string]config.CookieConf

func init() {
	j := conf.Jwt.Web.Cookie
	r := conf.Register.Cookie
	c := conf.Csrf.Cookie
	cookieOf = map[string]config.CookieConf{
		j.Name: j,
		r.Name: r,
		c.Name: c,
	}
}
