package cookie

import (
	"github.com/labstack/echo"

	"net/http"
	"time"

	"github.com/mamoroom/echo-mvc/app/config"
)

func SetCookie(c echo.Context, name string, value string) error {
	cookie_conf := _get_cookie_conf(name)
	cookie := &http.Cookie{
		Name:     cookie_conf.Name,
		Path:     "/", // deleteと仕様を合わせて直打ち
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
	cookie, err := c.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func DeleteCookie(c echo.Context, name string) error {
	cookie, err := c.Cookie(name)
	if err != nil {
		return err
	}
	cookie.MaxAge = -100
	cookie.Path = "/"
	c.SetCookie(cookie)
	return nil
}

func DeleteAllCookies(c echo.Context) error {
	for _, cookie := range c.Cookies() {
		err := DeleteCookie(c, cookie.Name)
		if err != nil {
			return err
		}
	}
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
