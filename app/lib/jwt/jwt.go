package jwt

import (
	jwt_go "github.com/dgrijalva/jwt-go"

	"github.com/mamoroom/echo-mvc/app/config"
)

var conf = config.Conf

//インターフェス引数 -> 値はオブジェクトへの参照
func CreateToken(claims jwt_go.Claims, key string) (string, error) {
	var t string
	token := jwt_go.NewWithClaims(jwt_go.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(key))
	if err != nil {
		return t, err
	}
	return t, nil
}
