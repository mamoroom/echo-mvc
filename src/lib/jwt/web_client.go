package jwt

import (
	jwt_go "github.com/dgrijalva/jwt-go"

	_ "errors"
	"fmt"
	"reflect"
	"time"
)

var web_singing_key = conf.Jwt.Web.SigningKey

type WebCustomClaims struct {
	jwt_go.StandardClaims
}

func CreateWebToken(id string) (string, error) {
	claims := &WebCustomClaims{
		jwt_go.StandardClaims{
			Id:        id,
			ExpiresAt: time.Now().Add(time.Duration(conf.Jwt.Web.ExpiresDuration) * time.Second).Unix(),
		},
	}
	return CreateToken(claims, web_singing_key)
}

func ParseWebToken(t string) (*WebCustomClaims, error) {
	custom_claims := &WebCustomClaims{}

	claims := reflect.ValueOf(custom_claims).Interface().(jwt_go.Claims)
	token, err := jwt_go.ParseWithClaims(t, claims, jwt_go.Keyfunc(verify_secret_key))
	if err != nil || !token.Valid {
		return custom_claims, fmt.Errorf("Could not parse token: %s", err.Error())
	}

	var ok bool
	if custom_claims, ok = token.Claims.(*WebCustomClaims); !ok {
		return custom_claims, fmt.Errorf("Invalid claims format: %v", custom_claims)
	}
	return custom_claims, nil
}

func verify_secret_key(token *jwt_go.Token) (interface{}, error) {
	// Check the signing method
	if _, ok := token.Method.(*jwt_go.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(web_singing_key), nil
}
