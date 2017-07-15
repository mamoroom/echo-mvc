package api

import (
	"github.com/labstack/echo"

	"github.com/mamoroom/echo-mvc/app/router/handler"
	"github.com/mamoroom/echo-mvc/app/view/res_json"
)

type ResNotLoggedIn struct {
	User NotLoggedInUser `json:"user"`
}
type NotLoggedInUser struct {
	IsLoggedIn bool `json:"is_logged_in"`
}

func NotLoggedInHandler(c echo.Context) error {
	res_jwt := c.Get(handler.GetResJwtContextKey()).(*handler.ResJwt)
	res_jwt.Data.Debug.Func = "NotLoggedInHandler"
	//必須
	r := &ResNotLoggedIn{
		User: NotLoggedInUser{
			IsLoggedIn: false,
		},
	}
	return res_json.Succeeded(c, r)
}

func LoginHandler(c echo.Context) error {
	//必須
	res_jwt, _ := c.Get(handler.GetResJwtContextKey()).(*handler.ResJwt)
	res_jwt.Data.Debug.Func = "LoginHandler"
	return res_json.Succeeded(c, res_jwt)
}
