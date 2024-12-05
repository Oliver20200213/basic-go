package main

import (
	"basic-go/practice/intenal/repository"
	"basic-go/practice/intenal/repository/dao"
	"basic-go/practice/intenal/service"
	"basic-go/practice/intenal/web"
	"basic-go/practice/intenal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

func main() {
	db := initDB()
	server := initServer()
	u := initUser(db)
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
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.Contains(origin, "http://127.0.0.1") {
				return true
			}
			return strings.Contains(origin, "公司所有的域名")
		},
	}))

	//单实例部署可以考虑memstore，内存存储（常用于开发和测试环境）
	store := memstore.NewStore([]byte("WqnrQSkoQQ0HVWdqR0suG8Td4uL4IDWE"),
		[]byte("pueKIkHTQsCIMa1N7mmzkTN6NmmHjIOP"))
	server.Use(sessions.Sessions("mysession", store))

	//多实例建议用redistore，redis存储
	//store := cookie.NewStore([]byte("secret"))   //cookie存储
	//server.Use(sessions.Sessions("ssid", store)) //ssid是一定用cookie的
	server.Use(middleware.NewLoginMiddlewareBuild().
		IgnorePath("/users/login").
		IgnorePath("/users/signup").
		Build())

	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}
