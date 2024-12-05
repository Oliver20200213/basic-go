package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddlewareBuild struct {
	paths []string
}

func NewLoginMiddlewareBuild() *LoginMiddlewareBuild {
	return &LoginMiddlewareBuild{}
}

func (l *LoginMiddlewareBuild) IgnorePath(path string) *LoginMiddlewareBuild {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuild) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
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
