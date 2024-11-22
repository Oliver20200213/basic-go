package main

import (
	"basic-go/webook/internal/web"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	u := web.NewUserHandler()
	//另一种分组方式
	//u.RegisterRoutesV1(server.Group("/users"))
	u.RegisterRoutes(server)
	server.Run(":8080")
}
