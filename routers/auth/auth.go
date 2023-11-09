package auth

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/sod-lol/small-cdn/controllers/auth"
	"github.com/sod-lol/small-cdn/services/redis"
)

func HandleAuthentication(ctx context.Context, authRoute *gin.RouterGroup) {

	redisDB := ctx.Value("redisDB").(*redis.Redis)

	authRoute.POST("/login", auth.HandleLogin(redisDB))
	authRoute.POST("/signup", auth.HandleSignUp())
	authRoute.POST("/refresh", auth.HandleRefreshToken(redisDB))
	authRoute.POST("/logout", auth.HandleLogout(redisDB))
}
