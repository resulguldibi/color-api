package util

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"resulguldibi/color-api/entity"
	"resulguldibi/color-api/types"
	"strings"

	b64 "encoding/base64"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"

	httpClient "resulguldibi/http-client/entity"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func HandleErr(ctx *gin.Context, err interface{}) {
	exp := &types.ExceptionMessage{}
	_ = json.Unmarshal([]byte(fmt.Sprint(err)), exp)
	responseSatus := PrepareResponseStatusWithMessage(false, exp.Message, exp.Code, exp.Stack)
	ctx.JSON(http.StatusBadRequest, responseSatus)
}

func PrepareResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func PrepareResponseStatus(err interface{}) entity.ResponseStatus {
	return entity.ResponseStatus{
		IsSucccess: false,
		Message:    fmt.Sprint(err),
	}
}

func PrepareResponseStatusWithMessage(isSucccess bool, message string, code string, stack string) entity.ResponseStatus {
	return entity.ResponseStatus{
		IsSucccess: isSucccess,
		Message:    message,
		Code:       code,
		Stack:      stack,
	}
}

func GenerateRandomNumber(max int) int {
	return rand.Intn(max)
}

func GenerateRandomColor() *entity.Color {

	color := &entity.Color{}
	color.R = GenerateRandomNumber(256)
	color.G = GenerateRandomNumber(256)
	color.B = GenerateRandomNumber(256)

	return color
}

func IsColorExist(colors []*entity.Color, color *entity.Color) bool {

	var isExist bool = false
	if colors != nil && len(colors) > 0 {
		for _, item := range colors {

			if item != nil && color != nil {

				if color.R == item.R && color.G == item.G && color.B == item.B {
					isExist = true
					break
				}
			}
		}
	}

	return isExist
}

func IsColorsEquals(color1 *entity.Color, color2 *entity.Color) bool {

	var isEquals bool = false
	if color1 != nil && color2 != nil {
		if color1.R == color2.R && color1.G == color2.G && color1.B == color2.B {
			isEquals = true
		}
	}

	return isEquals
}

func GenerateMixColor(colors []*entity.Color) *entity.Color {

	var r int = 0
	var g int = 0
	var b int = 0

	if colors != nil && len(colors) > 0 {
		for _, color := range colors {
			r = r + color.R
			g = g + color.G
			b = b + color.B
		}
	}

	color := &entity.Color{}
	length := len(colors)

	color.R = int(math.Floor(float64(r / length)))
	color.G = int(math.Floor(float64(g / length)))
	color.B = int(math.Floor(float64(b / length)))

	return color
}

func GenerateGuid() string {
	var id uuid.UUID
	var err error
	for {
		id, err = uuid.NewV4()
		if err == nil {
			break
		}
	}
	return id.String()
}

func IsMatchedColors(colors []*entity.Color, color *entity.Color) bool {
	mixedColor := GenerateMixColor(colors)
	isValid := IsColorsEquals(mixedColor, color)
	return isValid
}

func DecodeBase64(encodeData string) ([]byte, error) {
	encodeData = strings.Replace(encodeData, "-", "+", -1)
	encodeData = strings.Replace(encodeData, "_", "/", -1)
	res, _ := b64.RawStdEncoding.DecodeString(encodeData)
	//res, _ := b64.StdEncoding.DecodeString(encodeData)

	return res, nil
}

func GetBase64PayloadFromJWT(jwtToken string) string {
	res := strings.Split(jwtToken, ".")[1]
	return res
}

func GetGoogleIdTokenSignKey(httpClient httpClient.IHttpClient, idToken string) (string, error) {

	var googleOpenIDOAuthCertKey *entity.GoogleOpenIDOAuthCertKey
	googleJWTHeader, err := GetGoogleIdTokenHeaderInfo(idToken)
	if err != nil {
		return "", err
	}

	googleOpenIDOAuthCertKey, err = GetGoogleOpenIDOAuthCertKey(httpClient, googleJWTHeader)

	if googleOpenIDOAuthCertKey == nil {
		return "", types.NewBusinessException("google idtoken sign key exception", "exp.google.id.token.sign.key")
	}

	return googleOpenIDOAuthCertKey.N, nil
}

func GetGoogleIdTokenHeaderInfo(idToken string) (*entity.GoogleJWTHeader, error) {
	jwtToken := GetJWTTokenInfo(idToken)
	googleJWTHeader := &entity.GoogleJWTHeader{}

	headerBytes, err := DecodeBase64(jwtToken.Header)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(headerBytes, googleJWTHeader)

	if err != nil {
		return nil, err
	}

	return googleJWTHeader, nil
}

func GetJWTTokenInfo(jwtToken string) *entity.JWTToken {
	segments := strings.Split(jwtToken, ".")
	return &entity.JWTToken{Header: segments[0], Payload: segments[1], Signature: segments[2]}
}

func GetGoogleOpenIDConfiguration(httpClient httpClient.IHttpClient) (*entity.GoogleOpenIDConfiguration, error) {
	conf := &entity.GoogleOpenIDConfiguration{}
	err := httpClient.Get("https://accounts.google.com/.well-known/openid-configuration").EndStruct(conf)
	return conf, err
}

func GetGoogleOpenIDOAuthCertKey(httpClient httpClient.IHttpClient, jwtHeader *entity.GoogleJWTHeader) (*entity.GoogleOpenIDOAuthCertKey, error) {

	var googleOpenIDOAuthCertKey *entity.GoogleOpenIDOAuthCertKey
	googleOpenIDOAuthCertResponse, err := GetGoogleOpenIDOAuthCerts(httpClient)

	if err != nil {
		return nil, err
	}

	googleOpenIDOAuthCertKey = FindGoogleOpenIDOAuthCertKey(googleOpenIDOAuthCertResponse.Keys, jwtHeader)

	return googleOpenIDOAuthCertKey, nil
}

func GetGoogleOpenIDOAuthCerts(httpClient httpClient.IHttpClient) (*entity.GoogleOpenIDOAuthCertResponse, error) {
	conf, err := GetGoogleOpenIDConfiguration(httpClient)
	certResponse := &entity.GoogleOpenIDOAuthCertResponse{}
	err = httpClient.Get(conf.JwksUri).EndStruct(certResponse)
	return certResponse, err
}

func FindGoogleOpenIDOAuthCertKey(certList []*entity.GoogleOpenIDOAuthCertKey, jwtHeader *entity.GoogleJWTHeader) *entity.GoogleOpenIDOAuthCertKey {
	var foundCert *entity.GoogleOpenIDOAuthCertKey
	if certList != nil && len(certList) > 0 {
		for _, cert := range certList {
			if cert.Alg == jwtHeader.Alg && cert.Kid == jwtHeader.KID {
				foundCert = cert
				break
			}
		}
	}
	return foundCert
}
