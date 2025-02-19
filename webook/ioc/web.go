package ioc

import (
	"basic-go/webook/internal/web"
	"basic-go/webook/internal/web/middleware"
	"basic-go/webook/pkg/ginx/middlewares/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, UserHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	UserHdl.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		loginJWTHdl(),
		rateLimitHdl(redisClient),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{ //use方法注册middleware，use作用于全部路由
		//AllowOrigins: []string{"http://127.0.0.1:3000"}, //允许的访问的源
		//AllowOrigins: []string{"*"},  允许所有请求，不建议这种方式，之前前端是可以的现在前端在严格模式下不生效
		//AllowMethods: []string{"POST","GET"}, //如果不写，默认则是都支持
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"}, //不加这个前端是拿不到x-jwt-token的，意思是我给你的你才能拿到
		AllowCredentials: true,                    //是否允许你带cookie之类的东西

		//可以根据origin进行动态判断(注销掉上面的AllowOrigins配置，二选一配置)
		AllowOriginFunc: func(origin string) bool {
			//如果origin包含http://127.0.0.1就允许访问
			if strings.Contains(origin, "http://127.0.0.1") {
				//你的开发环境
				return true
			}
			return strings.Contains(origin, "公司的域名")
		},
		MaxAge: 12 * time.Hour,
		//这个选项指定浏览器在进行跨域预检请求（Preflight Request）时，能够缓存 CORS 信息的最大时间。也就是说，浏览器在接收到响应后，会在 12 小时内缓存这些 CORS 配置，而不再每次请求时都进行预检。
	})
}

func loginJWTHdl() gin.HandlerFunc {
	return middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").
		IgnorePaths("/users/login_sms/code/send").
		IgnorePaths("/users/login_sms").
		Build()
}

func rateLimitHdl(redisClient redis.Cmdable) gin.HandlerFunc {
	return ratelimit.NewBuilder(redisClient, time.Minute, 100).Build()
}
