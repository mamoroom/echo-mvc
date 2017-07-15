package handler

import (
	"github.com/labstack/echo"
)

func go_next(next echo.HandlerFunc, c echo.Context, context_key string, res interface{}) error {
	c.Set(context_key, res)
	return next(c)
}
