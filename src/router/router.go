package router

import (
	"github.com/labstack/echo"
	"github.com/mamoroom/echo-mvc/src/router/api"
)

func Init(e *echo.Echo) {
	api.Init(e)
}
