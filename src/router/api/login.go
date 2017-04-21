package api

import (
	"github.com/labstack/echo"
	"github.com/markbates/goth"

	"github.com/mamoroom/echo-mvc/src/lib/cookie"
	"github.com/mamoroom/echo-mvc/src/lib/jwt"
	"github.com/mamoroom/echo-mvc/src/lib/oauth"
	"github.com/mamoroom/echo-mvc/src/lib/util"
	"github.com/mamoroom/echo-mvc/src/models"
	"github.com/mamoroom/echo-mvc/src/router/handler"
	"github.com/mamoroom/echo-mvc/src/view/res_json"

	"errors"
	"fmt"
	_ "net/http"
	_ "reflect"
)

type ResNotLoggedIn struct {
	User NotLoggedInUser `json:"user"`
}
type NotLoggedInUser struct {
	IsLoggedIn bool `json:"is_logged_in"`
}

func NotLoggedInHandler(c echo.Context) error {
	//必須
	r := &ResNotLoggedIn{
		User: NotLoggedInUser{
			IsLoggedIn: false,
		},
	}
	return res_json.Succeeded(c, r)
}

func LoginHandler(c echo.Context) error {
	//必須
	res_jwt, _ := c.Get(handler.GetResJwtContextKey()).(*handler.ResJwt)
	res_jwt.Data.Debug.Func = "LoginHandler"
	return res_json.Succeeded(c, res_jwt)
}

func LoginGuestHandler(c echo.Context) error {
	res_jwt, ok := c.Get(handler.GetResJwtContextKey()).(*handler.ResJwt)
	if !ok {
		res_jwt = handler.NewResJwt()
	}
	res_jwt.Data.Debug.Func = "LoginGuestHandler"

	// ログイン済み
	if res_jwt.Data.SessionUserModel != nil {
		return res_json.Failed(c, "AlreadyLoggedInFailure")
	}

	// ユーザー選択言語を取得
	lang, err := cookie.GetCookie(c, conf.Register.Cookie.Name)
	if err != nil {
		return res_json.ErrorBadRequest(c, "NoCookieError", err, "Cannot read register cookie")
	}

	////// UserW //////
	user_w := models.NewUserW()
	user_w.Dbh.SetNewSession()
	rollback_func := func() error { return user_w.Dbh.Rollback() }
	commit_func := func() error { return user_w.Dbh.Commit() }
	defer user_w.Dbh.Close()
	defer fmt.Println("End Transaction.")

	// Transaction //
	if err := user_w.Dbh.BeginTx(); err != nil {
		handle_rollback_or_commit(rollback_func)
		return res_json.ErrorInternalServer(c, "DbTransactionError", err, "Could not begin transaction")
	}
	rows_affected, err := user_w.Insert(lang)
	if err != nil {
		handle_rollback_or_commit(rollback_func)
		return res_json.ErrorInternalServer(c, "DbTransactionError", err, "Could not insert data")
	}
	if rows_affected == 0 {
		handle_rollback_or_commit(rollback_func)
		return res_json.ErrorInternalServer(c, "DbTransactionError", errors.New("Rows afected = 0 on insert auth data"), "Could not insert data")
	}
	handle_rollback_or_commit(commit_func)
	///////////////////

	// cookie削除
	cookie.DeleteCookie(c, conf.Register.Cookie.Name)
	// token再発行
	return _go_login_handler_with_creating_token(c, user_w, res_jwt)
}

func LoginAuthHandler(c echo.Context) error {
	res_jwt, ok := c.Get(handler.GetResJwtContextKey()).(*handler.ResJwt)
	if !ok {
		res_jwt = handler.NewResJwt()
	}
	res_jwt.Data.Debug.Func = "LoginAuthHandler"

	// auth済
	if gothUser, err := oauth.CompleteUserAuth(c); err == nil {
		return _handle_login_auth(c, res_jwt, gothUser)
	}

	// auth url取得
	url, err := oauth.GetAuthURL(c)
	if err != nil {
		return res_json.ErrorBadRequest(c, "GetAuthUrlError", err, "Cannot get auth url")
	}
	return res_json.Redirect(c, url)
}

func LoginAuthCallbackHandler(c echo.Context) error {
	res_jwt, ok := c.Get(handler.GetResJwtContextKey()).(*handler.ResJwt)
	if !ok {
		res_jwt = handler.NewResJwt()
	}
	res_jwt.Data.Debug.Func = "LoginAuthCallbackHandler"

	if err := oauth.ValidateState(c); err != nil {
		return res_json.ErrorBadRequest(c, "InvalidCsrfTokenError", err, "Error CSRF token")
	}
	gothUser, err := oauth.CompleteUserAuth(c)
	if err != nil {
		return res_json.ErrorBadRequest(c, "InvalidCallbackRequestErr", err, "Invalid callback request")
	}
	return _handle_login_auth(c, res_jwt, gothUser)
}

