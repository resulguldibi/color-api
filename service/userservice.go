package service

import (
	"os"
	"resulguldibi/color-api/entity"

	jwt "github.com/dgrijalva/jwt-go"
)

func (service *UserService) SignIn(level int64) (*entity.Token, error) {
	var err error
	var tokenstring string

	// Embed User information to `token`
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &entity.User{
		Id: "12345",
	})
	// token -> string. Only server knows this secret (foobar).
	tokenstring, err = token.SignedString([]byte(os.Getenv("COLOR_API_JWT_KEY")))

	tokenData := &entity.Token{AccessToken: tokenstring}

	return tokenData, err
}
