package main

import (
	"basic-go/practice/internal/repository"
	"basic-go/practice/internal/repository/dao"
	"basic-go/practice/internal/service"
	"basic-go/practice/internal/web"
	"basic-go/practice/internal/web/middleware"
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
func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}
