package router

import (
	"github.com/mamoroom/echo-mvc/app/router/api"
	"github.com/mamoroom/echo-mvc/app/router/html"
	"github.com/labstack/echo"
)

func Init(e *echo.Echo) {
	api.Init(e)
	html.Init(e)
}
