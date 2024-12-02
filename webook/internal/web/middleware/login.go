package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
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
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
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
