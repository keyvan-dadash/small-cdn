package main

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/extra/pgdebug"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"github.com/sod-lol/small-cdn/core/models/user"
	"github.com/sod-lol/small-cdn/routers/auth"
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

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*user.User)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func dropTables(db *pg.DB) error {
	models := []interface{}{
		(*user.User)(nil),
	}

	for _, model := range models {
		err := db.Model(model).DropTable(&orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func reconstructSchema(db *pg.DB) error {

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

	session := postgres.CreateCDNDBClient(&postgres.CDNDBInfo{
		Addr:     "db",
		Username: "postgres",
		Password: "postgres",
		DBName:   "small-cdn",
		Port:     5432,
	})

	session.AddQueryHook(pgdebug.DebugHook{
		// Print all queries.
		Verbose: true,
	})

	defer session.Close()

	user.CreateUserRepo(session)

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

	redisAuth := redis.CreateRedisClient("redis-auth:6379", "", 0)
	ctxWithRedis := context.WithValue(root, "redisDB", redisAuth)

	cdn := createCDNRouter()

	cdn.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	auth.HandleAuthentication(ctxWithRedis, cdn.Group("/auth"))
	cdn.Run()
}
