package middleware

import (
	"basic-go/practice/internal/web"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePath(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//无需验证
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		//jwt验证
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("QAonYNt3DpoEojWkzJruRYmigFjmfn90"), nil
		})
		fmt.Println("claims::::::;:", claims.Uid)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid || claims.Uid == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//更新过期时间
		now := time.Now()
		//每10s过期一次，  测试时设置的过期时间是1分钟
		if claims.ExpiresAt.Sub(now) < time.Second*50 { //如果过期时间-当前时间小于50s，则进行更新
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString([]byte("QAonYNt3DpoEojWkzJruRYmigFjmfn90"))
			if err != nil {
				//记录到日志中
				log.Println("jwt续约失败", err)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}
		//将claims放到ctx，方便其他地方使用
		ctx.Set("claims", claims)

	}
}