func _handle_login_auth(c echo.Context, res_jwt *handler.ResJwt, gothUser goth.User) error {

	////// OAuthデータ確認 ///////
	user_r := models.NewUserR()
	err := user_r.FindUserByAuth(util.GetOauthProviderPrefix(gothUser.Provider), gothUser.UserID)
	if err != nil {
		return res_json.ErrorInternalServer(c, "DatabaseAccessError", err, "Could not get user data")
	}

	//--------------------//
	//    OAuthデータあり  //
	//--------------------//
	if !user_r.IsUserEntityEmpty() {

		//セーブデータあり
		//	 && 連携先のアカウントが紐づいたデータが見つかった場合 -> [todo] 現状はエラー
		if (res_jwt.Data.SessionUserModel != nil) && (res_jwt.Data.SessionUserModel.GetUserEntity().Id != user_r.GetUserEntity().Id) {
			//or そのアカウントでログインする -> new_token_id = user_r.GetUserEntity().Id
			return res_json.Failed(c, "DuplicateAccountFailure")
		}

		//セーブデータなし
		// && 連携先のデータあり -> [ログイン] token再発行
		return _go_login_handler_with_creating_token(c, user_r, res_jwt)
	}

	//--------------------//
	//    OAuthデータなし  //
	//--------------------//

	////// UserW //////
	user_w := models.NewUserW()
	user_w.Dbh.SetNewSession()
	rollback_func := func() error { return user_w.Dbh.Rollback() }
	commit_func := func() error { return user_w.Dbh.Commit() }
	defer user_w.Dbh.Close()
	defer fmt.Println("End Transaction.")

	// Transaction //
	if err := user_w.Dbh.BeginTx(); err != nil {
		handle_rollback_or_commit(rollback_func)
		return res_json.ErrorInternalServer(c, "DbTransactionError", err, "Could not begin transaction")
	}

	// セーブデータあり -> アカウント紐付け
	if res_jwt.Data.SessionUserModel != nil {
		//[todo] need test as user_slave
		user_w.Dbh.ForUpdate()
		if err := user_w.FindUserById(res_jwt.Data.SessionUserModel.GetUserEntity().Id); err != nil {
			handle_rollback_or_commit(rollback_func)
			return res_json.ErrorInternalServer(c, "DatabaseAccessError", err, "Could not get user data from master DB")
		}
		rows_affected, err := user_w.UpdateAuthByUserId(util.GetOauthProviderPrefix(gothUser.Provider), gothUser.UserID, gothUser.AccessToken, gothUser.RefreshToken, gothUser.ExpiresAt)
		if err != nil {
			handle_rollback_or_commit(rollback_func)
			return res_json.ErrorInternalServer(c, "DbTransactionError", err, "Could not insert user data")
		}
		if rows_affected == 0 {
			handle_rollback_or_commit(rollback_func)
			return res_json.ErrorInternalServer(c, "DbTransactionError", errors.New("Rows afected = 0 on insert auth data"), "Could not insert data")
		}
		handle_rollback_or_commit(commit_func)
		return _go_login_handler_updating_user_session(c, user_w, res_jwt)
	}

	// セーブデータなし -> 新規作成
	// ユーザー選択言語を取得
	lang, err := cookie.GetCookie(c, conf.Register.Cookie.Name)
	if err != nil {
		return res_json.ErrorBadRequest(c, "NoCookieError", err, "Cannot read register cookie")
	}

	rows_affected, err := user_w.InsertByAuth(lang, util.GetOauthProviderPrefix(gothUser.Provider), gothUser.UserID, gothUser.AccessToken, gothUser.RefreshToken, gothUser.ExpiresAt)
	if err != nil {
		handle_rollback_or_commit(rollback_func)
		return res_json.ErrorInternalServer(c, "DbTransactionError", err, "Could not insert data")
	}
	if rows_affected == 0 {
		handle_rollback_or_commit(rollback_func)
		return res_json.ErrorInternalServer(c, "DbTransactionError", errors.New("Rows afected = 0 on insert auth data"), "Could not insert data")
	}
	handle_rollback_or_commit(commit_func)
	///////////////////

	// cookie削除
	cookie.DeleteCookie(c, conf.Register.Cookie.Name)
	// token再発行
	return _go_login_handler_with_creating_token(c, user_w, res_jwt)
}

func _go_login_handler_with_creating_token(c echo.Context, user_model *models.UserModel, res_jwt *handler.ResJwt) error {
	t, err := jwt.CreateWebToken(util.CastUint64ToStr(user_model.GetUserEntity().Id))
	if err != nil {
		return res_json.ErrorInternalServer(c, "TokenCreateError", err, "Counld not create token")
	}
	cookie.SetCookie(c, conf.Jwt.Web.Cookie.Name, t)
	return _go_login_handler_updating_user_session(c, user_model, res_jwt)
}

func _go_login_handler_updating_user_session(c echo.Context, user_model *models.UserModel, res_jwt *handler.ResJwt) error {
	res_jwt.SetSessionUser(user_model)
	c.Set(handler.GetResJwtContextKey(), res_jwt)
	return LoginHandler(c)
}
