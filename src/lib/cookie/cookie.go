package cookie

import (
	"github.com/labstack/echo"

	"net/http"
	"time"

	"github.com/mamoroom/echo-mvc/src/config"
)

func SetCookie(c echo.Context, name string, value string) error {
	cookie_conf := _get_cookie_conf(name)
	cookie := &http.Cookie{
		Name:     cookie_conf.Name,
		Path:     cookie_conf.Path,
		Domain:   cookie_conf.Domain,
		MaxAge:   cookie_conf.MaxAge,
		Secure:   cookie_conf.Secure,
		HttpOnly: cookie_conf.HttpOnly,
		Expires:  time.Now().Add(time.Duration(cookie_conf.MaxAge) * time.Second),
	}
	cookie.Value = value
	c.SetCookie(cookie)
	return nil
}

func GetCookie(c echo.Context, name string) (string, error) {
	cookie_conf := _get_cookie_conf(name)
	cookie, err := c.Cookie(cookie_conf.Name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func DeleteCookie(c echo.Context, name string) error {
	cookie_conf := _get_cookie_conf(name)
	cookie := &http.Cookie{
		Name:     cookie_conf.Name,
		Path:     cookie_conf.Path,
		Domain:   cookie_conf.Domain,
		MaxAge:   cookie_conf.MaxAge,
		Secure:   cookie_conf.Secure,
		HttpOnly: cookie_conf.HttpOnly,
		Expires:  time.Now().Add(time.Duration(cookie_conf.MaxAge) * time.Second),
	}
	cookie.MaxAge = -100
	c.SetCookie(cookie)
	return nil
}

func _get_cookie_conf(name string) config.CookieConf {
	var cookie_conf config.CookieConf
	var ok bool

	if cookie_conf, ok = cookieOf[name]; !ok {
		panic(name + " is not exist on cookie conf.")
	}
	return cookie_conf
}
