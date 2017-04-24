package api

import (
	"github.com/labstack/echo"
	"github.com/markbates/goth"

	"github.com/mamoroom/echo-mvc/src/lib/oauth"
	_ "github.com/mamoroom/echo-mvc/src/lib/util"
	"github.com/mamoroom/echo-mvc/src/models"
	"github.com/mamoroom/echo-mvc/src/router/handler"
	"github.com/mamoroom/echo-mvc/src/view/res_json"

	"errors"
	"fmt"
)

type (
	Chat struct {
		IsTwLoggedIn bool `json:"is_tw_logged_in"`
		//TwLoginLink  string `json:"tw_login_link"`
	}
)

type ResChat struct {
	Chat Chat `json:"chat"`
	handler.ResJwt
}

func ChatHandler(c echo.Context) error {
	//必須
	res_jwt, _ := c.Get(handler.GetResJwtContextKey()).(*handler.ResJwt)
	res_jwt.Data.Debug.Func = "ChatHandler"

	res_chat := &ResChat{
		Chat{
			IsTwLoggedIn: false,
		},
		*res_jwt,
	}

	if !res_jwt.Data.SessionUserModel.GetUserEntity().HasTwitterUserId() {
		//res_chat.Chat.TwLoginLink = util.GetBaseUrl() + "/api/chat/auth/twitter-chat"
	} else {
		res_chat.Chat.IsTwLoggedIn = true
	}
	return res_json.Succeeded(c, res_chat)
}

func ChatAuthHandler(c echo.Context) error {
	res_jwt, _ := c.Get(handler.GetResJwtContextKey()).(*handler.ResJwt)
	res_jwt.Data.Debug.Func = "ChatAuthHandler"

	// auth済
	if gothUser, err := oauth.CompleteUserAuth(c); err == nil {
		return _handle_chat_auth(c, res_jwt, gothUser)
	}

	// auth url取得
	url, err := oauth.GetAuthURL(c)
	if err != nil {
		return res_json.ErrorBadRequest(c, "GetAuthUrlError", err, "Cannot get auth url")
	}
	return res_json.Redirect(c, url)
}

func ChatAuthCallbackHandler(c echo.Context) error {
	res_jwt, _ := c.Get(handler.GetResJwtContextKey()).(*handler.ResJwt)
	res_jwt.Data.Debug.Func = "ChatAuthCallbackHandler"

	gothUser, err := oauth.CompleteUserAuth(c)
	if err != nil {
		return res_json.ErrorBadRequest(c, "InvalidCallbackRequestErr", err, "Invalid callback request")
	}

	return _handle_chat_auth(c, res_jwt, gothUser)
}

func ChatAuthLogoutHandler(c echo.Context) error {
	res_jwt, _ := c.Get(handler.GetResJwtContextKey()).(*handler.ResJwt)
	res_jwt.Data.Debug.Func = "ChatAuthLogoutHandler"
	if !res_jwt.Data.SessionUserModel.GetUserEntity().HasTwitterUserId() {
		return res_json.Failed(c, "AlreadyLogoutFailure")
	}

	////// UserW //////
	user_w := models.NewUserW()
	user_w.Dbh.SetNewSession()
	rollback_func := func() error { return user_w.Dbh.Rollback() }
	commit_func := func() error { return user_w.Dbh.Commit() }
	defer user_w.Dbh.Close()
	defer fmt.Println("End Transaction.")

	user_w.Dbh.ForUpdate()
	if err := user_w.FindUserById(res_jwt.Data.SessionUserModel.GetUserEntity().Id); err != nil {
		handle_rollback_or_commit(rollback_func)
		return res_json.ErrorInternalServer(c, "DatabaseAccessError", err, "Could not get user data from master DB")
	}
	rows_affected, err := user_w.DeleteChatAuth()
	if err != nil {
		handle_rollback_or_commit(rollback_func)
		return res_json.ErrorInternalServer(c, "DbTransactionError", err, "Could not update empty value auth data")
	}
	if rows_affected == 0 {
		handle_rollback_or_commit(rollback_func)
		return res_json.ErrorInternalServer(c, "DbTransactionError", errors.New("Rows afected = 0 on update auth data"), "Could not update user data")
	}
	handle_rollback_or_commit(commit_func)
	///////////////////

	// Cookieを更新
	err = oauth.Logout(c)
	if err != nil {
		return res_json.ErrorInternalServer(c, "CookieError", err, "Could not delete cookie")

	}
	/*
		res_jwt.SetSessionUser(user_model)
		c.Set(handler.GetResJwtContextKey(), res_jwt)
	*/
	return res_json.Succeeded(c, nil)
}

func _handle_chat_auth(c echo.Context, res_jwt *handler.ResJwt, gothUser goth.User) error {
	////// OAuthデータ確認 ///////

	//--------------------//
	//    OAuthデータあり  //
	//--------------------//
	if res_jwt.Data.SessionUserModel.GetUserEntity().HasTwitterUserId() {
		return ChatHandler(c)
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

	//アカウント紐付け
	user_w.Dbh.ForUpdate()
	err := user_w.FindUserById(res_jwt.Data.SessionUserModel.GetUserEntity().Id)
	if err != nil {
		handle_rollback_or_commit(rollback_func)
		return res_json.ErrorInternalServer(c, "DatabaseAccessError", err, "Could not get user data from master DB")
	}
	rows_affected, err := user_w.UpdateTwitterAuth(gothUser.UserID, gothUser.AccessToken, gothUser.AccessTokenSecret, gothUser.AvatarURL)
	if err != nil {
		handle_rollback_or_commit(rollback_func)
		return res_json.ErrorInternalServer(c, "DatabaseAccessError", err, "Could not update user data")
	}
	if rows_affected == 0 {
		handle_rollback_or_commit(rollback_func)
		return res_json.ErrorInternalServer(c, "DbTransactionError", errors.New("Rows afected = 0 on insert auth data"), "Could not insert data")
	}
	handle_rollback_or_commit(commit_func)
	///////////////////

	err = res_jwt.SetSessionUser(user_w)
	if err != nil {
		return res_json.ErrorInternalServer(c, "ResObjectError", err, "Could not set user entity on res obj")
	}
	return ChatHandler(c)
}
