package authentication

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"

	"github.com/sod-lol/small-cdn/core/models/user"
	"github.com/sod-lol/small-cdn/middlewares/token"
	"github.com/sod-lol/small-cdn/services/redis"
)

type loginJsonExpect struct {
	Username string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

func HandleLogin(redisDB *redis.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {
		var authJson loginJsonExpect

		if err := c.ShouldBindJSON(&authJson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		retrivedUser := user.User{
			Username: authJson.Username,
		}

		if err := user.UserRepository.RetrieveUser(&retrivedUser); err != nil {
			if !errors.Is(err, pg.ErrNoRows) {
				logrus.Errorf("Cannot retrive user. error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}

			c.JSON(http.StatusNotFound, gin.H{
				"error": "user with given username not found",
			})

			return
		}

		if !retrivedUser.VerifyPassword(authJson.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "username or password is incorrect",
			})

			return
		}

		genratedToken, err := token.CreateToken(retrivedUser.Username)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		_, err = token.SaveTokenDetail(redisDB, genratedToken, retrivedUser.Username)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access":  genratedToken.GetAccessToken(),
			"refresh": genratedToken.GetRefreshToken(),
		})

	}
}

type signUpJsonExpect struct {
	Username string `form:"username" json:"username" xml:"username"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

func HandleSignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var signUpJson signUpJsonExpect

		if err := c.ShouldBindJSON(&signUpJson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		status, err := GenerateUser(signUpJson)

		if err != nil {
			c.JSON(status, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(status, gin.H{})

	}
}

type expectedRefreshTokenJson struct {
	Refresh string `form:"refresh" json:"refresh" xml:"refresh"  binding:"required"`
}

func HandleRefreshToken(redisDB *redis.Redis) gin.HandlerFunc {

	return func(c *gin.Context) {
		var refreshJson expectedRefreshTokenJson

		if err := c.ShouldBindJSON(&refreshJson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		refreshTokenDetail, err := token.ExtractRefreshTokenFrom(refreshJson.Refresh)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		newToken, err := token.CreateTokenBasedOnRefreshToken(refreshTokenDetail)

		if err != nil {
			logrus.Errorf("Cannot create token from refresh token. error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		_, err = token.SaveTokenDetail(redisDB, newToken, refreshTokenDetail.Username)

		if err != nil {
			logrus.Errorf("Cannot save genrated token. error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access":  newToken.GetAccessToken(),
			"refresh": newToken.GetRefreshToken(),
		})

	}
}

func HandleLogout(redisDB *redis.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetString("username") //come from middleware

		if username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "cannot retrieve username from access token",
			})
			return
		}

		refreshAndAccessUUID, err := token.GetRefreshAndAccessUUIDFrom(redisDB, username)

		if err != nil {
			logrus.Errorf("Cannot retrieve access and refresh token uuid from redis. error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		if _, err := token.DeleteAccessToken(redisDB, refreshAndAccessUUID.AccessUUID); err != nil {
			logrus.Errorf("Cannot delete access uuid from redis. error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		if _, err := token.DeleteRefreshToken(redisDB, refreshAndAccessUUID.RefreshUUID); err != nil {
			logrus.Errorf("Cannot delete refresh uuid from redis. error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		if _, err := redisDB.Delete(username).Result(); err != nil {
			logrus.Errorf("Cannot delete username from redis. error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusNoContent, gin.H{})

	}
}
