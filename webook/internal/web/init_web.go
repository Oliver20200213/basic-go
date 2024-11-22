package web

import "github.com/gin-gonic/gin"

/*
这是另一种路由注册方式
*/

func RegisterRoutes() *gin.Engine {
	server := gin.Default()
	RegisterUserRoutes(server)

	return server
}

func RegisterUserRoutes(server *gin.Engine) {
	u := &UserHandler{}
	server.POST("/users/signup", u.SignUp)
	//rest风格
	//server.PUT("/user", func(ctx *gin.Context) {
	//
	//})
	server.POST("/users/login", u.Login)
	server.POST("/users/edit", u.Edit)
	//rest 风格
	//server.POST("/users/:id",func(ctx *gin.Context) {
	//})

	server.GET("/users/profile", u.Profile)
	//REST风格
	//server.GET("/users/:id", func(ctx *gin.Context) {
	//})
}
