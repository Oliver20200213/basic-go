package main

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/web"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDB()
	server := initServer()
	u := initUser(db)
	u.RegisterRoutes(server)
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

	server.Use(cors.New(cors.Config{ //use方法注册middleware，use作用于全部路由
		//AllowOrigins: []string{"http://127.0.0.1:3000"}, //允许的访问的源
		//AllowOrigins: []string{"*"},  允许所有请求，不建议这种方式，之前前端是可以的现在前端在严格模式下不生效
		//AllowMethods: []string{"POST","GET"}, //如果不写，默认则是都支持
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"}, //响应里面带上x-jwt-token 意思就是允许正式业务请求头部携带改head的值
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
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:13316)/webook"))
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
