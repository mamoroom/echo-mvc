package main

import (
	"gopkg.in/go-playground/validator.v9"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/tylerb/graceful"

	"github.com/mamoroom/echo-mvc/app/config"
	"github.com/mamoroom/echo-mvc/app/lib/custom_validator"
	"github.com/mamoroom/echo-mvc/app/lib/logger"
	"github.com/mamoroom/echo-mvc/app/lib/util"
	"github.com/mamoroom/echo-mvc/app/router"
	_ "github.com/mamoroom/echo-mvc/app/router/handler"

	"fmt"
	"os"
	_ "reflect"
	"time"
)

var conf = config.Conf

func main() {
	env := os.Getenv("CONFIGOR_ENV")
	e := echo.New()
	e.Debug = conf.Echo.IsDebug

	// Res Header //
	//e.Use(middleware.RequestID())
	e.Use(middleware.Secure())
	//e.Use(handler.CSRF())

	e.Use(middleware.Gzip())

	// logger //
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${remote_ip}\t${latency_human}\t[${time_rfc3339}]\t\"${method}\t${uri}\"\t${status}\t${user_agent}\n",
		Output: logger.GetRotateLogWriter("server/echo" + "_" + env),
	}))

	switch conf.Echo.Log.SetLevel {
	case "DEBUG":
		e.Logger.SetLevel(log.DEBUG)
	case "INFO":
		e.Logger.SetLevel(log.INFO)
	case "WARN":
		e.Logger.SetLevel(log.WARN)
	case "ERROR":
		e.Logger.SetLevel(log.ERROR)
	case "OFF":
		e.Logger.SetLevel(log.OFF)
	default:
		panic("Invalid conf.Echo.Log.SetLevel | " + conf.Echo.Log.SetLevel)
	}

	fmt.Println("Echo log level=" + conf.Echo.Log.SetLevel)
	// validator //
	e.Validator = &custom_validator.CustomValidator{Validator: validator.New()}
	router.Init(e)

	//start server //
	//e.Logger.Fatal(e.Start())
	e.Server.Addr = util.GetPortStr()
	// [todo:] 調査
	graceful.ListenAndServe(e.Server, 20*time.Second)
}
