package api

import (
	"github.com/labstack/echo"

	"github.com/mamoroom/echo-mvc/src/lib/cookie"
	"github.com/mamoroom/echo-mvc/src/lib/util"
	"github.com/mamoroom/echo-mvc/src/router/handler"
	"github.com/mamoroom/echo-mvc/src/view/res_json"

	_ "fmt"
	_ "reflect"
)

type ReqRegister struct {
	Lang string `json:"lang" validate:"required"`
}

func RegisterHandler(c echo.Context) error {
	res_jwt, ok := c.Get(handler.GetResJwtContextKey()).(*handler.ResJwt)
	if !ok {
		res_jwt = handler.NewResJwt()
	}
	res_jwt.Data.Debug.Func = "RegisterHandler"

	// ログイン済み
	if res_jwt.Data.SessionUserModel != nil {
		return res_json.Failed(c, "AlreadyLoggedInFailure")
	}

	var req_reg = new(ReqRegister)
	if err := c.Bind(req_reg); err != nil {
		return res_json.ErrorBadRequest(c, "InvalidRequest", err, "Could not bind request body")
	}
	if err := c.Validate(req_reg); err != nil {
		return res_json.Failed(c, "ValidationFailure")
	}
	if !util.CheckRegexp(`^(ja|en|es|fr|it|de|ko|nl)$`, req_reg.Lang) {
		return res_json.Failed(c, "ValidationFailure")
	}

	cookie.SetCookie(c, conf.Register.Cookie.Name, req_reg.Lang)
	return res_json.Succeeded(c, req_reg)
}
