package main

import (
	"basic-go/practice/config"
	"basic-go/practice/internal/repository"
	"basic-go/practice/internal/repository/cache"
	"basic-go/practice/internal/repository/dao"
	"basic-go/practice/internal/service"
	"basic-go/practice/internal/web"
	"basic-go/practice/internal/web/middleware"
	"basic-go/practice/pkg/ginx/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDB()
	server := initServer()
	cache := initCache()
	u := initUser(db, cache)
	u.RegisterRoutes(server)
	server.Run(":8090")

}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:13316)/webook"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}

func initServer() *gin.Engine {
	server := gin.Default()

	// 引入redis限流
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Minute, 100).Build())

	server.Use(cors.New(cors.Config{
		AllowHeaders:  []string{"Content-Type", "Authorization"},
		ExposeHeaders: []string{"X-Jwt-Token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.Contains(origin, "http://127.0.0.1") {
				return true
			}
			return strings.Contains(origin, "公司所有域名")
		},
		MaxAge: 12 * time.Hour, //校验一次的有效期为12小时，
	}))

	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePath("/users/login").
		IgnorePath("/users/signup.lua").Build())
	return server
}

func initUser(db *gorm.DB, cache *cache.UserCache) *web.UserHandler {
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud, cache)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initCache() *cache.UserCache {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return cache.NewUserCache(redisClient)
}
