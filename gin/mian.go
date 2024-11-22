package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := gin.Default() //类似于启动一个逻辑上的服务器
	server.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello, go")
	})
	//go func() { //可以同时启动多个（运行在同一个进程中）
	//	server1 := gin.Default()
	//	server1.GET("/hello1", func(c *gin.Context) {
	//		c.String(http.StatusOK, "hello1, go")
	//	})
	//	server1.Run("127.0.0.1:8081")
	//}()

	server.POST("/post", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "this is post method!")
	})

	//restful风格
	//get /users/oliver 查询
	//delete /users/oliver 删除oliver
	//put /users/oliver 注册
	//post /users/oliver 修改 （也有人post用于注册和修改）
	//参数路由（也就是路径中携带的参数，注意：不是url中携带的查询参数）
	//http://127.0.0.1:8080/users/oliver
	server.GET("/users/:name", func(ctx *gin.Context) {
		name := ctx.Param("name") ////获取参数
		ctx.String(http.StatusOK, "hello,这是参数路由"+name)
	})
	//http://127.0.0.1:8080/views/home.html
	//注意通配符路由不能注册这种/users/* /user/*/a 也就是说*不能单独出现
	server.GET("/views/*.html", func(ctx *gin.Context) {
		page := ctx.Param(".html") //获取参数
		ctx.String(http.StatusOK, "hello,这是通配符路由"+page)
	})
	//编译都无法通过
	//server.GET("/views/*/*.html", func(ctx *gin.Context) {})
	//这样写是可以的
	//访问http://127.0.0.1:8080/items/或者http://127.0.0.1:8080/items是都可以访问到的
	//server.GET("/items/", func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, "hello,这是items")
	//})
	//这样也是允许的
	//http://127.0.0.1:8080/items/sfadsfabc
	server.GET("/items/*abc", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello,这是/items/*abc")
	})

	//查询参数 http://127.0.0.1:8080/order?id=110
	server.GET("/order", func(ctx *gin.Context) {
		oid := ctx.Query("id")
		ctx.String(http.StatusOK, "hello，这是查询参数"+oid)
	})

	server.Run("127.0.0.1:8080")
}
