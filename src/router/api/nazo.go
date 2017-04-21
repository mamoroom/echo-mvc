package api

import (
	"github.com/labstack/echo"

	"github.com/mamoroom/echo-mvc/src/router/handler"
	"github.com/mamoroom/echo-mvc/src/view/res_json"
)

func NazoHandler(c echo.Context) error {
	res_jwt := c.Get(handler.GetResJwtContextKey()).(*handler.ResJwt)
	res_jwt.Data.Debug.Func = "NazoHandler"
	return res_json.Succeeded(c, res_jwt)
}
