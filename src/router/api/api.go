package api

import (
	"github.com/echo-contrib/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/mamoroom/echo-mvc/src/config"
	"github.com/mamoroom/echo-mvc/src/lib/jwt"
	"github.com/mamoroom/echo-mvc/src/router/handler"
)

var conf = config.Conf

func Init(e *echo.Echo) {

	//For App
	a := e.Group("")
	a.Use(sessions.Sessions(conf.Oauth.App.Session.CookieName, sessions.NewCookieStore([]byte(conf.Oauth.App.Session.SecretKey))))
	//[todo] validate :provider
	/*
		a.GET("/app", AppIndex)
		a.GET(conf.Oauth.App.CallbackUri+":provider", AppAuth)
		a.GET(conf.Oauth.App.CallbackUri+":provider/callback", AppAuthCallback)
	*/

	{ //need authorization header
		a.Use(middleware.JWTWithConfig(jwt.GetMiddlewareJwtConfigForApp()))
		//a.GET("/app/reward", AppRewardIndex)
	}

	//For Web
	// Need Auth Group
	w_jwt := e.Group("")
	w_jwt.Use(sessions.Sessions(conf.Oauth.Login.Session.CookieName, sessions.NewCookieStore([]byte(conf.Oauth.Login.Session.SecretKey))))
	w_jwt.Use(handler.Jwt(NotLoggedInHandler))
	{
		w_jwt.GET("/api/login", LoginHandler)
		w_jwt.GET("/api/nazo", NazoHandler)
		w_jwt.GET("/api/chat", ChatHandler)
		w_jwt.POST("/api/opening", OpeningHandler)

		//[todo] 特別処理: CookieのMaxLength超えてしまうのでここだけ変える. バグならないかcheck...
		w_jwt.Use(sessions.Sessions(conf.Oauth.Chat.Session.CookieName, sessions.NewCookieStore([]byte(conf.Oauth.Chat.Session.SecretKey))))
		w_jwt.GET(conf.Oauth.Chat.CallbackUri+":provider", ChatAuthHandler)

	}

	//ByPass Group
	w_bypass := e.Group("")
	w_bypass.Use(sessions.Sessions(conf.Oauth.Login.Session.CookieName, sessions.NewCookieStore([]byte(conf.Oauth.Login.Session.SecretKey))))
	w_bypass.Use(handler.ByPassJwt())
	{

		w_bypass.POST("/api/register", RegisterHandler)
		w_bypass.POST("/api/login/guest", LoginGuestHandler)
		w_bypass.GET(conf.Oauth.Login.CallbackUri+":provider", LoginAuthHandler)
		w_bypass.GET(conf.Oauth.Login.CallbackUri+":provider/callback", LoginAuthCallbackHandler)

		//[todo] 特別処理: CookieのMaxLength超えてしまうのでここだけ変える. バグならないかcheck...
		w_bypass.Use(sessions.Sessions(conf.Oauth.Chat.Session.CookieName, sessions.NewCookieStore([]byte(conf.Oauth.Chat.Session.SecretKey))))
		w_bypass.GET(conf.Oauth.Chat.CallbackUri+":provider/callback", ChatAuthCallback)
	}
}

func handle_rollback_or_commit(f func() error) {
	err_tx := f()
	// for debug
	if err_tx != nil {
		panic(err_tx)
	}
}
