package handler

import (
	"github.com/labstack/echo"

	"github.com/mamoroom/echo-mvc/app/view/res_json"

	"errors"
	_ "fmt"
	"os"
	_ "time"
)

func DebugHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return _debug_handler_func(next)
	}
}

func _debug_handler_func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		env := os.Getenv("CONFIGOR_ENV")
		switch env {
		case "local":
			return next(c)
		case "dev":
			return next(c)
		case "stg":
			return next(c)
		default:
			return res_json.ErrorBadRequest(c, "InvalidEnv", errors.New("Invalid debug envrioment error."), "This env is invalid for debug: "+env)
		}
	}
}
