package entity

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

type ResponseStatus struct {
	IsSucccess bool   `json:"issuccess"`
	Message    string `json:"message"`
	Code       string `json:"code"`
	Stack      string `json:"stack"`
}

type Color struct {
	R          int    `json:"r"`
	G          int    `json:"g"`
	B          int    `json:"b"`
	Name       string `json:"name"`
	IsSelected bool   `json:"-"`
}

type User struct {
	Id        string `json:"id"`
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
	Email     string `json:"email"`
	Picture   string `json:"picture"`
	jwt.StandardClaims
}

type GoogleUser struct {
	Id           string `json:"sub"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	IsValidEmail bool   `json:"email_verified"`

	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	Picture   string `json:"picture"`

	jwt.StandardClaims
}

type JWTToken struct {
	Header    string
	Payload   string
	Signature string
}

type GoogleJWTHeader struct {
	Alg string `json:"alg"`
	KID string `json:"kid"`
}

type GoogleOpenIDConfiguration struct {
	Issuer  string `json:"issuer"`
	JwksUri string `json:"jwks_uri"`
}

type GoogleOpenIDOAuthCertResponse struct {
	Keys []*GoogleOpenIDOAuthCertKey `json:"keys"`
}

type GoogleOpenIDOAuthCertKey struct {
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type UserRaundStepNumber struct {
	Step int `json:"step"`
}

type StageInfo struct {
	UserStageCount            int `json:"userCount"`
	DefaultRaundCountPerStage int `json:"defaultCount"`
	Stage                     int `json:"stage"`
}

type UserRaundColorValidationAttempt struct {
	SendedColors []*Color `json:"sendedColors"`
	MixedColor   *Color   `json:"mixedColor"`
}

type Token struct {
	AccessToken string `json:"token"`
}

type IEntity interface {
	Do()
}

func (color Color) Do() {}

type SocketConnection struct {
	connection *websocket.Conn
}
