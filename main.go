package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/sod-lol/small-cdn/core/models/cache"
	"github.com/sod-lol/small-cdn/core/models/user"
	"github.com/sod-lol/small-cdn/middlewares/token"
	"github.com/sod-lol/small-cdn/routers/auth"
	cache_router "github.com/sod-lol/small-cdn/routers/cache"
	"github.com/sod-lol/small-cdn/services/pg"
	"github.com/sod-lol/small-cdn/services/redis"
)

// CDNRouter is router for whole earch project
type CDNRouter struct {
	*gin.Engine
}

var (
	once     sync.Once
	instance *CDNRouter
)

func createCDNRouter() *CDNRouter {

	once.Do(func() {
		cdn := gin.Default()
		instance = &CDNRouter{
			cdn,
		}
	})

	return instance
}

// GetRouter is function that return gin router of CDNRouter
func GetRouter() *CDNRouter {
	return instance
}

func createSchema(db *gorm.DB) error {
	return db.AutoMigrate(user.User{}, cache.CacheLog{})
}

func dropTables(db *gorm.DB) error {
	return db.Migrator().DropTable(user.User{}, cache.CacheLog{})
}

func reconstructSchema(db *gorm.DB) error {

	err := dropTables(db)
	if err != nil {
		panic(err)
	}
	err = createSchema(db)
	if err != nil {
		panic(err)
	}
	return nil
}

func main() {

	root := context.Background()
	defer root.Done()
	postgres_port := os.Getenv("POSTGRES_PORT")
	port, err := strconv.Atoi(postgres_port)
	if err != nil {
		panic(err)
	}

	session := postgres.CreateCDNDBClient(&postgres.CDNDBInfo{
		Addr:     os.Getenv("POSTGRES_ADDR"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   "small-cdn",
		Port:     int32(port),
		Log:      true,
	})

	db, err := session.DB()
	if err != nil {
		panic(err)
	} else {
		defer db.Close()
	}

	user.CreateUserRepo(session)
	cache.CreateCacheRepo(session)

	// if err := createSchema(session); err != nil {
	// 	panic(err)
	// }
	// uu, err := user.CreateUser("hiiiiii", "sdfs", "sdfksdj@sfkjl.com")
	// if err != nil {
	// 	panic(err)
	// }
	// user.UserRepository.InsertUser(uu)

	reconstructSchema(session)

	// if _, err := configAndSetupDB(session); err != nil {
	// 	logrus.Fatal("[Fatal](main) terminate program due to hot error during initialize tables. error: %v", err)
	// 	return
	// }

	redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PORT"))

	redisAuth := redis.CreateRedisClient(redisAddr, "", 0)
	ctxWithRedis := context.WithValue(root, "redisDB", redisAuth)

	cdn := createCDNRouter()

	jwtMiddleware := token.TokenMiddleWareAuth(redisAuth)

	auth.HandleAuthentication(ctxWithRedis, cdn.Group("/auth"))
	cache_router.HandleCacheing(ctxWithRedis, cdn.Group("/cache"), jwtMiddleware)
	cdn.Run()
}
