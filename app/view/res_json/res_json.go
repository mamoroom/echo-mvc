package res_json

import (
	_ "fmt"
	"github.com/labstack/echo"
	"net/http"
)

//interface
type RouterRes interface{}

//res object
type ResObject struct {
	Status    string `json:"status"`
	Type      string `json:"type,omitempty"`
	RouterRes `json:"data"`
}

type ErrorParam struct {
	Type    string
	Raw     error
	Message string
}

// StatusOK //
func Succeeded(c echo.Context, router_res_obj RouterRes) error {
	return _status_ok(c, &ResObject{
		Status:    "succeeded",
		RouterRes: router_res_obj,
	})
}

func Failed(c echo.Context, _type string) error {
	return _status_ok(c, &ResObject{
		Status: "failed",
		Type:   _type,
	})
}

//200
func _status_ok(c echo.Context, res_obj *ResObject) error {
	return c.JSON(http.StatusOK, res_obj)
}

// Client Error //
//400
func ErrorBadRequest(c echo.Context, error_type string, err error, msg string) error {
	return c.JSON(http.StatusBadRequest, echo.Map{"error_type": error_type, "error_message": msg, "error_raw": err.Error()})
}

// Server Error //
//500
func ErrorInternalServer(c echo.Context, error_type string, err error, msg string) error {
	return c.JSON(http.StatusInternalServerError, echo.Map{"error_type": error_type, "error_message": msg, "error_raw": err.Error()})
}
