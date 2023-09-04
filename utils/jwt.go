package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var Jwtsecret string = "verysecurepassword"

func GenToken(id string) string {
	jwt_token := jwt.New(jwt.GetSigningMethod("HS256"))
	jwt_token.Claims = jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	}
	token, _ := jwt_token.SignedString([]byte(Jwtsecret))
	return token
}
