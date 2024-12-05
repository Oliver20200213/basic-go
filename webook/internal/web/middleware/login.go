package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}
func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}
func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//不需要的校验
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		//不需要的校验
		//if ctx.Request.URL.Path == "/users/login" ||
		//	ctx.Request.URL.Path == "/users/signup" {
		//	return
		//}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")

		//每次访问进行登录校验的时候，同时刷新cookie的有效时间
		//如果1分钟刷新一次，如何知道一分钟过去了？
		//存储一个更新时间，update_time
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		updateTime := sess.Get("update_time")
		sess.Set("userId", id) //必须重新写入userId
		//这里再次配置MaxAge有以下考虑：
		//1.实现滑动过期机制，每次用户请求，都会重置会话有效期为 60 秒，
		//2.防止硬编码或不同路径的会话管理不一致
		//3.防止意外情况导致的会话被清除
		sess.Options(sessions.Options{
			MaxAge: 60,
		})
		now := time.Now().UnixMilli()
		if updateTime == nil {
			//说明还没有刷新过，刚登陆，还没刷新过
			sess.Set("update_time", now)
			sess.Save()
			return
		}
		// updateTime是有的
		updateTimeVal, ok := updateTime.(int64)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if now-updateTimeVal > 10*1000 { //当是超过10的请求之后就会更新update_time
			sess.Set("update_time", now)
			sess.Save()
			return
		}

	}
}

// 另一种实现方式
func (l *LoginMiddlewareBuilder) BuildV1() gin.HandlerFunc {
	//用Go的方式编码解码  需要先注册一下否则会报错
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		//不需要的校验
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		//不需要的校验
		//if ctx.Request.URL.Path == "/users/login" ||
		//	ctx.Request.URL.Path == "/users/signup" {
		//	return
		//}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")

		//每次访问进行登录校验的时候，同时刷新cookie的有效时间
		//如果1分钟刷新一次，如何知道一分钟过去了？
		//存储一个更新时间，update_time
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		updateTime := sess.Get("update_time")
		sess.Set("userId", id) //必须重新写入userId
		//这里再次配置MaxAge有以下考虑：
		//1.实现滑动过期机制，每次用户请求，都会重置会话有效期为 60 秒，
		//2.防止硬编码或不同路径的会话管理不一致
		//3.防止意外情况导致的会话被清除
		sess.Options(sessions.Options{
			MaxAge: 60,
		})
		now := time.Now()
		if updateTime == nil {
			//说明还没有刷新过，刚登陆，还没刷新过
			sess.Set("update_time", now)
			err := sess.Save()
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			return
		}
		// updateTime是有的
		updateTimeVal, ok := updateTime.(time.Time)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if now.Sub(updateTimeVal) > time.Second*10 { //当是超过10的请求之后就会更新update_time
			sess.Set("update_time", now)
			err := sess.Save()
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			return
		}

	}
}

// 实现版本1，最差
var IgnorePaths []string

func CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//不需要登录校验
		for _, path := range IgnorePaths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

// 实现版本2，内部用问题不大
func CheckLoginV1(paths []string,
	abc int,
	cde string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(paths) == 0 {
			paths = []string{""}
		}
		//不需要登录校验
		for _, path := range paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
