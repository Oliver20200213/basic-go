package main

import (
	"basic-go/webook/config"
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/service/sms/memory"
	"basic-go/webook/internal/web"
	"basic-go/webook/internal/web/middleware"
	"basic-go/webook/pkg/ginx/middlewares/ratelimit"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	redis2 "github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDB()
	server := initServer()
	rdb := initRedis()
	u := initUser(db, rdb)
	u.RegisterRoutes(server)

	//server := gin.Default()
	//server.GET("/hello", func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, "这是hello Go页面！")
	//})

	//server.RunTLS()  https
	server.Run(":8080")

}

func initServer() *gin.Engine {
	server := gin.Default()

	server.Use(func(c *gin.Context) { //先注册先执行后注册后执行
		fmt.Println("这是第一个middleware")
	})
	server.Use(func(ctx *gin.Context) {
		fmt.Println("这是第二个middleware")
	})

	//引入redis限流
	redisClient := redis2.NewClient(&redis2.Options{
		//对应k8s-redis-service中metadata的name以及port
		Addr: config.Config.Redis.Addr,
	})
	//1分钟中内只能允许100个请求
	server.Use(ratelimit.NewBuilder(redisClient, time.Minute, 100).Build())

	server.Use(cors.New(cors.Config{ //use方法注册middleware，use作用于全部路由
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
		MaxAge: 12 * time.Hour, //profile的有效期
	}))

	//session实现步骤1
	//单实例可以使用memstore,将数据存储到内存中(常用于测试或开发环境中)
	//多实例使用redistore
	//第一个参数是authentication key（加密key） 第二个是encryption key（加密value）
	//store := memstore.NewStore([]byte("QAonYNt3DpoEojWkzJruRYmigFjmfn90"),
	//	[]byte("pueKIkHTQsCIMa1N7mmzkTN6NmmHjIOP"))
	//server.Use(sessions.Sessions("ssid", store))

	//store, err := redis.NewStore(16, //第一个参数表示最大空闲连接数量，实际中随便写，16,32都行
	//	"tcp", "localhost:6379", "", //协议 登录地址和端口 密码
	//	[]byte("QAonYNt3DpoEojWkzJruRYmigFjmfn90"), //authentication key 指身份认证  32位或64位的key都可以
	//	[]byte("pueKIkHTQsCIMa1N7mmzkTN6NmmHjIOP")) //Encryption 是指数据加密
	//if err != nil {
	//	panic(err)
	//}
	//server.Use(sessions.Sessions("ssid", store))

	//根据面向接口编程 实现自定义sqlxstore
	//mystore := &sqlx_store.Store{}
	//server.Use(sessions.Sessions("ssid", mystore))

	//store := cookie.NewStore([]byte("secret"))   //存储的地方，数据存储到cookie
	//server.Use(sessions.Sessions("ssid", store)) //mysession是cookie中的名字，store是值
	//session实现步骤3
	//server.Use(middlewares.NewLoginMiddlewareBuilder().Build())
	//链式调用,session最好的实现
	//server.Use(middlewares.NewLoginMiddlewareBuilder().
	//	IgnorePaths("/users/login").
	//	IgnorePaths("/users/signup.lua").Build())

	//版本1
	////忽略sss路径
	//middlewares.IgnorePaths = []string{"sss"}
	//server.Use(middlewares.CheckLogin())
	////又有一个server不能忽略sss这个路径,此时v1版本无法实现
	//server1 := gin.Default()
	//server1.Use(middlewares.CheckLogin())

	//使用JWT
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/signup.lua").
		IgnorePaths("/users/login").
		IgnorePaths("/users/login_sms/code/send").
		IgnorePaths("/users/login_sms").
		Build())

	return server
}

func initUser(db *gorm.DB, rdb redis2.Cmdable) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	uc := cache.NewUserCache(rdb) //初始化用户缓存 user cache
	repo := repository.NewRepository(ud, uc)
	svc := service.NewUserService(repo)
	codeCache := cache.NewCodeCache(rdb)
	codeRepo := repository.NewCodeRepository(codeCache)
	smsSvc := memory.NewService()
	codeSvc := service.NewCodeService(codeRepo, smsSvc)
	u := web.NewUserHandler(svc, codeSvc)
	return u
}

func initDB() *gorm.DB {
	//db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:13316)/webook"))
	//连接到k8s的mysql，需要使用k8s-mysql-service.yaml中metadata的name以及port来连接
	//db, err := gorm.Open(mysql.Open("root:root@tcp(webook-mysql:3309)/webook"))
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		//只会再初始化过程中 panic
		//panic相当于证个goroutine结束
		//一段初始化过程出错，应用就不要启动了
		panic(err)
	}

	//初始化建表语句
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}

func initRedis() redis2.Cmdable {
	redisClient := redis2.NewClient(&redis2.Options{
		Addr: config.Config.Redis.Addr,
	})
	return redisClient
}
