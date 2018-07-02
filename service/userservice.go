package service

import (
	"fmt"
	"net/url"
	"os"
	"resulguldibi/color-api/contract"
	"resulguldibi/color-api/entity"
	"resulguldibi/color-api/types"
	"resulguldibi/color-api/util"

	jwt "github.com/dgrijalva/jwt-go"
)

func (service *UserService) GetGoogleOAuthTokenResponse(googleOAuht2token string) (*entity.Token, error) {
	var err error
	var tokenstring string

	data := url.Values{}
	data.Add("code", googleOAuht2token)
	data.Add("client_id", os.Getenv("COLOR_API_GOOGLE_CLIENT_ID"))
	data.Add("client_secret", os.Getenv("COLOR_API_GOOGLE_CLIENT_SECRET"))
	data.Add("redirect_uri", os.Getenv("COLOR_API_GOOGLE_REDIRECT_URI"))
	data.Add("grant_type", "authorization_code")

	googleTokenResponse := &contract.GetGoogleOAuthTokenResponse{}
	err = service.httpClient.PostUrlEncoded("https://www.googleapis.com/oauth2/v4/token", data).EndStruct(googleTokenResponse)

	if err != nil {
		fmt.Println("err -> ", err)
		return nil, types.NewBusinessException("google oauth2 token exception", "exp.google.oauth2.token")
	}

	user := entity.GoogleUser{}

	jwt.ParseWithClaims(googleTokenResponse.IdToken, &user, func(token *jwt.Token) (interface{}, error) {

		key, err := util.GetGoogleIdTokenSignKey(service.httpClient, googleTokenResponse.IdToken)

		if err != nil {
			return nil, err
		}

		return []byte(key), nil
	})

	if err != nil {
		fmt.Println("err -> ", err.Error())
		return nil, types.NewBusinessException("google oauth2 token exception", "exp.google.oauth2.token")
	}

	// Embed User information to `token`
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &entity.User{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Picture:   user.Picture,
	})
	// token -> string. Only server knows this secret (foobar).
	tokenstring, err = token.SignedString([]byte(os.Getenv("COLOR_API_JWT_KEY")))

	if err != nil {
		fmt.Println("err -> ", err.Error())
		return nil, types.NewBusinessException("google oauth2 token exception", "exp.google.oauth2.token")
	}

	tokenData := &entity.Token{AccessToken: tokenstring}

	return tokenData, err
}
