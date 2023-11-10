package cache

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/sod-lol/small-cdn/controllers/cache"
)

func HandleCacheing(ctx context.Context, cacheRoute *gin.RouterGroup, middlewares ...gin.HandlerFunc) {
	cacheRoute.Use(middlewares...)

	cacheRoute.POST("/add", cache.HandleAddCache())
	cacheRoute.POST("/list", cache.HandleListOfCacheFiles())
}
