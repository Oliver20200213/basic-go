package main

import (
	"basic-go/practice/wetest/intenal/repository"
	"basic-go/practice/wetest/intenal/repository/dao"
	"basic-go/practice/wetest/intenal/service"
	"basic-go/practice/wetest/intenal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

func main() {
	db := initDB()
	server := initServer()
	u := initUserHandler(db)

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
			return strings.Contains(origin, "all domain name")
		},
	}))

	return server
}
func initUserHandler(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}
