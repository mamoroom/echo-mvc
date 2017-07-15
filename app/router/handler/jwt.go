package handler

import (
	"github.com/labstack/echo"

	"github.com/mamoroom/echo-mvc/app/lib/cookie"
	"github.com/mamoroom/echo-mvc/app/lib/custom_time"
	"github.com/mamoroom/echo-mvc/app/lib/jwt"
	"github.com/mamoroom/echo-mvc/app/lib/util"
	"github.com/mamoroom/echo-mvc/app/models"
	"github.com/mamoroom/echo-mvc/app/models/entity"
	"github.com/mamoroom/echo-mvc/app/view/res_json"

	"encoding/json"
	"errors"
	_ "fmt"
	_ "reflect"
	_ "runtime"
	_ "time"
)

type (
	User struct {
		IsLoggedIn bool         `json:"is_logged_in"`
		LoginType  string       `json:"login_type"`
		Entity     *entity.User `json:"entity"`
	}
	Data struct {
		UserId           string            `json:"user_id"`
		SessionUserModel *models.UserModel `json:"-"`
		Debug            Debug             `json:"server"`
	}

	Debug struct {
		Func      string    `json:"func"`
		LoginLink LoginLink `json:"login_link"`
		Msg       string    `json:"msg"`
		Error     string    `json:"error"`
	}
	LoginLink struct {
		Guest  string `json:"guest"`
		Google string `json:"google"`
		//L5ID string `json:"l5id"`
	}
)

func GetResJwtContextKey() string {
	return "res_jwt_handler"
}

func NewResJwt() *ResJwt {
	return &ResJwt{
		User: User{
			IsLoggedIn: false,
		},
		Data: Data{
			Debug: Debug{
				LoginLink: LoginLink{
					Guest:  util.GetBaseUrl() + "/api/login/guest",
					Google: util.GetBaseUrl() + "/api/login/auth/google-login",
				},
			},
		},
	}
}

type ResJwt struct {
	User User `json:"user"`

	// 内部利用, jsonの表示はリリース時に消す
	Data Data `json:"-"`
}

func (res *ResJwt) Marshal2Json() ([]byte, error) {
	return json.Marshal(res)
}

func (res *ResJwt) SetSessionUser(user_model *models.UserModel) error {
	res.ClearSessionUserModel()
	if user_model.IsEntityNil() || user_model.IsEntityEmpty() {
		return errors.New("Could not set empty or nil user to response")
	}
	res.User.Entity = user_model.GetEntity()
	switch res.User.LoginType {
	case "guest":
		res.Data.Debug.LoginLink.Guest = ""
	default:
		res.Data.Debug.LoginLink.Guest = ""
		res.Data.Debug.LoginLink.Google = ""
	}
	res.User.IsLoggedIn = true

	// for internal system
	res.Data.SessionUserModel = user_model
	// いらないかも...
	res.Data.UserId = util.CastUint64ToStr(user_model.GetEntity().Id)
	return nil
}

func (res *ResJwt) ClearSessionUserModel() {
	res.Data.SessionUserModel = nil
}

func ByPassJwtHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return _jwt_handler_func(next, next)
	}
}

func JwtHandler(chg_func echo.HandlerFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return _jwt_handler_func(next, chg_func)
	}
}

func _jwt_handler_func(next echo.HandlerFunc, chg_func echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var _id string
		res := NewResJwt()
		c.Set(custom_time.GetCustomTimeContextKey(), custom_time.Now())

		////// Saveデータ確認 ///////
		//正常系
		token, err := cookie.GetCookie(c, conf.Jwt.Web.Cookie.Name)
		if err != nil {
			res.Data.Debug.Msg = "Cannot read cookies"
			res.Data.Debug.Error = err.Error()
			return go_next(chg_func, c, GetResJwtContextKey(), res)
		}
		claims, err := jwt.ParseWebToken(token)
		if err != nil {
			res.Data.Debug.Msg = "Token invalid"
			res.Data.Debug.Error = err.Error()
			return go_next(chg_func, c, GetResJwtContextKey(), res)
		}

		//[todo] token expiration時の実装

		//異常系
		if _id = claims.Id; len(_id) == 0 {
			return res_json.ErrorBadRequest(c, "InvalidTokenError", err, "Token was found but invalid id format is used")
		}

		////// OAuthデータ確認 ///////
		// DBアクセス
		user_r := models.NewUserR()
		id, err := util.CastStrToInt64(_id)
		if err != nil {
			return res_json.ErrorBadRequest(c, "InvalidTokenError", err, "Token was found but invalid id format is userd for DB")
		}
		err = user_r.GetById(id)
		if err != nil {
			return res_json.ErrorBadRequest(c, "DatabaseAccessError", err, "Token was found but could not get user data")
		} else if user_r.IsEntityEmpty() {
			return res_json.ErrorBadRequest(c, "DatabaseAccessError", errors.New("User data is empty"), "Token was found but user data is empty")
		}

		//for Response
		err = res.SetSessionUser(user_r)
		if err != nil {
			return res_json.ErrorInternalServer(c, "ResObjectError", err, "Could not set session user")
		}

		return go_next(next, c, GetResJwtContextKey(), res)
	}
}
