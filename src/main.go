package main

import (
	"gopkg.in/go-playground/validator.v9"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/tylerb/graceful"

	"github.com/mamoroom/echo-mvc/src/config"
	"github.com/mamoroom/echo-mvc/src/lib/custom_validator"
	"github.com/mamoroom/echo-mvc/src/lib/util"
	"github.com/mamoroom/echo-mvc/src/router"
	_ "github.com/mamoroom/echo-mvc/src/router/handler"

	_ "fmt"
	_ "os"
	_ "reflect"
	"time"
)

var conf = config.Conf

func main() {
	e := echo.New()
	e.Debug = conf.Echo.IsDebug

	// Res Header //
	//e.Use(middleware.RequestID())
	e.Use(middleware.Secure())
	//e.Use(handler.CSRF())

	e.Use(middleware.Gzip())

	// logger //
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${remote_ip} ${latency_human} [${time_rfc3339}] \"${method} ${uri}\" ${status} ${user_agent}\n",
	}))
	e.Logger.SetLevel(log.DEBUG)
	//e.Logger.SetLevel(reflect.ValueOf(log).FieldByName(conf.Echo.Log.SetLevel))

	// validator //
	e.Validator = &custom_validator.CustomValidator{Validator: validator.New()}
	router.Init(e)

	//start server //
	//e.Logger.Fatal(e.Start())
	e.Server.Addr = util.GetPortStr()
	// [todo:] 調査
	graceful.ListenAndServe(e.Server, 5*time.Second)
}
