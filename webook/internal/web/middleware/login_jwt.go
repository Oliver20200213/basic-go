package middleware

import (
	"basic-go/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

// JWT登录校验

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

// IgnorePaths 中间方法，用于构建部分
func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

// Build终结方法，返回你最终希望的数据
func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//不需要的校验
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		//使用JWT来校验
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//此时tokenHeader的格式：Bearer xdsfdsadsfd最前面是有个Bearer的
		segs := strings.Split(tokenHeader, " ")
		//segs := strings.SplitN(tokenHeader, ".", 2) 或者直接切2次
		if len(segs) != 2 {
			//没有登录,有人捣乱
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claims := &web.UserClaims{}
		//ParseWithClaims里面一定要传入claims指针
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("QAonYNt3DpoEojWkzJruRYmigFjmfn90"), nil
		})
		if err != nil {
			//没有登录，Bearer xxxxx
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//默认是不用自己手动验证过期时间的，会自己验证如果过期token.Valid会为false
		//如果实在想自己验证，也可以自己比较验证
		//if claims.ExpiresAt.Before(time.Now()){
		//	//过期了
		//}
		//err为nil，token不为nil
		if token == nil || !token.Valid || claims.Uid == 0 {
			//没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			//严重的安全问题
			//是需要家监控
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//更新过期时间
		//每十秒钟刷新一次
		now := time.Now()
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			//需要重新生成token
			tokenStr, err = token.SignedString([]byte("QAonYNt3DpoEojWkzJruRYmigFjmfn90"))
			if err != nil {
				//记录日志
				log.Println("jwt续约失败", err)
			}
			//重新返回给前端
			ctx.Header("x-jwt-token", tokenStr)
		}

		ctx.Set("claims", claims)
	}
}
