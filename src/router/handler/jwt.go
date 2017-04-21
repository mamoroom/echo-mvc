package handler

import (
	"github.com/labstack/echo"

	"github.com/mamoroom/echo-mvc/src/lib/cookie"
	"github.com/mamoroom/echo-mvc/src/lib/jwt"
	"github.com/mamoroom/echo-mvc/src/lib/util"
	"github.com/mamoroom/echo-mvc/src/models"
	"github.com/mamoroom/echo-mvc/src/view/res_json"

	"encoding/json"
	"errors"
	_ "fmt"
	_ "reflect"
	_ "runtime"
	"time"
)

var (
	contextKey = "res_jwt_handler"
)

type (
	User struct {
		IsLoggedIn bool       `json:"is_logged_in"`
		LoginType  string     `json:"login_type"`
		Entity     UserEntity `json:"entity"`
	}
	UserEntity struct {
		Name              string    `json:"name"`
		Lang              string    `json:"lang"`
		Points            uint64    `json:"points"`
		Coins             uint64    `json:"coins"`
		IsTutorialDone    bool      `json:"is_tutorial_done"`
		IsFirstGachaDone  bool      `json:"is_first_gacha_done"`
		IsNotificationsOn bool      `json:"is_notifications_on"`
		SeqLoginCnt       uint8     `json:"seq_login_cnt"`
		UpdatedAt         time.Time `json:"updated_at"`
	}
	Data struct {
		UserId           string            `json:"user_id"`
		SessionUserModel *models.UserModel `json:"-"`
		Debug            Debug             `json:"debug"`
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

type ResJwt struct {
	User User `json:"user"`

	// 内部利用
	Data Data `json:"_data"`
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

func (res *ResJwt) Marshal2Json() ([]byte, error) {
	return json.Marshal(res)
}

func (res *ResJwt) SetSessionUser(user_model *models.UserModel) error {
	res.ClearSessionUserModel()
	if user_model.IsUserEntityNil() || user_model.IsUserEntityEmpty() {
		return errors.New("Could not set empty or nil user to response")
	}
	res.User.Entity.Name = user_model.GetUserEntity().Name
	res.User.Entity.Lang = user_model.GetUserEntity().Lang
	res.User.Entity.Points = user_model.GetUserEntity().Points
	res.User.Entity.Coins = user_model.GetUserEntity().Coins
	res.User.Entity.IsTutorialDone = user_model.GetUserEntity().IsTutorialDone
	res.User.Entity.IsTutorialDone = user_model.GetUserEntity().IsTutorialDone
	res.User.Entity.IsFirstGachaDone = user_model.GetUserEntity().IsFirstGachaDone
	res.User.Entity.IsNotificationsOn = user_model.GetUserEntity().IsNotificationsOn
	res.User.Entity.SeqLoginCnt = user_model.GetUserEntity().SeqLoginCnt
	res.User.LoginType = user_model.GetUserEntity().GetLoginType()
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
	res.Data.UserId = util.CastUint64ToStr(user_model.GetUserEntity().Id)
	return nil
}

func (res *ResJwt) ClearSessionUserModel() {
	res.Data.SessionUserModel = nil
}

func ByPassJwt() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return _jwt_handler_func(next, next)
	}
}

func Jwt(chg_func echo.HandlerFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return _jwt_handler_func(next, chg_func)
	}
}

func _jwt_handler_func(next echo.HandlerFunc, chg_func echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var _id string
		res := NewResJwt()

		////// Saveデータ確認 ///////
		//正常系
		token, err := cookie.GetCookie(c, conf.Jwt.Web.Cookie.Name)
		if err != nil {
			res.Data.Debug.Msg = "Cannot read cookies"
			res.Data.Debug.Error = err.Error()
			return go_next(chg_func, c, res)
		}
		claims, err := jwt.ParseWebToken(token)
		if err != nil {
			res.Data.Debug.Msg = "Token invalid"
			res.Data.Debug.Error = err.Error()
			return go_next(chg_func, c, res)
		}

		//[todo] token expiration時の実装

		//異常系
		if _id = claims.Id; len(_id) == 0 {
			return res_json.ErrorBadRequest(c, "InvalidTokenError", err, "Token was found but invalid id format is used")
		}

		////// OAuthデータ確認 ///////
		// DBアクセス
		user_r := models.NewUserR()
		id, err := util.CastStrToUint64(_id)
		if err != nil {
			return res_json.ErrorBadRequest(c, "InvalidTokenError", err, "Token was found but invalid id format is userd for DB")
		}
		err = user_r.FindUserById(id)
		if err != nil {
			return res_json.ErrorBadRequest(c, "DatabaseAccessError", err, "Token was found but could not get user data")
		} else if user_r.IsUserEntityEmpty() {
			return res_json.ErrorBadRequest(c, "DatabaseAccessError", errors.New("User data is empty"), "Token was found but user data is empty")
		}

		//for Response
		err = res.SetSessionUser(user_r)
		if err != nil {
			return res_json.ErrorInternalServer(c, "ResObjectError", err, "Could not set session user")
		}

		return go_next(next, c, res)
	}
}

func GetResJwtContextKey() string {
	return contextKey
}

func go_next(next echo.HandlerFunc, c echo.Context, res *ResJwt) error {
	c.Set(GetResJwtContextKey(), res)
	return next(c)
}
