package token

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sod-lol/small-cdn/services/redis"
)

// TokenMiddleWareAuth is middleware that do token authentication
func TokenMiddleWareAuth(redis *redis.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, userID, err := AccessTokenValidation(redis, c.Request)

		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Set("username", username)
		c.Set("userID", userID)

		c.Next()
	}
}

// AccessTokenValidation is function that validate access token and return username
func AccessTokenValidation(redis *redis.Redis, r *http.Request) (string, uint, error) {

	accessToken, err := VerifyAccessToken(r)

	if err != nil {
		return "", 0, err
	}

	if isContain, _ := redis.Contain(accessToken.AccessTokenUUID); !isContain {
		return "", 0, fmt.Errorf("token is invalid or expired")
	}

	return accessToken.Username, accessToken.UserID, nil

}

// VerifyAccessToken is function that verify and check for sanity of access token
func VerifyAccessToken(r *http.Request) (*AccessTokenDetail, error) {

	extractedAccessToken, err := ExtractTokenFromRequest(r)

	if err != nil {
		return nil, err
	}

	accessToken, err := ExtractAccessTokenFrom(extractedAccessToken)

	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

// ExtractTokenFromRequest is function that extract token from header of request
func ExtractTokenFromRequest(r *http.Request) (string, error) {
	bearToken := r.Header.Get("Authorization")

	splitedTokenAndBear := strings.Split(bearToken, " ")

	if len(splitedTokenAndBear) == 2 {
		return splitedTokenAndBear[1], nil
	}

	return "", fmt.Errorf("cannot split bear from token")
}
