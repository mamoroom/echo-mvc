package jwt

import (
	jwt_go "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	_ "errors"
	"fmt"
	"time"
)

var app_singing_key = conf.Jwt.App.SigningKey

type AppCustomClaims struct {
	jwt_go.StandardClaims
}

func GetMiddlewareJwtConfigForApp() middleware.JWTConfig {
	return middleware.JWTConfig{
		Claims:      &AppCustomClaims{},
		SigningKey:  []byte(app_singing_key),
		TokenLookup: "header:" + conf.Jwt.App.HeaderName,
		ContextKey:  conf.Jwt.App.ContextKey,
	}
}

func CreateAppToken(id string) (string, error) {
	claims := &WebCustomClaims{
		jwt_go.StandardClaims{
			Id: id,
			//ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			ExpiresAt: time.Now().Add(time.Duration(conf.Jwt.Web.ExpiresDuration) * time.Second).Unix(),
		},
	}
	return CreateToken(claims, app_singing_key)
}

func ParseAppToken(c echo.Context) (*AppCustomClaims, error) {
	user, ok := c.Get(conf.Jwt.App.ContextKey).(*jwt_go.Token)
	if !ok {
		return nil, fmt.Errorf("No token data received")
	}
	custom_claims, ok := user.Claims.(*AppCustomClaims)
	if !ok {
		return nil, fmt.Errorf("Invalid claims format: %v", custom_claims)
	}
	return custom_claims, nil
}
