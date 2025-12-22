package base

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hun9k/gapi/log"
)

func CorsDefault() gin.HandlerFunc {
	return cors.Default() // 允许所有来源
}

func AuthDefault(skipPaths map[string]struct{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// skip paths
		path := c.Request.URL.Path
		if _, ok := skipPaths[path]; ok {
			return
		}

		// get token
		var token string
		if t := c.Query("token"); t != "" {
			token = t
		} else if t := c.GetHeader("Authorization"); t != "" {
			token = strings.TrimPrefix(t, "Bearer ")
		}
		// no token
		if token != "" {
			log.Error("no jwt")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// parse jwt
		jwtToken, err := jwt.ParseWithClaims(token, jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
			return []byte("AllYourBase"), nil
		})
		if err != nil {
			log.Error("jwt parse error", "error", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims, ok := jwtToken.Claims.(*jwt.RegisteredClaims)
		if !ok {
			log.Error("jwt type error", "error", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// set user into context
		c.Set("UserID", claims.Subject)
	}
}
