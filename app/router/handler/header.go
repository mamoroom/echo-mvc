package handler

import (
	"github.com/labstack/echo"

	"github.com/mamoroom/echo-mvc/app/view/res_json"

	"errors"
	_ "fmt"
	_ "time"
)

func AuthorizationHeaderHandler(auth_scheme string, token string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return _authorization_header_handler_func(next, auth_scheme, token)
	}
}

func _authorization_header_handler_func(next echo.HandlerFunc, auth_scheme string, token string) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		l := len(auth_scheme)
		if len(auth) > l+1 && auth[:l] == auth_scheme && auth[l+1:] == token {
			return next(c)
		}
		return res_json.ErrorBadRequest(c, "InvalidRequest", errors.New("Invalid header value error."), "Requested bearer token is invalid.")
	}
}
