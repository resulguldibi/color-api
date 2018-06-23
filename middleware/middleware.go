package middleware

import (
	"errors"
	"net/http"
	"os"
	"resulguldibi/color-api/entity"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func UseUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authorization := c.Request.Header.Get("Authorization")

		if isAuthorizationRequired(c.Request.URL.Path) {

			// get user-id from client in http request header with key "Authorization" in JWT format
			if authorization == "" {
				c.AbortWithError(http.StatusUnauthorized, errors.New("No Authorization header found"))
			}

			user := entity.User{}
			token, err := jwt.ParseWithClaims(authorization, &user, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("COLOR_API_JWT_KEY")), nil
			})

			if err != nil {
				panic(err)
			}

			if token.Valid {
				c.Set("User", user)
			}
		}

		c.Next()
	}
}

func isAuthorizationRequired(path string) bool {
	unAuthoriziedPaths := []string{"/signin", "/signup"}
	var isAuthorizationRequired bool = true
	if unAuthoriziedPaths != nil && len(unAuthoriziedPaths) > 0 && path != "" {
		for _, unAuthoriziedPath := range unAuthoriziedPaths {
			if unAuthoriziedPath == path {
				isAuthorizationRequired = false
				break
			}
		}
	}

	return isAuthorizationRequired
}