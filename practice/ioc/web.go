package ioc

import (
	"basic-go/practice/internal/web"
	"basic-go/practice/internal/web/middleware"
	"basic-go/practice/pkg/ginx/middlewares/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWebServer(mdl []gin.HandlerFunc, UserHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdl...)
	UserHdl.RegisterRoutes(server)
	return server
}

func InitMiddleware(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		loginJWTHdl(),
		rateLimitHdl(redisClient),
	}

}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowHeaders:  []string{"Content-Type", "Authorization"},
		ExposeHeaders: []string{"X-Jwt-Token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.Contains(origin, "http://127.0.0.1") {
				return true
			}
			return strings.Contains(origin, "公司所有域名")
		},
		MaxAge: 12 * time.Hour, //校验一次的有效期为12小时，
	})
}

func loginJWTHdl() gin.HandlerFunc {
	return middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePath("/users/login").
		IgnorePath("/users/signup").
		IgnorePath("/users/login_sms/send/code").Build()
}

func rateLimitHdl(redisClient redis.Cmdable) gin.HandlerFunc {
	return ratelimit.NewBuilder(redisClient, time.Minute*15, 100).Build()
}
