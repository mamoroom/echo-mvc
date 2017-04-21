package api

import (
	"github.com/labstack/echo"

	"github.com/mamoroom/echo-mvc/src/lib/jwt"
	"github.com/mamoroom/echo-mvc/src/lib/oauth"
	"github.com/mamoroom/echo-mvc/src/lib/util"
	"github.com/mamoroom/echo-mvc/src/models"
	"github.com/mamoroom/echo-mvc/src/view/res_json"

	"net/http"
	"strconv"
)

func AppIndex(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "oauth_link": "http://localhost:8080/app/auth/google-app"})
}

func AppRewardIndex(c echo.Context) error {
	claims, err := jwt.ParseAppToken(c)
	if err != nil {
		return res_json.ErrorBadRequest(c, "TokenInvalidError", err, "Token invalid")
	}
	return c.JSON(http.StatusOK, echo.Map{"id": claims.Id})
}

func AppAuth(c echo.Context) error {
	// auth済
	if _, err := oauth.CompleteUserAuth(c); err == nil {
		// AppAuthCallback -> ////// OAuthデータ確認 /////// 配下の処理に飛ばす
	}

	// auth url取得
	url, err := oauth.GetAuthURL(c)
	if err != nil {
		return res_json.ErrorBadRequest(c, "GetAuthUrlError", err, "Cannot get auth url")
	}
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func AppAuthCallback(c echo.Context) error {
	if err := oauth.ValidateState(c); err != nil {
		return res_json.ErrorBadRequest(c, "InvalidCsrfTokenError", err, "Error CSRF token")
	}
	gothUser, err := oauth.CompleteUserAuth(c)
	if err != nil {
		return res_json.ErrorBadRequest(c, "InvalidCallbackRequestErr", err, "Invalid callback request")
	}

	////// OAuthデータ確認 ///////
	user_r := models.NewUserR()
	err = user_r.FindUserByAuth(util.GetOauthProviderPrefix(gothUser.Provider), gothUser.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"err": err.Error()})
	}
	if user_r.IsUserEntityEmpty() {
		return c.JSON(http.StatusOK, echo.Map{"status": "failed", "error_type": "UserEntityNotFound", "error_message": "Register and login at campaign sight."})
	}

	t, err := jwt.CreateAppToken(strconv.FormatUint(user_r.GetUserEntity().Id, 10))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"err": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"t": t})
}
