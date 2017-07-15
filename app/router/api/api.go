package api

import (
	"github.com/echo-contrib/sessions"
	"github.com/labstack/echo"
	_ "github.com/labstack/echo/middleware"

	"github.com/mamoroom/echo-mvc/app/config"
	_ "github.com/mamoroom/echo-mvc/app/lib/jwt"
	"github.com/mamoroom/echo-mvc/app/models"
	"github.com/mamoroom/echo-mvc/app/router/handler"
)

var conf = config.Conf

func Init(e *echo.Echo) {

	//For Lambda
	l := e.Group("")
	l.Use(handler.AuthorizationHeaderHandler("Bearer", conf.Lambda.AuthBearerToken))
	{
	}

	//For Web
	w_jwt := e.Group("")
	w_jwt.Use(sessions.Sessions(conf.Oauth.Login.Session.CookieName, sessions.NewCookieStore([]byte(conf.Oauth.Login.Session.SecretKey))))
	w_jwt.Use(handler.JwtHandler(NotLoggedInHandler))
	{

		// debug
		w_jwt_debug := w_jwt.Group("")
		w_jwt_debug.Use(handler.DebugHandler())
		{
			w_jwt_debug.GET("/api/debug", DebugMockHandler)
		}

		w_jwt_chat := w_jwt.Group("")
		w_jwt_chat.Use(sessions.Sessions(conf.Oauth.Post.Session.CookieName, sessions.NewCookieStore([]byte(conf.Oauth.Post.Session.SecretKey))))
		{
		}
	}

	//ByPass Group
	w_bypass := e.Group("")
	w_bypass.Use(sessions.Sessions(conf.Oauth.Login.Session.CookieName, sessions.NewCookieStore([]byte(conf.Oauth.Login.Session.SecretKey))))
	w_bypass.Use(handler.ByPassJwtHandler())
	{
	}
}

func handle_rollback_or_commit(f func() error) {
	err_tx := f()
	// for debug
	if err_tx != nil {
		panic(err_tx)
	}
}

func go_login_handler(c echo.Context, user_model *models.UserModel, res_jwt *handler.ResJwt) error {
	res_jwt.SetSessionUser(user_model)
	c.Set(handler.GetResJwtContextKey(), res_jwt)
	return LoginHandler(c)
}
