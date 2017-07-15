package api

import (
	"github.com/labstack/echo"

	"github.com/mamoroom/echo-mvc/app/router/handler"
	"github.com/mamoroom/echo-mvc/app/view/res_json"
)

// 共通レスポンス
type ResCommon struct {
	Msg      string      `json:"msg"`
	NewValue interface{} `json:"new_value"`
}

// DebugItemConnectDsCodeHandler
type ReqDebugUserAdd struct {
	AddNum  int64  `json:"add_num" validate:"required"`
	AddType string `json:"add_type" validate:"required"`
}

type add_user_param struct {
	msg       string
	action    func(add_num int64)
	get_value func() int64
}

func DebugMockHandler(c echo.Context) error {
	res_jwt := c.Get(handler.GetResJwtContextKey()).(*handler.ResJwt)
	res_jwt.Data.Debug.Func = "DebugMockHandler"
	return res_json.Succeeded(c, &ResCommon{
		Msg:      "",
		NewValue: "",
	})
}
